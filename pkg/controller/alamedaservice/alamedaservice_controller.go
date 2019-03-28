package alamedaservice

import (
	"context"
	"time"

	federatoraiv1alpha1 "github.com/containers-ai/federatorai-operator/pkg/apis/federatorai/v1alpha1"
	"github.com/containers-ai/federatorai-operator/pkg/component"
	"github.com/containers-ai/federatorai-operator/pkg/lib/resourceapply"

	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/alamedaserviceparamter"
	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/updateenvvar"
	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/updateparamter"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	alamedaServiceLockName = "alamedaservice-lock"
)

var (
	_   reconcile.Reconciler = &ReconcileAlamedaService{}
	log                      = logf.Log.WithName("controller_alamedaservice")
	//name                             = "kroos-installnamespace"
	gracePeriod = int64(3)
)

// Add creates a new AlamedaService Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	kubeClient, _ := kubernetes.NewForConfig(mgr.GetConfig())
	return &ReconcileAlamedaService{
		client:       mgr.GetClient(),
		scheme:       mgr.GetScheme(),
		apiextclient: apiextension.NewForConfigOrDie(mgr.GetConfig()),
		kubeClient:   kubeClient,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("alamedaservice-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}
	// Watch for changes to primary resource AlamedaService
	err = c.Watch(&source.Kind{Type: &federatoraiv1alpha1.AlamedaService{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner AlamedaService
	/*
		err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &federatoraiv1alpha1.AlamedaService{},
		})
		if err != nil {
			return err
		}

		err = c.Watch(&source.Kind{Type: &corev1.ServiceAccount{}}, &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &federatoraiv1alpha1.AlamedaService{},
		})
		if err != nil {
			return err
		}

		err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &federatoraiv1alpha1.AlamedaService{},
		})
		if err != nil {
			return err
		}

		err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &federatoraiv1alpha1.AlamedaService{},
		})
		if err != nil {
			return err
		}
		err = c.Watch(&source.Kind{Type: &corev1.PersistentVolumeClaim{}}, &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &federatoraiv1alpha1.AlamedaService{},
		})
		if err != nil {
			return err
		}
	*/
	return nil
}

// ReconcileAlamedaService reconciles a AlamedaService object
type ReconcileAlamedaService struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client       client.Client
	scheme       *runtime.Scheme
	apiextclient apiextension.Interface
	kubeClient   *kubernetes.Clientset
}

// Reconcile reads that state of the cluster for a AlamedaService object and makes changes based on the state read
// and what is in the AlamedaService.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileAlamedaService) Reconcile(request reconcile.Request) (reconcile.Result, error) {

	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling AlamedaService")

	// Fetch the AlamedaService instance
	instance := &federatoraiv1alpha1.AlamedaService{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			uninstallResource := alamedaserviceparamter.GetUnInstallResource()
			//r.UninstallDeployment(instance,uninstallResource)
			//r.UninstallService(instance,uninstallResource)
			//r.UninstallConfigMap(instance,uninstallResource)
			r.DeleteRegisterTestsCRD(uninstallResource)
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	needToReconsile, err := r.needToReconsile(instance)
	if err != nil {
		log.V(-1).Info("check if AlamedaService needs to reconsile failed, retry reconciling", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}

	if !needToReconsile {
		log.Info("AlamedaService doe not need to reconsile", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name)
		err := r.syncAlamedaServiceActive(instance, false)
		if err != nil {
			log.V(-1).Info("reconsile AlamedaService failed", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
			return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
		}
		return reconcile.Result{}, nil
	}

	if err := r.syncAlamedaServiceActive(instance, true); err != nil {
		log.V(-1).Info("sync AlamedaService failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}
	asp := alamedaserviceparamter.NewAlamedaServiceParamter(instance)
	installResource := asp.GetInstallResource()
	r.RegisterTestsCRD(installResource)
	if err := r.syncConfigMap(instance, asp, installResource); err != nil {
		log.V(-1).Info("sync configMap failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}
	if err := r.syncService(instance, asp, installResource); err != nil {
		log.V(-1).Info("sync service failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}
	if err := r.syncDeployment(instance, asp, installResource); err != nil {
		log.V(-1).Info("sync deployment failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}
	// if EnableExecution Or EnableGUI has been changed to false
	if !asp.EnableExecution {
		log.Info("EnableExecution has been changed to false")
		excutionResource := alamedaserviceparamter.GetExcutionResource()
		if err := r.UninstallExecutionComponent(instance, excutionResource); err != nil {
			log.V(-1).Info("retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
			return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
		}
	}
	if !asp.EnableGUI {
		log.Info("EnableGUI has been changed to false")
		guiResource := alamedaserviceparamter.GetGUIResource()
		if err := r.UninstallGUIComponent(instance, guiResource); err != nil {
			log.V(-1).Info("retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
			return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
		}
	}
	return reconcile.Result{}, nil
}
func (r *ReconcileAlamedaService) RegisterTestsCRD(resource *alamedaserviceparamter.Resource) {
	for _, fileString := range resource.CustomResourceDefinitionList {
		crd := component.RegistryCustomResourceDefinition(fileString)
		_, _, _ = resourceapply.ApplyCustomResourceDefinition(r.apiextclient.ApiextensionsV1beta1(), crd)
	}
}
func (r *ReconcileAlamedaService) DeleteRegisterTestsCRD(resource *alamedaserviceparamter.Resource) {
	for _, fileString := range resource.CustomResourceDefinitionList {
		crd := component.RegistryCustomResourceDefinition(fileString)
		_, _, _ = resourceapply.DeleteCustomResourceDefinition(r.apiextclient.ApiextensionsV1beta1(), crd)
	}
}
func (r *ReconcileAlamedaService) syncConfigMap(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter, resource *alamedaserviceparamter.Resource) error {
	for _, fileString := range resource.ConfigMapList {
		resourceCM := component.NewConfigMap(fileString)
		if err := controllerutil.SetControllerReference(instance, resourceCM, r.scheme); err != nil {
			return errors.Errorf("Fail resourceCM SetControllerReference: %s", err.Error())
		}
		resourceCM.Namespace = instance.Namespace
		foundCM := &corev1.ConfigMap{}

		err := r.client.Get(context.TODO(), types.NamespacedName{Name: resourceCM.Name, Namespace: resourceCM.Namespace}, foundCM)
		if err != nil && k8sErrors.IsNotFound(err) {
			log.Info("Creating a new Resource ConfigMap... ", "resourceCM.Name", resourceCM.Name)
			resourceCM = updateenvvar.AssignConfigMap(resourceCM, instance.Namespace)
			err = r.client.Create(context.TODO(), resourceCM)
			if err != nil {
				return errors.Errorf("create configMap %s/%s failed: %s", resourceCM.Namespace, resourceCM.Name, err.Error())
			} else {
				log.Info("Successfully Creating Resource ConfigMap", "resourceCM.Name", resourceCM.Name)
			}
		} else if err != nil {
			return errors.Errorf("get configMap %s/%s failed: %s", resourceCM.Namespace, resourceCM.Name, err.Error())
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) syncService(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter, resource *alamedaserviceparamter.Resource) error {
	for _, fileString := range resource.ServiceList {
		resourceSV := component.NewService(fileString)
		if err := controllerutil.SetControllerReference(instance, resourceSV, r.scheme); err != nil {
			return errors.Errorf("Fail resourceSV SetControllerReference: %s", err.Error())
		}
		resourceSV.Namespace = instance.Namespace
		foundSV := &corev1.Service{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: resourceSV.Name, Namespace: resourceSV.Namespace}, foundSV)
		if err != nil && k8sErrors.IsNotFound(err) {
			log.Info("Creating a new Resource Service... ", "resourceSV.Name", resourceSV.Name)
			err = r.client.Create(context.TODO(), resourceSV)
			if err != nil {
				return errors.Errorf("create service %s/%s failed: %s", resourceSV.Namespace, resourceSV.Name, err.Error())
			} else {
				log.Info("Successfully Creating Resource Service", "resourceSV.Name", resourceSV.Name)
			}
		} else if err != nil {
			return errors.Errorf("get service %s/%s failed: %s", resourceSV.Namespace, resourceSV.Name, err.Error())
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) syncDeployment(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter, resource *alamedaserviceparamter.Resource) error {
	for _, fileString := range resource.DeploymentList {
		resourceDep := component.NewDeployment(fileString)
		if err := controllerutil.SetControllerReference(instance, resourceDep, r.scheme); err != nil {
			return errors.Errorf("Fail resourceDep SetControllerReference: %s", err.Error())
		}
		resourceDep.Namespace = instance.Namespace
		foundDep := &appsv1.Deployment{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: resourceDep.Name, Namespace: resourceDep.Namespace}, foundDep)
		if err != nil && k8sErrors.IsNotFound(err) {
			log.Info("Creating a new Resource Deployment... ", "resourceDep.Name", resourceDep.Name)
			resourceDep = updateenvvar.AssignDeployment(resourceDep, instance.Namespace)
			resourceDep = updateparamter.ProcessImageVersion(resourceDep, asp.Version)
			resourceDep = updateparamter.ProcessPrometheusService(resourceDep, asp.PrometheusService)
			err = r.client.Create(context.TODO(), resourceDep)
			if err != nil {
				return errors.Errorf("create deployment %s/%s failed: %s", resourceDep.Namespace, resourceDep.Name, err.Error())
			} else {
				log.Info("Successfully Creating Resource Deployment", "resourceDep.Name", resourceDep.Name)
			}
		} else if err != nil {
			return errors.Errorf("get deployment %s/%s failed: %s", resourceDep.Namespace, resourceDep.Name, err.Error())
		} else {
			update := updateparamter.MatchAlamedaServiceParamter(foundDep, asp.Version, asp.PrometheusService)
			if update {
				log.Info("Update Resource Deployment:", "resourceDep.Name", foundDep.Name)
				foundDep = updateparamter.ProcessImageVersion(foundDep, asp.Version)
				foundDep = updateparamter.ProcessPrometheusService(foundDep, asp.PrometheusService)
				err = r.client.Update(context.TODO(), foundDep)
				if err != nil {
					return errors.Errorf("update deployment %s/%s failed: %s", foundDep.Namespace, foundDep.Name, err.Error())
				}
				log.Info("Successfully Update Resource Deployment", "resourceDep.Name", foundDep.Name)
			}
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) UninstallDeployment(instance *federatoraiv1alpha1.AlamedaService, resource *alamedaserviceparamter.Resource) error {

	for _, fileString := range resource.DeploymentList {
		resourceDep := component.NewDeployment(fileString)
		resourceDep.Namespace = instance.Namespace
		err := r.client.Delete(context.TODO(), resourceDep)
		if err != nil && k8sErrors.IsNotFound(err) {
			return nil
		} else if err != nil {
			return errors.Errorf("delete deployment %s/%s failed: %s", resourceDep.Namespace, resourceDep.Name, err.Error())
		}
	}

	return nil
}

func (r *ReconcileAlamedaService) UninstallService(instance *federatoraiv1alpha1.AlamedaService, resource *alamedaserviceparamter.Resource) error {

	for _, fileString := range resource.ServiceList {
		resourceSVC := component.NewService(fileString)
		resourceSVC.Namespace = instance.Namespace
		err := r.client.Delete(context.TODO(), resourceSVC)
		if err != nil && k8sErrors.IsNotFound(err) {
			return nil
		} else if err != nil {
			return errors.Errorf("delete service %s/%s failed: %s", resourceSVC.Namespace, resourceSVC.Name, err.Error())
		}
	}

	return nil
}

func (r *ReconcileAlamedaService) UninstallConfigMap(instance *federatoraiv1alpha1.AlamedaService, resource *alamedaserviceparamter.Resource) error {

	for _, fileString := range resource.ConfigMapList {
		resourceCM := component.NewConfigMap(fileString)
		resourceCM.Namespace = instance.Namespace
		err := r.client.Delete(context.TODO(), resourceCM)
		if err != nil && k8sErrors.IsNotFound(err) {
			return nil
		} else if err != nil {
			return errors.Errorf("delete comfigMap %s/%s failed: %s", resourceCM.Namespace, resourceCM.Name, err.Error())
		}
	}

	return nil
}

func (r *ReconcileAlamedaService) UninstallGUIComponent(instance *federatoraiv1alpha1.AlamedaService, resource *alamedaserviceparamter.Resource) error {

	if err := r.UninstallDeployment(instance, resource); err != nil {
		return errors.Wrapf(err, "uninstall gui component failed")
	}

	if err := r.UninstallService(instance, resource); err != nil {
		return errors.Wrapf(err, "uninstall gui component failed")
	}

	if err := r.UninstallConfigMap(instance, resource); err != nil {
		return errors.Wrapf(err, "uninstall gui component failed")
	}

	return nil
}

func (r *ReconcileAlamedaService) UninstallExecutionComponent(instance *federatoraiv1alpha1.AlamedaService, resource *alamedaserviceparamter.Resource) error {

	if err := r.UninstallDeployment(instance, resource); err != nil {
		return errors.Wrapf(err, "uninstall execution component failed")
	}

	if err := r.UninstallService(instance, resource); err != nil {
		return errors.Wrapf(err, "uninstall execution component failed")
	}

	return nil
}

func (r *ReconcileAlamedaService) needToReconsile(alamedaService *federatoraiv1alpha1.AlamedaService) (bool, error) {

	lock, err := r.getAlamedaServiceLock(alamedaService.Namespace, alamedaServiceLockName)
	if err == nil {
		if lockIsOwnedByAlamedaService(lock, alamedaService) {
			return true, nil
		} else {
			return false, nil
		}
	} else if k8sErrors.IsNotFound(err) {
		err = r.createAlamedaServiceLock(alamedaService)
		if err == nil {
			return true, nil
		} else if k8sErrors.IsAlreadyExists(err) {
			return false, nil
		} else if err != nil {
			return false, errors.Wrap(err, "check if needs to reconsile failed")
		}
	} else if err != nil {
		return false, errors.Wrap(err, "check if needs to reconsile failed")
	}

	return false, nil
}

func (r *ReconcileAlamedaService) getAlamedaServiceLock(ns, name string) (rbacv1.ClusterRole, error) {

	lock := rbacv1.ClusterRole{}
	err := r.client.Get(context.Background(), types.NamespacedName{Name: name}, &lock)
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			return lock, err
		} else {
			return lock, errors.Errorf("get AlamedaService lock failed: %s", err.Error())
		}
	}

	return lock, nil
}

func (r *ReconcileAlamedaService) createAlamedaServiceLock(alamedaService *federatoraiv1alpha1.AlamedaService) error {

	lock := rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: alamedaServiceLockName,
		},
	}
	if err := controllerutil.SetControllerReference(alamedaService, &lock, r.scheme); err != nil {
		return errors.Errorf("create AlamedaService lock failed: %s", err)
	}

	err := r.client.Create(context.Background(), &lock)
	if err != nil {
		if k8sErrors.IsAlreadyExists(err) {
			return err
		} else {
			return errors.Errorf("create AlamedaService lock failed: %s", err.Error())
		}
	}

	return nil
}

func (r *ReconcileAlamedaService) syncAlamedaServiceActive(alamedaService *federatoraiv1alpha1.AlamedaService, active bool) error {

	copyAlamedaService := alamedaService.DeepCopy()

	if active {
		copyAlamedaService.Status.Conditions = []federatoraiv1alpha1.AlamedaServiceStatusCondition{
			federatoraiv1alpha1.AlamedaServiceStatusCondition{
				Paused: !active,
			},
		}
	} else {
		copyAlamedaService.Status.Conditions = []federatoraiv1alpha1.AlamedaServiceStatusCondition{
			federatoraiv1alpha1.AlamedaServiceStatusCondition{
				Paused:  !active,
				Message: "Other AlamedaService is active.",
			},
		}
	}

	if err := r.client.Update(context.Background(), copyAlamedaService); err != nil {
		return errors.Errorf("update AlamedaService active failed: %s", err.Error())
	}

	return nil
}

func lockIsOwnedByAlamedaService(lock rbacv1.ClusterRole, alamedaService *federatoraiv1alpha1.AlamedaService) bool {

	for _, ownerReference := range lock.OwnerReferences {
		if ownerReference.UID == alamedaService.UID {
			return true
		}
	}

	return false
}
