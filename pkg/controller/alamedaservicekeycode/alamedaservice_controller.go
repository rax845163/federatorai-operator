package alamedaservicekeycode

import (
	"context"
	"fmt"
	"sync"
	"time"

	federatoraiv1alpha1 "github.com/containers-ai/federatorai-operator/pkg/apis/federatorai/v1alpha1"
	client_datahub "github.com/containers-ai/federatorai-operator/pkg/client/datahub"
	"github.com/containers-ai/federatorai-operator/pkg/component"
	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/alamedaserviceparamter"
	repository_keycode "github.com/containers-ai/federatorai-operator/pkg/repository/keycode"
	repository_keycode_datahub "github.com/containers-ai/federatorai-operator/pkg/repository/keycode/datahub"
	"github.com/containers-ai/federatorai-operator/pkg/util"

	"github.com/pkg/errors"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var (
	_               reconcile.Reconciler = &ReconcileAlamedaServiceKeycode{}
	log                                  = logf.Log.WithName("controller_alamedaservicekeycode")
	requeueDuration                      = 1 * time.Second
	finalizerList                        = []string{"keycode.alamedaservices.federatorai.containers.ai"}
)

// Add creates a new AlamedaService Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileAlamedaServiceKeycode{
		client: mgr.GetClient(),
		scheme: mgr.GetScheme(),

		datahubClientMap:     make(map[string]client_datahub.Client),
		datahubClientMapLock: sync.Mutex{},
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("alamedaservicekeycode-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}
	// Watch for changes to primary resource AlamedaService
	err = c.Watch(&source.Kind{Type: &federatoraiv1alpha1.AlamedaService{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}
	return nil
}

// ReconcileAlamedaServiceKeycode reconciles a AlamedaService object
type ReconcileAlamedaServiceKeycode struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme

	datahubClientMap     map[string]client_datahub.Client
	datahubClientMapLock sync.Mutex
}

// Reconcile reconcile AlamedaService's keycode
func (r *ReconcileAlamedaServiceKeycode) Reconcile(request reconcile.Request) (reconcile.Result, error) {

	log.Info("Reconcile Keycode")

	var reconcileError = reconcileError{}
	var state federatoraiv1alpha1.KeycodeState
	var signatureData string
	defer func() {

		instance := &federatoraiv1alpha1.AlamedaService{}
		err := r.client.Get(context.TODO(), client.ObjectKey{Namespace: request.Namespace, Name: request.Name}, instance)
		if err != nil && k8sErrors.IsNotFound(err) {
			addr, err := r.getDatahubAddressByNamespace(request.Namespace)
			if err != nil {
				log.V(-1).Info("Get datahub address failed, skip deleting datahub client", "AlamedaService.Namespace", request.Namespace, "AlamedaService.Name", request.Name, "error", err.Error())
			}
			r.deleteDatahubClient(addr)
			return
		} else if err != nil {
			log.V(-1).Info("Get AlamedaService failed, skip writing status", "AlamedaService.Namespace", request.Namespace, "AlamedaService.Name", request.Name, "error", err.Error())
			return
		}

		instance.SetStatusKeycodeLastErrorMessage(reconcileError.Error())
		instance.SetStatusKeycodeState(state)
		instance.Spec.Keycode.SignatureData = signatureData

		// Get keycodeRepository
		keycodeRepository, err := r.getKeycodeRepository(request.Namespace)
		if err != nil {
			log.V(-1).Info("Get keycode summary failed, will not write keycode summary into AlamedaService's status", "AlamedaService.Namespace", request.Namespace, "AlamedaService.Name", request.Name, "error", err.Error())
		} else {
			detail, err := keycodeRepository.GetKeycodeDetail("")
			if err != nil {
				log.V(-1).Info("Get keycode summary failed, write empty keycode summary into AlamedaService's status", "AlamedaService.Namespace", request.Namespace, "AlamedaService.Name", request.Name, "error", err.Error())
			}
			instance.SetStatusKeycodeSummary(detail.Summary())
		}

		if err := r.client.Update(context.Background(), instance); err != nil {
			log.V(-1).Info("Update AlamedaService status failed", "AlamedaService.Namespace", request.Namespace, "AlamedaService.Name", request.Name, "error", err.Error())
		}
	}()

	// Fetch the AlamedaService instance
	alamedaService := &federatoraiv1alpha1.AlamedaService{}
	err := r.client.Get(context.TODO(), request.NamespacedName, alamedaService)
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			log.Info("AlamedaService not found, skip keycode reconciling", "AlamedaService.Namespace", request.Namespace, "AlamedaService.Name", request.Name)
			reconcileError.err = nil
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.V(-1).Info("Get AlamedaService failed, retry reconciling keycode", "AlamedaService.Namespace", alamedaService.Namespace, "AlamedaService.Name", alamedaService.Name, "error", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: requeueDuration}, nil
	}
	state = alamedaService.Status.KeycodeStatus.State
	signatureData = alamedaService.Spec.Keycode.SignatureData

	// Get keycodeRepository
	keycodeRepository, err := r.getKeycodeRepository(alamedaService.Namespace)
	if err != nil {
		log.V(-1).Info("Get licese repository failed, retry reconciling keycode", "AlamedaService.Namespace", alamedaService.Namespace, "AlamedaService.Name", alamedaService.Name, "error", err.Error())
		reconcileError.err = errors.Wrap(err, "get keycode repository instance failed")
		return reconcile.Result{Requeue: true, RequeueAfter: requeueDuration}, nil
	}

	// Handle deletion of AlamedaService
	if alamedaService.DeletionTimestamp != nil {
		if err := r.deleteAlamedaServiceDependencies(keycodeRepository, alamedaService); err != nil {
			log.V(-1).Info("handle AlamedaService deletion failed, retry reconciling keycode", "AlamedaService.Namespace", alamedaService.Namespace, "AlamedaService.Name", alamedaService.Name, "error", err.Error())
			reconcileError.err = errors.Wrap(err, "handle AlamedaService deletion failed")
			return reconcile.Result{Requeue: true, RequeueAfter: requeueDuration}, nil
		}
		return reconcile.Result{Requeue: false}, nil
	}
	if err := r.setupFinalizers(alamedaService); err != nil {
		log.V(-1).Info("setup finalizers to AlamedaService failed, retry reconciling keycode", "AlamedaService.Namespace", alamedaService.Namespace, "AlamedaService.Name", alamedaService.Name, "error", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: requeueDuration}, nil
	}

	// There are two conditions to handle,
	// first, keycode is empty
	// seconde, keycode is not empty
	if alamedaService.IsCodeNumberEmpty() {
		if err := r.handleEmptyKeycode(keycodeRepository, alamedaService); err != nil {
			log.V(-1).Info("Handle empty keycode failed, retry reconciling keycode", "AlamedaService.Namespace", alamedaService.Namespace, "AlamedaService.Name", alamedaService.Name, "error", err.Error())
			reconcileError.err = errors.Wrap(err, "handle empty keycode failed")
			return reconcile.Result{Requeue: true, RequeueAfter: requeueDuration}, nil
		}
		log.Info("Handle empty keycode done", "AlamedaService.Namespace", alamedaService.Namespace, "AlamedaService.Name", alamedaService.Name)
		return reconcile.Result{}, nil
	}

	// Check if need to reconcile keycode
	if !alamedaService.IsCodeNumberUpdated() &&
		alamedaService.Status.KeycodeStatus.State == federatoraiv1alpha1.KeycodeStateDone {
		log.Info("Keycode not changed, skip reconciling")
		return reconcile.Result{}, nil
	}

	// If keycode is updated, do the update process no matter what the current state is
	if alamedaService.IsCodeNumberUpdated() {
		if err := r.updateKeycode(keycodeRepository, alamedaService); err != nil {
			log.V(-1).Info("Update keycode failed, retry reconciling keycode", "AlamedaService.Namespace", alamedaService.Namespace, "AlamedaService.Name", alamedaService.Name, "error", err.Error())
			reconcileError.err = errors.Wrap(err, "update keycode failed")
			return reconcile.Result{Requeue: true, RequeueAfter: requeueDuration}, nil
		}
		log.Info("Update keycode done, start polling registration data", "AlamedaService.Namespace", alamedaService.Namespace, "AlamedaService.Name", alamedaService.Name)
		return reconcile.Result{Requeue: true, RequeueAfter: requeueDuration}, nil
	}

	// Process keycode by the current state
	switch alamedaService.Status.KeycodeStatus.State {
	case federatoraiv1alpha1.KeycodeStateDefault, federatoraiv1alpha1.KeycodeStateWaitingKeycode:
		if err := r.handleKeycode(keycodeRepository, alamedaService); err != nil {
			log.V(-1).Info("Handling keycode failed, retry reconciling keycode", "AlamedaService.Namespace", alamedaService.Namespace, "AlamedaService.Name", alamedaService.Name, "error", err.Error())
			reconcileError.err = errors.Wrap(err, "handle keycode failed")
			return reconcile.Result{Requeue: true, RequeueAfter: requeueDuration}, nil
		}
		state = federatoraiv1alpha1.KeycodeStatePollingRegistrationData
		log.Info("Handling keycode done, start polling registration data", "AlamedaService.Namespace", alamedaService.Namespace, "AlamedaService.Name", alamedaService.Name)
		return reconcile.Result{Requeue: true, RequeueAfter: requeueDuration}, nil
	case federatoraiv1alpha1.KeycodeStatePollingRegistrationData:
		// This state will move to "federatoraiv1alpha1.KeycodeStateDone" state if the keycode detail is registered

		// Poll registration data from keycode repository
		registrationData, err := keycodeRepository.GetRegistrationData()
		if err != nil {
			log.V(-1).Info("Polling registration data failed, retry reconciling keycode", "AlamedaService.Namespace", alamedaService.Namespace, "AlamedaService.Name", alamedaService.Name, "error", err.Error())
			reconcileError.err = errors.Wrap(err, "poll registration data failed")
			return reconcile.Result{Requeue: true, RequeueAfter: requeueDuration}, nil
		}

		// Get keycode defailt from keycode repository
		detail, err := keycodeRepository.GetKeycodeDetail("")
		if err != nil {
			log.V(-1).Info("Polling registration data failed, retry reconciling keycode", "AlamedaService.Namespace", alamedaService.Namespace, "AlamedaService.Name", alamedaService.Name, "error", err.Error())
			reconcileError.err = errors.Wrap(err, "poll registration data failed")
			return reconcile.Result{Requeue: true, RequeueAfter: requeueDuration}, nil
		}
		if detail.Registered {
			signatureData = registrationData
			state = federatoraiv1alpha1.KeycodeStateDone
			log.Info("Keycode has been registered, move state to \"Done\"", "AlamedaService.Namespace", alamedaService.Namespace, "AlamedaService.Name", alamedaService.Name)
		} else {
			// Update AlamedaService keycode state to WaitingSignatureData and registration data
			if err := r.updateAlamedaServiceStatusKeycodeRegistrationData(alamedaService.Namespace, alamedaService.Name, registrationData); err != nil {
				log.V(-1).Info("Update AlamedaService with registration data failed, retry reconciling keycode", "AlamedaService.Namespace", alamedaService.Namespace, "AlamedaService.Name", alamedaService.Name, "error", err.Error())
				reconcileError.err = errors.Wrap(err, "update registration data failed")
				return reconcile.Result{Requeue: true, RequeueAfter: requeueDuration}, nil
			}
			state = federatoraiv1alpha1.KeycodeStateWaitingSignatureData
			log.Info("Polling registration data done, waiting signature data", "AlamedaService.Namespace", alamedaService.Namespace, "AlamedaService.Name", alamedaService.Name)
		}

		return reconcile.Result{}, nil
	case federatoraiv1alpha1.KeycodeStateWaitingSignatureData:
		if alamedaService.Spec.Keycode.SignatureData == "" {
			log.Info("Waiting signature data to be filled in, skip reconciling", "AlamedaService.Namespace", alamedaService.Namespace, "AlamedaService.Name", alamedaService.Name)
			return reconcile.Result{}, nil
		}
		if err := r.handleSignatureData(keycodeRepository, alamedaService); err != nil {
			log.V(-1).Info("Handling signature data failed, retry reconciling keycode", "AlamedaService.Namespace", alamedaService.Namespace, "AlamedaService.Name", alamedaService.Name, "error", err.Error())
			reconcileError.err = errors.Wrap(err, "handle signature data  failed")
			return reconcile.Result{Requeue: true, RequeueAfter: requeueDuration}, nil
		}
		state = federatoraiv1alpha1.KeycodeStateDone
		log.Info("Handling signature data done", "AlamedaService.Namespace", alamedaService.Namespace, "AlamedaService.Name", alamedaService.Name)
		return reconcile.Result{}, nil
	case federatoraiv1alpha1.KeycodeStateDone:
		if err := r.updateKeycode(keycodeRepository, alamedaService); err != nil {
			log.V(-1).Info("Update keycode failed, retry reconciling keycode", "AlamedaService.Namespace", alamedaService.Namespace, "AlamedaService.Name", alamedaService.Name, "error", err.Error())
			reconcileError.err = errors.Wrap(err, "update keycode failed")
			return reconcile.Result{Requeue: true, RequeueAfter: requeueDuration}, nil
		}
		state = federatoraiv1alpha1.KeycodeStatePollingRegistrationData
		log.Info("Update keycode done, start polling registration data", "AlamedaService.Namespace", alamedaService.Namespace, "AlamedaService.Name", alamedaService.Name)
		return reconcile.Result{Requeue: true, RequeueAfter: requeueDuration}, nil
	default:
		log.Info("Unknown keycode state, skip reconciling", "AlamedaService.Namespace", alamedaService.Namespace, "AlamedaService.Name", alamedaService.Name, "state", alamedaService.Status.KeycodeStatus.State)
		return reconcile.Result{}, nil
	}
}

func (r *ReconcileAlamedaServiceKeycode) setupFinalizers(alamedaService *federatoraiv1alpha1.AlamedaService) error {

	needToAppend := false
	for _, finalizer := range finalizerList {
		if !util.StringInSlice(finalizer, alamedaService.Finalizers) {
			needToAppend = true
			break
		}
	}
	if needToAppend {
		appendFinalizers(alamedaService, finalizerList)
		if err := r.client.Update(context.Background(), alamedaService); err != nil {
			return err
		}
	}

	return nil
}

func (r *ReconcileAlamedaServiceKeycode) deleteAlamedaServiceDependencies(keycodeRepository repository_keycode.Interface, alamedaService *federatoraiv1alpha1.AlamedaService) error {

	if err := r.deleteKeycode(keycodeRepository, alamedaService); err != nil {
		return errors.Wrap(err, "delete keycode failed")
	}
	datahubAddress, err := r.getDatahubAddressByNamespace(alamedaService.Namespace)
	if err != nil {
		return errors.Wrap(err, "get datahub address failed")
	}
	r.deleteDatahubClient(datahubAddress)

	// Remove finalizers from AlamedaService
	deleteFinalizers(alamedaService, finalizerList)
	if err := r.client.Update(context.Background(), alamedaService); err != nil {
		return errors.Errorf("remove finalizers from AlamedaService failed: %s", err.Error())
	}

	return nil
}

func (r *ReconcileAlamedaServiceKeycode) deleteKeycode(keycodeRepository repository_keycode.Interface, alamedaService *federatoraiv1alpha1.AlamedaService) error {

	codeNum := alamedaService.Status.KeycodeStatus.CodeNumber
	if codeNum != "" {
		if err := keycodeRepository.DeleteKeycode(codeNum); err != nil {
			return err
		}
	}

	return nil
}

func (r *ReconcileAlamedaServiceKeycode) deleteDatahubClient(datahubAddr string) {

	if _, exist := r.datahubClientMap[datahubAddr]; exist {
		r.datahubClientMapLock.Lock()
		delete(r.datahubClientMap, datahubAddr)
		r.datahubClientMapLock.Unlock()
	}
}

func (r *ReconcileAlamedaServiceKeycode) handleEmptyKeycode(keycodeRepository repository_keycode.Interface, alamedaService *federatoraiv1alpha1.AlamedaService) error {

	if !alamedaService.IsCodeNumberUpdated() {
		return nil
	}

	// Skip handling if keycode has never been applied
	if alamedaService.Status.KeycodeStatus.State == federatoraiv1alpha1.KeycodeStateDefault {
		return nil
	}

	// Delete keycode to keycode repository
	prevAppliedKeycode := alamedaService.Status.KeycodeStatus.CodeNumber
	if prevAppliedKeycode != "" {
		if err := keycodeRepository.DeleteKeycode(prevAppliedKeycode); err != nil {
			return errors.Wrap(err, "delete keycode from keycode repository failed")
		}
	}

	instance := &federatoraiv1alpha1.AlamedaService{}
	if err := r.client.Get(context.TODO(), client.ObjectKey{Namespace: alamedaService.Namespace, Name: alamedaService.Name}, instance); err != nil {
		return errors.New(err.Error())
	}
	instance.Spec.Keycode = federatoraiv1alpha1.KeycodeSpec{}
	instance.Status.KeycodeStatus = federatoraiv1alpha1.KeycodeStatus{State: federatoraiv1alpha1.KeycodeStateWaitingKeycode}
	if err := r.client.Update(context.Background(), instance); err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func (r *ReconcileAlamedaServiceKeycode) handleKeycode(keycodeRepository repository_keycode.Interface, alamedaService *federatoraiv1alpha1.AlamedaService) error {

	// Apply keycode to keycode repository
	keycode := alamedaService.Spec.Keycode.CodeNumber
	if err := keycodeRepository.SendKeycode(keycode); err != nil {
		return errors.Wrap(err, "send keycode to keycode repository failed")
	}

	// Move AlamedaService keycode state to PollingRegistrationData and update the current handling keycode
	instance := &federatoraiv1alpha1.AlamedaService{}
	if err := r.client.Get(context.TODO(), client.ObjectKey{Namespace: alamedaService.Namespace, Name: alamedaService.Name}, instance); err != nil {
		return errors.Wrap(err, "move AlamedaService keycode state to polling falied")
	}
	instance.Spec.Keycode.SignatureData = ""
	instance.SetStatusCodeNumber(keycode)
	if err := r.client.Update(context.Background(), instance); err != nil {
		return errors.Wrap(err, "move AlamedaService keycode state to polling falied")
	}

	return nil
}

func (r *ReconcileAlamedaServiceKeycode) updateKeycode(keycodeRepository repository_keycode.Interface, alamedaService *federatoraiv1alpha1.AlamedaService) error {

	prevKeycode := alamedaService.Status.KeycodeStatus.CodeNumber
	if prevKeycode != "" {
		if err := keycodeRepository.DeleteKeycode(prevKeycode); err != nil {
			return errors.Wrap(err, fmt.Sprintf("delete previous keycode \"%s\" failed", prevKeycode))
		}
	}

	instance := &federatoraiv1alpha1.AlamedaService{}
	if err := r.client.Get(context.TODO(), client.ObjectKey{Namespace: alamedaService.Namespace, Name: alamedaService.Name}, instance); err != nil {
		return errors.New(err.Error())
	}
	instance.Status.KeycodeStatus = federatoraiv1alpha1.KeycodeStatus{State: federatoraiv1alpha1.KeycodeStateWaitingKeycode}
	if err := r.client.Update(context.Background(), instance); err != nil {
		return errors.New(err.Error())
	}

	if err := r.handleKeycode(keycodeRepository, alamedaService); err != nil {
		return errors.Wrap(err, "handle keycode failed")
	}

	return nil
}

func (r *ReconcileAlamedaServiceKeycode) handleSignatureData(keycodeRepository repository_keycode.Interface, alamedaService *federatoraiv1alpha1.AlamedaService) error {

	// Sending registration data to keycode repository
	err := keycodeRepository.SendSignatureData(alamedaService.Spec.Keycode.SignatureData)
	if err != nil {
		return errors.Wrap(err, "send signature data to keycode repository failed")
	}

	return nil
}

func (r *ReconcileAlamedaServiceKeycode) getKeycodeRepository(namespace string) (repository_keycode.Interface, error) {

	datahubAddress, err := r.getDatahubAddressByNamespace(namespace)
	if err != nil {
		return nil, errors.Wrap(err, "get Datahub address failed")
	}
	if _, exist := r.datahubClientMap[datahubAddress]; !exist {
		r.datahubClientMapLock.Lock()
		datahubClientConfig := client_datahub.NewDefaultConfig()
		datahubClientConfig.Address = datahubAddress
		r.datahubClientMap[datahubAddress] = client_datahub.NewDatahubClient(datahubClientConfig)
		r.datahubClientMapLock.Unlock()
	}
	datahubClient := r.datahubClientMap[datahubAddress]
	keycodeRepository := repository_keycode_datahub.NewKeycodeRepository(&datahubClient)

	return keycodeRepository, nil
}

func (r *ReconcileAlamedaServiceKeycode) getDatahubAddressByNamespace(namespace string) (string, error) {

	componentFactory := component.ComponentConfig{NameSpace: namespace}

	// Get datahub client instance
	datahubServiceAssetName := alamedaserviceparamter.GetAlamedaDatahubService()
	datahubService := componentFactory.NewService(datahubServiceAssetName)
	datahubAddress, err := util.GetServiceAddress(datahubService, "grpc")
	if err != nil {
		return "", err
	}
	return datahubAddress, nil
}

func (r *ReconcileAlamedaServiceKeycode) updateAlamedaServiceStatusKeycodeRegistrationData(namespace, name, data string) error {

	alamedaService := &federatoraiv1alpha1.AlamedaService{}
	r.client.Get(context.TODO(), client.ObjectKey{Namespace: namespace, Name: name}, alamedaService)
	alamedaService.SetStatusRegistrationData(data)
	if err := r.client.Update(context.Background(), alamedaService); err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func (r *ReconcileAlamedaServiceKeycode) updaAlamedaServiceStatusKeycodeSummary(namespace, name, summary string) error {

	alamedaService := &federatoraiv1alpha1.AlamedaService{}
	r.client.Get(context.TODO(), client.ObjectKey{Namespace: namespace, Name: name}, alamedaService)
	alamedaService.SetStatusKeycodeSummary(summary)
	if err := r.client.Update(context.Background(), alamedaService); err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func (r *ReconcileAlamedaServiceKeycode) updateAlamedaServiceKeycodeState(namespace, name string, state federatoraiv1alpha1.KeycodeState) error {

	alamedaService := &federatoraiv1alpha1.AlamedaService{}
	r.client.Get(context.TODO(), client.ObjectKey{Namespace: namespace, Name: name}, alamedaService)
	alamedaService.SetStatusKeycodeState(state)
	if err := r.client.Update(context.Background(), alamedaService); err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func (r *ReconcileAlamedaServiceKeycode) updateAlamedaServiceKeycodeLastError(namespace, name string, err error) error {

	alamedaService := &federatoraiv1alpha1.AlamedaService{}
	r.client.Get(context.TODO(), client.ObjectKey{Namespace: namespace, Name: name}, alamedaService)

	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	alamedaService.SetStatusKeycodeLastErrorMessage(errMsg)
	if err := r.client.Update(context.Background(), alamedaService); err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func (r *ReconcileAlamedaServiceKeycode) updateAlamedaServiceKeycodeSummary(namespace, name string) error {

	// Get keycodeRepository
	keycodeRepository, err := r.getKeycodeRepository(namespace)
	if err != nil {
		return errors.Wrap(err, "get keycode repository failed")
	}

	// Get summary of keycode
	detail, err := keycodeRepository.GetKeycodeDetail("")
	if err != nil {
		return errors.Wrap(err, "get keycode summay from keycode repository failed")
	}

	// Update AlamedaService keycode state to WaitingSignatureData and keycode summary
	if err := r.updaAlamedaServiceStatusKeycodeSummary(namespace, name, detail.Summary()); err != nil {
		return errors.Wrap(err, "update AlamedaService to Done status failed")
	}

	return nil
}

func appendFinalizers(alamedaService *federatoraiv1alpha1.AlamedaService, finalizers []string) {

	existFinalizers := alamedaService.GetFinalizers()

	appendList := make([]string, 0)
	for _, finalizer := range finalizers {
		if !util.StringInSlice(finalizer, existFinalizers) {
			appendList = append(appendList, finalizer)
		}
	}

	alamedaService.Finalizers = append(alamedaService.Finalizers, appendList...)
}

func deleteFinalizers(alamedaService *federatoraiv1alpha1.AlamedaService, finalizers []string) {

	existFinalizers := alamedaService.GetFinalizers()

	preservedList := make([]string, 0)
	for _, existFinalizer := range existFinalizers {
		if !util.StringInSlice(existFinalizer, finalizers) {
			preservedList = append(preservedList, existFinalizer)
		}
	}

	alamedaService.Finalizers = preservedList
}
