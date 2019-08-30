package alamedaservice

import (
	"context"
	"os"
	"time"

	autoscaling_v1alpha1 "github.com/containers-ai/alameda/operator/pkg/apis/autoscaling/v1alpha1"
	federatoraiv1alpha1 "github.com/containers-ai/federatorai-operator/pkg/apis/federatorai/v1alpha1"
	"github.com/containers-ai/federatorai-operator/pkg/component"
	"github.com/containers-ai/federatorai-operator/pkg/lib/resourceapply"
	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec"
	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/alamedaserviceparamter"
	"github.com/containers-ai/federatorai-operator/pkg/updateresource"
	"github.com/containers-ai/federatorai-operator/pkg/util"

	certmanagerv1alpha1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"

	routev1 "github.com/openshift/api/route/v1"
	securityv1 "github.com/openshift/api/security/v1"
	"github.com/pkg/errors"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	ingressv1beta1 "k8s.io/api/extensions/v1beta1"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	apiregistrationv1beta1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1beta1"

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
	_               reconcile.Reconciler = &ReconcileAlamedaService{}
	log                                  = logf.Log.WithName("controller_alamedaservice")
	componentConfig *component.ComponentConfig
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
	util.Disable_operand_resource_protection = os.Getenv("DISABLE_OPERAND_RESOURCE_PROTECTION")
	if util.Disable_operand_resource_protection != "true" {
		err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
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
	}
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
			// uninstallResource := alamedaserviceparamter.GetUnInstallResource()
			// r.UninstallDeployment(instance,uninstallResource)
			// r.UninstallService(instance,uninstallResource)
			// r.UninstallConfigMap(instance,uninstallResource)
			// r.uninstallCustomResourceDefinition(uninstallResource)
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.V(-1).Info("get AlamedaService failed, retry reconciling", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, err
	}

	r.InitAlamedaService(instance)

	// Check if AlamedaService need to reconcile, currently only reconcile one AlamedaService in one cluster
	needToReconcile, err := r.needToReconcile(instance)
	if err != nil {
		log.V(-1).Info("check if AlamedaService needs to reconcile failed, retry reconciling", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}
	if !needToReconcile {
		log.Info("AlamedaService doe not need to reconcile", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name)
		err := r.updateAlamedaServiceActivation(instance, false)
		if err != nil {
			log.V(-1).Info("reconcile AlamedaService failed", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
			return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
		}
		return reconcile.Result{}, nil
	} else {
		if err := r.updateAlamedaServiceActivation(instance, true); err != nil {
			log.V(-1).Info("reconcile AlamedaService failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
			return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
		}
	}
	if flag, _ := r.checkAlamedaServiceSpecIsChange(instance, request.NamespacedName); !flag && util.Disable_operand_resource_protection == "true" {
		log.Info("AlamedaService spec is not changed, skip reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name)
		return reconcile.Result{}, nil
	}
	asp := alamedaserviceparamter.NewAlamedaServiceParamter(instance)
	ns, err := r.getNamespace(request.Namespace)
	if err != nil {
		log.V(-1).Info("get namespace failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}
	componentConfig = r.newComponentConfig(ns)
	installResource := asp.GetInstallResource()
	if err = r.syncCustomResourceDefinition(instance, asp, installResource); err != nil {
		log.Error(err, "create crd failed")
	}
	if err := r.updateNamespaceLabel(instance.Namespace); err != nil {
		log.V(-1).Info("update Namespace's label failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}
	if err := r.createAlamedaNotificationChannels(instance, installResource); err != nil {
		log.V(-1).Info("create AlamedaNotificationChannels failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
	}
	if err := r.createAlamedaNotificationTopics(instance, installResource); err != nil {
		log.V(-1).Info("create AlamedaNotificationTopic failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
	}
	if err := r.syncPodSecurityPolicy(instance, asp, installResource); err != nil {
		log.V(-1).Info("sync podSecurityPolicy failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}
	if err := r.syncSecurityContextConstraints(instance, asp, installResource); err != nil {
		log.V(-1).Info("sync securityContextConstraint failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
	}
	if err := r.syncClusterRole(instance, asp, installResource); err != nil {
		log.V(-1).Info("sync clusterRole failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}
	if err := r.syncServiceAccount(instance, asp, installResource); err != nil {
		log.V(-1).Info("sync serviceAccount failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}
	if err := r.syncClusterRoleBinding(instance, asp, installResource); err != nil {
		log.V(-1).Info("sync clusterRoleBinding failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}
	if err := r.syncRole(instance, asp, installResource); err != nil {
		log.V(-1).Info("sync Role failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}
	if err := r.syncRoleBinding(instance, asp, installResource); err != nil {
		log.V(-1).Info("sync RoleBinding failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}
	if err := r.syncIssuer(instance, asp, installResource); err != nil {
		log.V(-1).Info("sync Issuer failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}
	if err := r.syncCertificate(instance, asp, installResource); err != nil {
		log.V(-1).Info("sync Certificate failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}
	if err := r.createSecret(instance, asp, installResource); err != nil {
		log.V(-1).Info("create secret failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}
	if err := r.createPersistentVolumeClaim(instance, asp, installResource); err != nil {
		log.V(-1).Info("create PersistentVolumeClaim failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}
	if err := r.syncConfigMap(instance, asp, installResource); err != nil {
		log.V(-1).Info("sync configMap failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}
	if err := r.syncService(instance, asp, installResource); err != nil {
		log.V(-1).Info("sync service failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}
	if err := r.syncMutatingWebhookConfiguration(instance, asp, installResource); err != nil {
		log.V(-1).Info("sync MutatingWebhookConfiguration failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}
	if err := r.syncValidatingWebhookConfiguration(instance, asp, installResource); err != nil {
		log.V(-1).Info("sync ValidatingWebhookConfiguration failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}
	if err := r.syncAPIService(instance, asp, installResource); err != nil {
		log.V(-1).Info("sync APIService failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}
	if err := r.syncDeployment(instance, asp, installResource); err != nil {
		log.V(-1).Info("sync deployment failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}
	if err := r.syncStatefulSet(instance, asp, installResource); err != nil {
		log.V(-1).Info("sync statefulset failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}
	if err := r.syncIngress(instance, asp, installResource); err != nil {
		log.V(-1).Info("sync Ingress failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}
	if err := r.syncRoute(instance, asp, installResource); err != nil {
		log.V(-1).Info("sync route failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
	}
	if err := r.syncDaemonSet(instance, asp, installResource); err != nil {
		log.V(-1).Info("sync DaemonSet failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}
	// if EnableExecution Or EnableGUI has been changed to false
	//Uninstall Execution Component
	if !asp.EnableExecution {
		log.Info("EnableExecution has been changed to false")
		excutionResource := alamedaserviceparamter.GetExcutionResource()
		if err := r.uninstallExecutionComponent(instance, excutionResource); err != nil {
			log.V(-1).Info("retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
			return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
		}
	}
	//Uninstall GUI Component
	if !asp.EnableGUI {
		log.Info("EnableGUI has been changed to false")
		guiResource := alamedaserviceparamter.GetGUIResource()
		if err := r.uninstallGUIComponent(instance, guiResource); err != nil {
			log.V(-1).Info("retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
			return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
		}
	}
	//Uninstall dispatcher Component
	if !asp.EnableDispatcher {
		log.Info("EnableDispatcher has been changed to false")
		dispatcherResource := alamedaserviceparamter.GetDispatcherResource()
		if err := r.uninstallExecutionComponent(instance, dispatcherResource); err != nil {
			log.V(-1).Info("retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
			return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
		}
	}
	//Uninstall alameter Component
	if !asp.EnableFedemeter {
		log.Info("EnableFedemeter has been changed to false")
		fedemeterResource := alamedaserviceparamter.GetFedemeterResource()
		if err := r.uninstallFedemeterComponent(instance, fedemeterResource); err != nil {
			log.V(-1).Info("retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
			return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
		}
	}
	//Uninstall weavescope components
	if !asp.EnableWeavescope {
		weavescopeResource := alamedaserviceparamter.GetWeavescopeResource()
		if err := r.uninstallResource(weavescopeResource); err != nil {
			log.V(-1).Info("retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
			return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
		}
	}
	//Uninstall PersistentVolumeClaim Source
	pvcResource := asp.GetUninstallPersistentVolumeClaimSource()
	if err := r.uninstallPersistentVolumeClaim(instance, pvcResource); err != nil {
		log.V(-1).Info("retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}
	if !asp.SelfDriving {
		log.Info("selfDriving has been changed to false")
		selfDrivingResource := alamedaserviceparamter.GetSelfDrivingRsource()
		if err := r.uninstallScalerforAlameda(instance, selfDrivingResource); err != nil {
			log.V(-1).Info("retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
			return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
		}
	} else { //install Alameda Scaler
		if err := r.createScalerforAlameda(instance, asp, installResource); err != nil {
			log.V(-1).Info("create scaler for alameda failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
			return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
		}
	}

	if err = r.updateAlamedaService(instance, request.NamespacedName, asp); err != nil {
		log.Error(err, "Update AlamedaService failed, retry reconciling AlamedaService", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name, "msg", err.Error())
		return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
	}

	log.Info("Reconcile done.", "AlamedaService.Namespace", instance.Namespace, "AlamedaService.Name", instance.Name)
	return reconcile.Result{}, nil
}

func (r *ReconcileAlamedaService) InitAlamedaService(alamedaService *federatoraiv1alpha1.AlamedaService) {
	if alamedaService.Spec.EnableDispatcher == nil {
		enableTrue := true
		alamedaService.Spec.EnableDispatcher = &enableTrue
	}
}

func (r *ReconcileAlamedaService) getNamespace(namespaceName string) (corev1.Namespace, error) {
	namespace := corev1.Namespace{}
	if err := r.client.Get(context.TODO(), client.ObjectKey{Name: namespaceName}, &namespace); err != nil {
		return namespace, errors.Errorf("get namespace %s failed: %s", namespaceName, err.Error())
	}
	return namespace, nil
}

func (r *ReconcileAlamedaService) newComponentConfig(namespace corev1.Namespace) *component.ComponentConfig {
	podTemplateConfig := component.NewDefaultPodTemplateConfig(namespace)
	componentConfg := component.NewComponentConfig(namespace.Name, podTemplateConfig)
	return componentConfg
}

func (r *ReconcileAlamedaService) createScalerforAlameda(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter, resource *alamedaserviceparamter.Resource) error {
	for _, fileString := range resource.AlamedaScalerList {
		resourceScaler := componentConfig.NewAlamedaScaler(fileString)
		if err := controllerutil.SetControllerReference(instance, resourceScaler, r.scheme); err != nil {
			return errors.Errorf("Fail resourceScaler SetControllerReference: %s", err.Error())
		}
		foundScaler := &autoscaling_v1alpha1.AlamedaScaler{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: resourceScaler.Name, Namespace: resourceScaler.Namespace}, foundScaler)
		if err != nil && k8sErrors.IsNotFound(err) {
			log.Info("Creating a new Resource Scaler... ", "resourceScaler.Name", resourceScaler.Name)
			err = r.client.Create(context.TODO(), resourceScaler)
			if err != nil {
				return errors.Errorf("create Scaler %s/%s failed: %s", resourceScaler.Namespace, resourceScaler.Name, err.Error())
			}
			log.Info("Successfully Creating Resource Scaler", "resourceScaler.Name", resourceScaler.Name)
		} else if err != nil {
			return errors.Errorf("get Scaler %s/%s failed: %s", resourceScaler.Namespace, resourceScaler.Name, err.Error())
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) syncCustomResourceDefinition(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter, resource *alamedaserviceparamter.Resource) error {
	for _, fileString := range resource.CustomResourceDefinitionList {
		crd := componentConfig.RegistryCustomResourceDefinition(fileString)
		/*if err := controllerutil.SetControllerReference(instance, crd, r.scheme); err != nil {
			return errors.Errorf("Fail resourceCRB SetControllerReference: %s", err.Error())
		}*/
		_, err := resourceapply.ApplyCustomResourceDefinition(r.apiextclient.ApiextensionsV1beta1(), crd, asp)
		if err != nil {
			return errors.Wrapf(err, "syncCustomResourceDefinition faild: CustomResourceDefinition.Name: %s", crd.Name)
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) uninstallCustomResourceDefinition(resource *alamedaserviceparamter.Resource) {
	for _, fileString := range resource.CustomResourceDefinitionList {
		crd := componentConfig.RegistryCustomResourceDefinition(fileString)
		_, _, _ = resourceapply.DeleteCustomResourceDefinition(r.apiextclient.ApiextensionsV1beta1(), crd)
	}
}

func (r *ReconcileAlamedaService) syncClusterRoleBinding(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter, resource *alamedaserviceparamter.Resource) error {
	for _, FileStr := range resource.ClusterRoleBindingList {
		resourceCRB := componentConfig.NewClusterRoleBinding(FileStr)
		if err := controllerutil.SetControllerReference(instance, resourceCRB, r.scheme); err != nil {
			return errors.Errorf("Fail resourceCRB SetControllerReference: %s", err.Error())
		}
		foundCRB := &rbacv1.ClusterRoleBinding{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: resourceCRB.Name}, foundCRB)
		if err != nil && k8sErrors.IsNotFound(err) {
			log.Info("Creating a new Resource ClusterRoleBinding... ", "resourceCRB.Name", resourceCRB.Name)
			err = r.client.Create(context.TODO(), resourceCRB)
			if err != nil {
				return errors.Errorf("create clusterRoleBinding %s/%s failed: %s", resourceCRB.Namespace, resourceCRB.Name, err.Error())
			}
			log.Info("Successfully Creating Resource ClusterRoleBinding", "resourceCRB.Name", resourceCRB.Name)
		} else if err != nil {
			return errors.Errorf("get clusterRoleBinding %s/%s failed: %s", resourceCRB.Namespace, resourceCRB.Name, err.Error())
		} else {
			err = r.client.Update(context.TODO(), resourceCRB)
			if err != nil {
				return errors.Errorf("Update clusterRoleBinding %s/%s failed: %s", resourceCRB.Namespace, resourceCRB.Name, err.Error())
			}
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) updateNamespaceLabel(namespace string) error {

	labelsForCertManager := map[string]string{
		"certmanager.k8s.io/disable-validation": "true",
	}

	ctx := context.TODO()
	instance := corev1.Namespace{}
	if err := r.client.Get(ctx, types.NamespacedName{Name: namespace}, &instance); err != nil {
		return errors.Errorf("get Namesapce failed: %s", err.Error())
	}

	newInstance := *instance.DeepCopy()
	if newInstance.Labels == nil {
		newInstance.Labels = make(map[string]string)
	}
	for k, v := range labelsForCertManager {
		newInstance.Labels[k] = v
	}
	if err := r.client.Update(ctx, &newInstance); err != nil {
		return errors.Errorf("update Namespace failed: %s", err.Error())
	}

	return nil
}

func (r *ReconcileAlamedaService) createAlamedaNotificationChannels(instance *federatoraiv1alpha1.AlamedaService, resource *alamedaserviceparamter.Resource) error {
	for _, file := range resource.AlamedaNotificationChannelList {
		src, err := componentConfig.NewAlamedaNotificationChannel(file)
		if err != nil {
			return errors.Errorf("get AlamedaNotificationChannel failed: file: %s, error: %s", file, err.Error())
		}
		if err := controllerutil.SetControllerReference(instance, src, r.scheme); err != nil {
			return errors.Errorf("Fail AlamedaNotificationChannel SetControllerReference: %s", err.Error())
		}
		err = r.client.Create(context.TODO(), src)
		if err != nil && !k8sErrors.IsAlreadyExists(err) {
			return errors.Errorf("create AlamedaNotificationChannel %s failed: %s", src.GetName(), err.Error())
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) createAlamedaNotificationTopics(instance *federatoraiv1alpha1.AlamedaService, resource *alamedaserviceparamter.Resource) error {
	for _, file := range resource.AlamedaNotificationTopic {
		src, err := componentConfig.NewAlamedaNotificationTopic(file)
		if err != nil {
			return errors.Errorf("get AlamedaNotificationTopic failed: file: %s, error: %s", file, err.Error())
		}
		if err := controllerutil.SetControllerReference(instance, src, r.scheme); err != nil {
			return errors.Errorf("Fail AlamedaNotificationTopic SetControllerReference: %s", err.Error())
		}
		err = r.client.Create(context.TODO(), src)
		if err != nil && !k8sErrors.IsAlreadyExists(err) {
			return errors.Errorf("create AlamedaNotificationTopic %s failed: %s", src.GetName(), err.Error())
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) syncPodSecurityPolicy(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter, resource *alamedaserviceparamter.Resource) error {
	for _, FileStr := range resource.PodSecurityPolicyList {
		resourcePSP := componentConfig.NewPodSecurityPolicy(FileStr)
		if err := controllerutil.SetControllerReference(instance, resourcePSP, r.scheme); err != nil {
			return errors.Errorf("Fail resourcePSP SetControllerReference: %s", err.Error())
		}
		foundPSP := &v1beta1.PodSecurityPolicy{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: resourcePSP.Name}, foundPSP)
		if err != nil && k8sErrors.IsNotFound(err) {
			log.Info("Creating a new Resource PodSecurityPolicy... ", "resourcePSP.Name", resourcePSP.Name)
			err = r.client.Create(context.TODO(), resourcePSP)
			if err != nil {
				return errors.Errorf("create PodSecurityPolicy %s/%s failed: %s", resourcePSP.Namespace, resourcePSP.Name, err.Error())
			}
			log.Info("Successfully Creating Resource PodSecurityPolicy", "resourcePSP.Name", resourcePSP.Name)
		} else if err != nil {
			return errors.Errorf("get PodSecurityPolicy %s/%s failed: %s", resourcePSP.Namespace, resourcePSP.Name, err.Error())
		} else {
			err = r.client.Update(context.TODO(), resourcePSP)
			if err != nil {
				return errors.Errorf("Update PodSecurityPolicy %s/%s failed: %s", resourcePSP.Namespace, resourcePSP.Name, err.Error())
			}
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) syncSecurityContextConstraints(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter, resource *alamedaserviceparamter.Resource) error {
	for _, FileStr := range resource.SecurityContextConstraintsList {
		resourceSCC := componentConfig.NewSecurityContextConstraints(FileStr)
		if err := controllerutil.SetControllerReference(instance, resourceSCC, r.scheme); err != nil {
			return errors.Errorf("Fail resourceSCC SetControllerReference: %s", err.Error())
		}
		//process resource SecurityContextConstraints according to AlamedaService CR
		resourceSCC = processcrdspec.ParamterToSecurityContextConstraints(resourceSCC, asp)
		foundSCC := &securityv1.SecurityContextConstraints{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: resourceSCC.Name}, foundSCC)
		if err != nil && k8sErrors.IsNotFound(err) {
			log.Info("Creating a new Resource SecurityContextConstraints... ", "resourceSCC.Name", resourceSCC.Name)
			err = r.client.Create(context.TODO(), resourceSCC)
			if err != nil {
				return errors.Errorf("create SecurityContextConstraints %s/%s failed: %s", resourceSCC.Namespace, resourceSCC.Name, err.Error())
			}
			log.Info("Successfully Creating Resource SecurityContextConstraints", "resourceSCC.Name", resourceSCC.Name)
		} else if err != nil {
			return errors.Errorf("get SecurityContextConstraints %s/%s failed: %s", resourceSCC.Namespace, resourceSCC.Name, err.Error())
		} else {
			err = r.client.Update(context.TODO(), resourceSCC)
			if err != nil {
				return errors.Errorf("Update SecurityContextConstraints %s/%s failed: %s", resourceSCC.Namespace, resourceSCC.Name, err.Error())
			}
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) syncDaemonSet(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter, resource *alamedaserviceparamter.Resource) error {
	for _, FileStr := range resource.DaemonSetList {
		resourceDS := componentConfig.NewDaemonSet(FileStr)
		if err := controllerutil.SetControllerReference(instance, resourceDS, r.scheme); err != nil {
			return errors.Errorf("Fail resourceDS SetControllerReference: %s", err.Error())
		}
		//process resource DaemonSet according to AlamedaService CR
		resourceDS = processcrdspec.ParamterToDaemonSet(resourceDS, asp)
		foundDS := &appsv1.DaemonSet{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: resourceDS.Name, Namespace: resourceDS.Namespace}, foundDS)

		if err != nil && k8sErrors.IsNotFound(err) {
			log.Info("Creating a new Resource DaemonSet... ", "resourceDS.Namespace", resourceDS.Namespace, "resourceDS.Name", resourceDS.Name)
			err = r.client.Create(context.TODO(), resourceDS)
			if err != nil {
				return errors.Errorf("create DaemonSet %s/%s failed: %s", resourceDS.Namespace, resourceDS.Name, err.Error())
			}
			log.Info("Successfully Creating Resource DaemonSet", "resourceDS.Namespace", resourceDS.Namespace, "resourceDS.Name", resourceDS.Name)
		} else if err != nil {
			return errors.Errorf("get DaemonSet %s/%s failed: %s", resourceDS.Namespace, resourceDS.Name, err.Error())
		} else {
			if updateresource.MisMatchResourceDaemonSet(foundDS, resourceDS) {
				log.Info("Update Resource DaemonSet:", "foundDS.Name", foundDS.Name)
				err = r.client.Delete(context.TODO(), foundDS)
				if err != nil {
					return errors.Errorf("delete DaemonSet %s/%s failed: %s", foundDS.Namespace, foundDS.Name, err.Error())
				}
				err = r.client.Create(context.TODO(), resourceDS)
				if err != nil {
					return errors.Errorf("create DaemonSet %s/%s failed: %s", foundDS.Namespace, foundDS.Name, err.Error())
				}
				log.Info("Successfully Update Resource DaemonSet", "resourceDS.Namespace", resourceDS.Namespace, "resourceDS.Name", resourceDS.Name)
			}
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) syncClusterRole(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter, resource *alamedaserviceparamter.Resource) error {
	for _, FileStr := range resource.ClusterRoleList {
		resourceCR := componentConfig.NewClusterRole(FileStr)
		if err := controllerutil.SetControllerReference(instance, resourceCR, r.scheme); err != nil {
			return errors.Errorf("Fail resourceCR SetControllerReference: %s", err.Error())
		}
		foundCR := &rbacv1.ClusterRole{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: resourceCR.Name}, foundCR)
		if err != nil && k8sErrors.IsNotFound(err) {
			log.Info("Creating a new Resource ClusterRole... ", "resourceCR.Name", resourceCR.Name)
			err = r.client.Create(context.TODO(), resourceCR)
			if err != nil {
				return errors.Errorf("create clusterRole %s/%s failed: %s", resourceCR.Namespace, resourceCR.Name, err.Error())
			}
			log.Info("Successfully Creating Resource ClusterRole", "resourceCR.Name", resourceCR.Name)
		} else if err != nil {
			return errors.Errorf("get clusterRole %s/%s failed: %s", resourceCR.Namespace, resourceCR.Name, err.Error())
		} else {
			err = r.client.Update(context.TODO(), resourceCR)
			if err != nil {
				return errors.Errorf("Update clusterRole %s/%s failed: %s", resourceCR.Namespace, resourceCR.Name, err.Error())
			}
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) syncServiceAccount(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter, resource *alamedaserviceparamter.Resource) error {
	for _, FileStr := range resource.ServiceAccountList {
		resourceSA := componentConfig.NewServiceAccount(FileStr)
		if err := controllerutil.SetControllerReference(instance, resourceSA, r.scheme); err != nil {
			return errors.Errorf("Fail resourceSA SetControllerReference: %s", err.Error())
		}
		foundSA := &corev1.ServiceAccount{}

		err := r.client.Get(context.TODO(), types.NamespacedName{Name: resourceSA.Name, Namespace: resourceSA.Namespace}, foundSA)
		if err != nil && k8sErrors.IsNotFound(err) {
			log.Info("Creating a new Resource ServiceAccount... ", "resourceSA.Name", resourceSA.Name)
			err = r.client.Create(context.TODO(), resourceSA)
			if err != nil {
				return errors.Errorf("create serviceAccount %s/%s failed: %s", resourceSA.Namespace, resourceSA.Name, err.Error())
			}
			log.Info("Successfully Creating Resource ServiceAccount", "resourceSA.Name", resourceSA.Name)
		} else if err != nil {
			return errors.Errorf("get serviceAccount %s/%s failed: %s", resourceSA.Namespace, resourceSA.Name, err.Error())
		} else {
			err = r.client.Update(context.TODO(), resourceSA)
			if err != nil {
				return errors.Errorf("Update serviceAccount %s/%s failed: %s", resourceSA.Namespace, resourceSA.Name, err.Error())
			}
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) syncRole(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter, resource *alamedaserviceparamter.Resource) error {
	for _, FileStr := range resource.RoleList {
		resourceCR := componentConfig.NewRole(FileStr)
		if err := controllerutil.SetControllerReference(instance, resourceCR, r.scheme); err != nil {
			return errors.Errorf("Fail resourceCR SetControllerReference: %s", err.Error())
		}
		foundCR := &rbacv1.Role{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Namespace: resourceCR.Namespace, Name: resourceCR.Name}, foundCR)
		if err != nil && k8sErrors.IsNotFound(err) {
			log.Info("Creating a new Resource Role... ", "resourceCR.Namespace", resourceCR.Namespace, "resourceCR.Name", resourceCR.Name)
			err = r.client.Create(context.TODO(), resourceCR)
			if err != nil {
				return errors.Errorf("create Role %s/%s failed: %s", resourceCR.Namespace, resourceCR.Name, err.Error())
			}
			log.Info("Successfully Creating Resource Role", "resourceCR.Namespace", resourceCR.Namespace, "resourceCR.Name", resourceCR.Name)
		} else if err != nil {
			return errors.Errorf("get Role %s/%s failed: %s", resourceCR.Namespace, resourceCR.Name, err.Error())
		} else {
			err = r.client.Update(context.TODO(), resourceCR)
			if err != nil {
				return errors.Errorf("Update Role %s/%s failed: %s", resourceCR.Namespace, resourceCR.Name, err.Error())
			}
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) syncRoleBinding(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter, resource *alamedaserviceparamter.Resource) error {
	for _, FileStr := range resource.RoleBindingList {
		resourceCR := componentConfig.NewRoleBinding(FileStr)
		if err := controllerutil.SetControllerReference(instance, resourceCR, r.scheme); err != nil {
			return errors.Errorf("Fail resourceCR SetControllerReference: %s", err.Error())
		}
		foundCR := &rbacv1.RoleBinding{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Namespace: resourceCR.Namespace, Name: resourceCR.Name}, foundCR)
		if err != nil && k8sErrors.IsNotFound(err) {
			log.Info("Creating a new Resource RoleBinding... ", "resourceCR.Namespace", resourceCR.Namespace, "resourceCR.Name", resourceCR.Name)
			err = r.client.Create(context.TODO(), resourceCR)
			if err != nil {
				return errors.Errorf("create RoleBinding %s/%s failed: %s", resourceCR.Namespace, resourceCR.Name, err.Error())
			}
			log.Info("Successfully Creating Resource RoleBinding", "resourceCR.Namespace", resourceCR.Namespace, "resourceCR.Name", resourceCR.Name)
		} else if err != nil {
			return errors.Errorf("get RoleBinding %s/%s failed: %s", resourceCR.Namespace, resourceCR.Name, err.Error())
		} else {
			err = r.client.Update(context.TODO(), resourceCR)
			if err != nil {
				return errors.Errorf("Update RoleBinding %s/%s failed: %s", resourceCR.Namespace, resourceCR.Name, err.Error())
			}
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) syncIssuer(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter, resource *alamedaserviceparamter.Resource) error {
	for _, FileStr := range resource.IssuerList {
		issuer, err := componentConfig.NewIssuer(FileStr)
		if err != nil {
			return errors.Wrap(err, "new Issuer failed")
		} else if issuer == nil {
			return errors.Errorf("new Issuer failed, Issuer is nil")
		}
		if err := controllerutil.SetControllerReference(instance, issuer, r.scheme); err != nil {
			return errors.Errorf("Fail issuer SetControllerReference: %s", err.Error())
		}
		existIssuer := &certmanagerv1alpha1.Issuer{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Namespace: issuer.Namespace, Name: issuer.Name}, existIssuer)
		if err != nil && k8sErrors.IsNotFound(err) {
			log.Info("Creating a new Resource Issuer... ", "namespace", issuer.Namespace, "name", issuer.Name)
			err = r.client.Create(context.TODO(), issuer)
			if err != nil {
				return errors.Errorf("create Issuer %s/%s failed: %s", issuer.Namespace, issuer.Name, err.Error())
			}
			log.Info("Successfully Creating Resource Issuer", "namespace", issuer.Namespace, "name", issuer.Name)
		} else if err != nil {
			return errors.Errorf("get Issuer %s/%s failed: %s", issuer.Namespace, issuer.Name, err.Error())
		} else {
			newIssuer := updateresource.UpdateIssuer(*issuer, *existIssuer)
			err = r.client.Update(context.TODO(), &newIssuer)
			if err != nil {
				return errors.Errorf("update Issuer %s/%s failed: %s", issuer.Namespace, issuer.Name, err.Error())
			}
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) syncCertificate(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter, resource *alamedaserviceparamter.Resource) error {
	for _, FileStr := range resource.CertificateList {
		certificate, err := componentConfig.NewCertificate(FileStr)
		if err != nil {
			return errors.Wrap(err, "new Certificate failed")
		} else if certificate == nil {
			return errors.Errorf("new Certificate failed, Certificate is nil")
		}
		if err := controllerutil.SetControllerReference(instance, certificate, r.scheme); err != nil {
			return errors.Errorf("Fail certificate SetControllerReference: %s", err.Error())
		}
		existCertificate := &certmanagerv1alpha1.Certificate{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Namespace: certificate.Namespace, Name: certificate.Name}, existCertificate)
		if err != nil && k8sErrors.IsNotFound(err) {
			log.Info("Creating a new Resource Certificate... ", "namespace", certificate.Namespace, "name", certificate.Name)
			err = r.client.Create(context.TODO(), certificate)
			if err != nil {
				return errors.Errorf("create Certificate %s/%s failed: %s", certificate.Namespace, certificate.Name, err.Error())
			}
			log.Info("Successfully create Resource Certificate", "namespace", certificate.Namespace, "name", certificate.Name)
		} else if err != nil {
			return errors.Errorf("get Certificate %s/%s failed: %s", certificate.Namespace, certificate.Name, err.Error())
		} else {
			newCertificate := updateresource.UpdateCertificate(*certificate, *existCertificate)
			err = r.client.Update(context.TODO(), &newCertificate)
			if err != nil {
				return errors.Errorf("update Certificate %s/%s failed: %s", certificate.Namespace, certificate.Name, err.Error())
			}
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) createPersistentVolumeClaim(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter, resource *alamedaserviceparamter.Resource) error {
	for _, FileStr := range resource.PersistentVolumeClaimList {
		resourcePVC := componentConfig.NewPersistentVolumeClaim(FileStr)
		//process resource configmap into desire configmap
		resourcePVC = processcrdspec.ParamterToPersistentVolumeClaim(resourcePVC, asp)
		if err := controllerutil.SetControllerReference(instance, resourcePVC, r.scheme); err != nil {
			return errors.Errorf("Fail resourcePVC SetControllerReference: %s", err.Error())
		}
		foundPVC := &corev1.PersistentVolumeClaim{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: resourcePVC.Name, Namespace: resourcePVC.Namespace}, foundPVC)
		if err != nil && k8sErrors.IsNotFound(err) {
			log.Info("Creating a new Resource PersistentVolumeClaim... ", "resourcePVC.Name", resourcePVC.Name)
			err = r.client.Create(context.TODO(), resourcePVC)
			if err != nil {
				return errors.Errorf("create PersistentVolumeClaim %s/%s failed: %s", resourcePVC.Namespace, resourcePVC.Name, err.Error())
			}
			log.Info("Successfully Creating Resource PersistentVolumeClaim", "resourcePVC.Name", resourcePVC.Name)
		} else if err != nil {
			return errors.Errorf("get PersistentVolumeClaim %s/%s failed: %s", resourcePVC.Namespace, resourcePVC.Name, err.Error())
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) createSecret(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter, resource *alamedaserviceparamter.Resource) error {
	secret, err := componentConfig.NewAdmissionControllerSecret()
	if err != nil {
		return errors.Errorf("build AdmissionController secret failed: %s", err.Error())
	}
	if err := controllerutil.SetControllerReference(instance, secret, r.scheme); err != nil {
		return errors.Errorf("set controller reference to secret %s/%s failed: %s", secret.Namespace, secret.Name, err.Error())
	}
	err = r.client.Create(context.TODO(), secret)
	if err != nil && k8sErrors.IsAlreadyExists(err) {
		log.Info("create secret failed: secret is already exists", "secret.Namespace", secret.Namespace, "secret.Name", secret.Name)
	} else if err != nil {
		return errors.Errorf("get secret %s/%s failed: %s", secret.Namespace, secret.Name, err.Error())
	}
	secret, err = componentConfig.NewInfluxDBSecret()
	if err != nil {
		return errors.Errorf("build InfluxDB secret failed: %s", err.Error())
	}
	if err := controllerutil.SetControllerReference(instance, secret, r.scheme); err != nil {
		return errors.Errorf("set controller reference to secret %s/%s failed: %s", secret.Namespace, secret.Name, err.Error())
	}
	err = r.client.Create(context.TODO(), secret)
	if err != nil && k8sErrors.IsAlreadyExists(err) {
		log.Info("create secret failed: secret is already exists", "secret.Namespace", secret.Namespace, "secret.Name", secret.Name)
	} else if err != nil {
		return errors.Errorf("get secret %s/%s failed: %s", secret.Namespace, secret.Name, err.Error())
	}
	secret, err = componentConfig.NewfedemeterSecret()
	if err != nil {
		return errors.Errorf("build Fedemeter secret failed: %s", err.Error())
	}
	if err := controllerutil.SetControllerReference(instance, secret, r.scheme); err != nil {
		return errors.Errorf("set controller reference to secret %s/%s failed: %s", secret.Namespace, secret.Name, err.Error())
	}
	err = r.client.Create(context.TODO(), secret)
	if err != nil && k8sErrors.IsAlreadyExists(err) {
		log.Info("create secret failed: secret is already exists", "secret.Namespace", secret.Namespace, "secret.Name", secret.Name)
	} else if err != nil {
		return errors.Errorf("get secret %s/%s failed: %s", secret.Namespace, secret.Name, err.Error())
	}
	return nil
}

func (r *ReconcileAlamedaService) syncConfigMap(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter, resource *alamedaserviceparamter.Resource) error {
	for _, fileString := range resource.ConfigMapList {
		resourceCM := componentConfig.NewConfigMap(fileString)
		if err := controllerutil.SetControllerReference(instance, resourceCM, r.scheme); err != nil {
			return errors.Errorf("Fail resourceCM SetControllerReference: %s", err.Error())
		}
		//process resource configmap into desire configmap
		resourceCM = processcrdspec.ParamterToConfigMap(resourceCM, asp)
		foundCM := &corev1.ConfigMap{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: resourceCM.Name, Namespace: resourceCM.Namespace}, foundCM)
		if err != nil && k8sErrors.IsNotFound(err) {
			log.Info("Creating a new Resource ConfigMap... ", "resourceCM.Name", resourceCM.Name)
			err = r.client.Create(context.TODO(), resourceCM)
			if err != nil {
				return errors.Errorf("create configMap %s/%s failed: %s", resourceCM.Namespace, resourceCM.Name, err.Error())
			}
			log.Info("Successfully Creating Resource ConfigMap", "resourceCM.Name", resourceCM.Name)
		} else if err != nil {
			return errors.Errorf("get configMap %s/%s failed: %s", resourceCM.Namespace, resourceCM.Name, err.Error())
		} else {
			if updateresource.MisMatchResourceConfigMap(foundCM, resourceCM) {
				log.Info("Update Resource Service:", "foundCM.Name", foundCM.Name)
				err = r.client.Update(context.TODO(), foundCM)
				if err != nil {
					return errors.Errorf("update configMap %s/%s failed: %s", foundCM.Namespace, foundCM.Name, err.Error())
				} else {
					if foundCM.Name == util.GrafanaDatasourcesName { //if modify grafana-datasource then delete Deployment(Temporary strategy)
						grafanaDep := componentConfig.NewDeployment(util.GrafanaYaml)
						err = r.deleteDeploymentWhenModifyConfigMapOrService(grafanaDep)
						if err != nil {
							errors.Errorf("delete Deployment when modify ConfigMap %s/%s failed: %s", grafanaDep.Namespace, grafanaDep.Name, err.Error())
						}
					}
				}
				log.Info("Successfully Update Resource CinfigMap", "resourceCM.Name", foundCM.Name)
			}
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) syncService(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter, resource *alamedaserviceparamter.Resource) error {
	for _, fileString := range resource.ServiceList {
		resourceSV := componentConfig.NewService(fileString)
		if err := controllerutil.SetControllerReference(instance, resourceSV, r.scheme); err != nil {
			return errors.Errorf("Fail resourceSV SetControllerReference: %s", err.Error())
		}
		if err := processcrdspec.ParamterToService(resourceSV, asp); err != nil {
			return errors.Wrapf(err, "process service (%s/%s) failed", resourceSV.Namespace, resourceSV.Name)
		}
		foundSV := &corev1.Service{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: resourceSV.Name, Namespace: resourceSV.Namespace}, foundSV)
		if err != nil && k8sErrors.IsNotFound(err) {
			log.Info("Creating a new Resource Service... ", "resourceSV.Name", resourceSV.Name)
			err = r.client.Create(context.TODO(), resourceSV)
			if err != nil {
				return errors.Errorf("create service %s/%s failed: %s", resourceSV.Namespace, resourceSV.Name, err.Error())
			}
			log.Info("Successfully Creating Resource Service", "resourceSV.Name", resourceSV.Name)
		} else if err != nil {
			return errors.Errorf("get service %s/%s failed: %s", resourceSV.Namespace, resourceSV.Name, err.Error())
		} else {
			if updateresource.MisMatchResourceService(foundSV, resourceSV) {
				log.Info("Update Resource Service:", "foundSV.Name", foundSV.Name)
				err = r.client.Delete(context.TODO(), foundSV)
				if err != nil {
					return errors.Errorf("delete service %s/%s failed: %s", foundSV.Namespace, foundSV.Name, err.Error())
				}
				err = r.client.Create(context.TODO(), resourceSV)
				if err != nil {
					return errors.Errorf("create service %s/%s failed: %s", foundSV.Namespace, foundSV.Name, err.Error())
				}
				log.Info("Successfully Update Resource Service", "resourceSV.Name", foundSV.Name)
			}
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) syncMutatingWebhookConfiguration(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter, resource *alamedaserviceparamter.Resource) error {
	for _, fileString := range resource.MutatingWebhookConfigurationList {
		mutatingWebhookConfiguration, err := componentConfig.NewMutatingWebhookConfiguration(fileString)
		if err != nil {
			return errors.Wrap(err, "new MutatingWebhookConfiguration failed")
		}
		if err := controllerutil.SetControllerReference(instance, mutatingWebhookConfiguration, r.scheme); err != nil {
			return errors.Errorf("Fail MutatingWebhookConfiguration SetControllerReference: %s", err.Error())
		}
		existMutatingWebhookConfiguration := &admissionregistrationv1beta1.MutatingWebhookConfiguration{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: mutatingWebhookConfiguration.Name}, existMutatingWebhookConfiguration)
		if err != nil && k8sErrors.IsNotFound(err) {
			log.Info("Creating a new Resource MutatingWebhookConfiguration... ", "name", mutatingWebhookConfiguration.Name)
			err = r.client.Create(context.TODO(), mutatingWebhookConfiguration)
			if err != nil {
				return errors.Errorf("create MutatingWebhookConfiguration %s failed: %s", mutatingWebhookConfiguration.Name, err.Error())
			}
			log.Info("Successfully Creating Resource MutatingWebhookConfiguration", "name", mutatingWebhookConfiguration.Name)
		} else if err != nil {
			return errors.Errorf("get MutatingWebhookConfiguration %s failed: %s", mutatingWebhookConfiguration.Name, err.Error())
		} else {
			log.V(-1).Info("MutatingWebhookConfiguration is existing, skip creating", "name", mutatingWebhookConfiguration.Name)
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) syncValidatingWebhookConfiguration(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter, resource *alamedaserviceparamter.Resource) error {
	for _, fileString := range resource.ValidatingWebhookConfigurationList {
		validatingWebhookConfiguration, err := componentConfig.NewValidatingWebhookConfiguration(fileString)
		if err != nil {
			return errors.Wrap(err, "new ValidatingWebhookConfigurationList failed")
		}
		if err := controllerutil.SetControllerReference(instance, validatingWebhookConfiguration, r.scheme); err != nil {
			return errors.Errorf("Fail ValidatingWebhookConfiguration SetControllerReference: %s", err.Error())
		}
		existValidatingWebhookConfiguration := &admissionregistrationv1beta1.ValidatingWebhookConfiguration{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: validatingWebhookConfiguration.Name}, existValidatingWebhookConfiguration)
		if err != nil && k8sErrors.IsNotFound(err) {
			log.Info("Creating a new Resource ValidatingWebhookConfiguration... ", "name", validatingWebhookConfiguration.Name)
			err = r.client.Create(context.TODO(), validatingWebhookConfiguration)
			if err != nil {
				return errors.Errorf("create ValidatingWebhookConfiguration %s failed: %s", validatingWebhookConfiguration.Name, err.Error())
			}
			log.Info("Successfully Creating Resource ValidatingWebhookConfiguration", "name", validatingWebhookConfiguration.Name)
		} else if err != nil {
			return errors.Errorf("get ValidatingWebhookConfiguration %s failed: %s", validatingWebhookConfiguration.Name, err.Error())
		} else {
			log.V(-1).Info("ValidatingWebhookConfiguration is existing, skip creating", "name", validatingWebhookConfiguration.Name)
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) syncAPIService(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter, resource *alamedaserviceparamter.Resource) error {
	for _, fileString := range resource.APIServiceList {
		apiService, err := componentConfig.NewAPIService(fileString)
		if err != nil {
			return errors.Wrap(err, "new APIService failed")
		}
		if err := controllerutil.SetControllerReference(instance, apiService, r.scheme); err != nil {
			return errors.Errorf("Fail APIService SetControllerReference: %s", err.Error())
		}
		existAPIService := &apiregistrationv1beta1.APIService{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: apiService.Name}, existAPIService)
		if err != nil && k8sErrors.IsNotFound(err) {
			log.Info("Creating a new Resource APIService... ", "name", apiService.Name)
			err = r.client.Create(context.TODO(), apiService)
			if err != nil {
				return errors.Errorf("create APIService %s failed: %s", apiService.Name, err.Error())
			}
			log.Info("Successfully Creating Resource APIService", "name", apiService.Name)
		} else if err != nil {
			return errors.Errorf("get APIService %s failed: %s", apiService.Name, err.Error())
		} else {
			log.V(-1).Info("APIService is existing, skip creating", "name", apiService.Name)
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) syncDeployment(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter, resource *alamedaserviceparamter.Resource) error {
	for _, fileString := range resource.DeploymentList {
		resourceDep := componentConfig.NewDeployment(fileString)
		if err := controllerutil.SetControllerReference(instance, resourceDep, r.scheme); err != nil {
			return errors.Errorf("Fail resourceDep SetControllerReference: %s", err.Error())
		}
		//process resource deployment into desire deployment
		resourceDep = processcrdspec.ParamterToDeployment(resourceDep, asp)
		foundDep := &appsv1.Deployment{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: resourceDep.Name, Namespace: resourceDep.Namespace}, foundDep)
		if err != nil && k8sErrors.IsNotFound(err) {
			log.Info("Creating a new Resource Deployment... ", "resourceDep.Name", resourceDep.Name)
			err = r.client.Create(context.TODO(), resourceDep)
			if err != nil {
				return errors.Errorf("create deployment %s/%s failed: %s", resourceDep.Namespace, resourceDep.Name, err.Error())
			}
			log.Info("Successfully Creating Resource Deployment", "resourceDep.Name", resourceDep.Name)
			continue
		} else if err != nil {
			return errors.Errorf("get deployment %s/%s failed: %s", resourceDep.Namespace, resourceDep.Name, err.Error())
		} else {
			if updateresource.MisMatchResourceDeployment(foundDep, resourceDep) {
				log.Info("Update Resource Deployment:", "resourceDep.Name", foundDep.Name)
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

func (r *ReconcileAlamedaService) syncRoute(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter, resource *alamedaserviceparamter.Resource) error {
	for _, FileStr := range resource.RouteList {
		resourceRT := componentConfig.NewRoute(FileStr)
		if err := controllerutil.SetControllerReference(instance, resourceRT, r.scheme); err != nil {
			return errors.Errorf("Fail resourceRT SetControllerReference: %s", err.Error())
		}
		foundRT := &routev1.Route{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: resourceRT.Name, Namespace: resourceRT.Namespace}, foundRT)
		if err != nil && k8sErrors.IsNotFound(err) {
			log.Info("Creating a new Resource Route... ", "resourceRT.Name", resourceRT.Name)
			err = r.client.Create(context.TODO(), resourceRT)
			if err != nil {
				return errors.Errorf("create route %s/%s failed: %s", resourceRT.Namespace, resourceRT.Name, err.Error())
			}
			log.Info("Successfully Creating Resource route", "resourceRT.Name", resourceRT.Name)
		} else if err != nil {
			return errors.Errorf("get route %s/%s failed: %s", resourceRT.Namespace, resourceRT.Name, err.Error())
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) syncIngress(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter, resource *alamedaserviceparamter.Resource) error {
	for _, FileStr := range resource.IngressList {
		resourceIG := componentConfig.NewIngress(FileStr)
		if err := controllerutil.SetControllerReference(instance, resourceIG, r.scheme); err != nil {
			return errors.Errorf("Fail resourceIG SetControllerReference: %s", err.Error())
		}
		foundIG := &ingressv1beta1.Ingress{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: resourceIG.Name, Namespace: resourceIG.Namespace}, foundIG)
		if err != nil && k8sErrors.IsNotFound(err) {
			log.Info("Creating a new Resource Route... ", "resourceIG.Name", resourceIG.Name)
			err = r.client.Create(context.TODO(), resourceIG)
			if err != nil {
				return errors.Errorf("create route %s/%s failed: %s", resourceIG.Namespace, resourceIG.Name, err.Error())
			}
			log.Info("Successfully Creating Resource route", "resourceRT.Name", resourceIG.Name)
		} else if err != nil {
			return errors.Errorf("get route %s/%s failed: %s", resourceIG.Namespace, resourceIG.Name, err.Error())
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) syncStatefulSet(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter, resource *alamedaserviceparamter.Resource) error {
	for _, FileStr := range resource.StatefulSetList {
		resourceSS := componentConfig.NewStatefulSet(FileStr)
		if err := controllerutil.SetControllerReference(instance, resourceSS, r.scheme); err != nil {
			return errors.Errorf("Fail resourceSS SetControllerReference: %s", err.Error())
		}
		resourceSS = processcrdspec.ParamterToStatefulset(resourceSS, asp)
		foundSS := &appsv1.StatefulSet{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: resourceSS.Name, Namespace: resourceSS.Namespace}, foundSS)
		if err != nil && k8sErrors.IsNotFound(err) {
			log.Info("Creating a new Resource Route... ", "resourceSS.Name", resourceSS.Name)
			err = r.client.Create(context.TODO(), resourceSS)
			if err != nil {
				return errors.Errorf("create route %s/%s failed: %s", resourceSS.Namespace, resourceSS.Name, err.Error())
			}
			log.Info("Successfully Creating Resource route", "resourceSS.Name", resourceSS.Name)
		} else if err != nil {
			return errors.Errorf("get route %s/%s failed: %s", resourceSS.Namespace, resourceSS.Name, err.Error())
		} else {
			log.Info("Update Resource StatefulSet:", "resourceSS.Name", resourceSS.Name)
			if updateresource.MisMatchResourceStatefulSet(foundSS, resourceSS) {
				log.Info("Update Resource StatefulSet:", "name", foundSS.Name)
				err = r.client.Update(context.TODO(), foundSS)
				if err != nil {
					return errors.Errorf("update statefulSet %s/%s failed: %s", foundSS.Namespace, foundSS.Name, err.Error())
				}
				log.Info("Successfully Update Resource StatefulSet", "name", foundSS.Name)
			}
			log.Info("Updating Resource StatefulSet", "name", resourceSS.Name)
			err = r.client.Update(context.TODO(), resourceSS)
			if err != nil {
				return errors.Errorf("update StatefulSet %s/%s failed: %s", resourceSS.Namespace, resourceSS.Name, err.Error())
			}
			log.Info("Successfully Update Resource StatefulSet", "resourceSS.Name", resourceSS.Name)
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) uninstallStatefulSet(instance *federatoraiv1alpha1.AlamedaService, resource *alamedaserviceparamter.Resource) error {
	for _, fileString := range resource.StatefulSetList {
		resourceSS := componentConfig.NewStatefulSet(fileString)
		err := r.client.Delete(context.TODO(), resourceSS)
		if err != nil && k8sErrors.IsNotFound(err) {
			return nil
		} else if err != nil {
			return errors.Errorf("delete statefulset %s/%s failed: %s", resourceSS.Namespace, resourceSS.Name, err.Error())
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) uninstallIngress(instance *federatoraiv1alpha1.AlamedaService, resource *alamedaserviceparamter.Resource) error {
	for _, fileString := range resource.IngressList {
		resourceIG := componentConfig.NewIngress(fileString)
		err := r.client.Delete(context.TODO(), resourceIG)
		if err != nil && k8sErrors.IsNotFound(err) {
			return nil
		} else if err != nil {
			return errors.Errorf("delete ingress %s/%s failed: %s", resourceIG.Namespace, resourceIG.Name, err.Error())
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) uninstallRoute(instance *federatoraiv1alpha1.AlamedaService, resource *alamedaserviceparamter.Resource) error {
	for _, fileString := range resource.RouteList {
		resourceRT := componentConfig.NewRoute(fileString)
		err := r.client.Delete(context.TODO(), resourceRT)
		if err != nil && k8sErrors.IsNotFound(err) {
			return nil
		} else if err != nil {
			return errors.Errorf("delete route %s/%s failed: %s", resourceRT.Namespace, resourceRT.Name, err.Error())
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) uninstallDeployment(instance *federatoraiv1alpha1.AlamedaService, resource *alamedaserviceparamter.Resource) error {
	for _, fileString := range resource.DeploymentList {
		resourceDep := componentConfig.NewDeployment(fileString)
		err := r.client.Delete(context.TODO(), resourceDep)
		if err != nil && k8sErrors.IsNotFound(err) {
			return nil
		} else if err != nil {
			return errors.Errorf("delete deployment %s/%s failed: %s", resourceDep.Namespace, resourceDep.Name, err.Error())
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) uninstallService(instance *federatoraiv1alpha1.AlamedaService, resource *alamedaserviceparamter.Resource) error {
	for _, fileString := range resource.ServiceList {
		resourceSVC := componentConfig.NewService(fileString)
		err := r.client.Delete(context.TODO(), resourceSVC)
		if err != nil && k8sErrors.IsNotFound(err) {
			return nil
		} else if err != nil {
			return errors.Errorf("delete service %s/%s failed: %s", resourceSVC.Namespace, resourceSVC.Name, err.Error())
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) uninstallConfigMap(instance *federatoraiv1alpha1.AlamedaService, resource *alamedaserviceparamter.Resource) error {
	for _, fileString := range resource.ConfigMapList {
		resourceCM := componentConfig.NewConfigMap(fileString)
		err := r.client.Delete(context.TODO(), resourceCM)
		if err != nil && k8sErrors.IsNotFound(err) {
			return nil
		} else if err != nil {
			return errors.Errorf("delete comfigMap %s/%s failed: %s", resourceCM.Namespace, resourceCM.Name, err.Error())
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) uninstallSecret(instance *federatoraiv1alpha1.AlamedaService, resource *alamedaserviceparamter.Resource) error {
	for _, fileString := range resource.SecretList {
		resourceSec, _ := componentConfig.NewSecret(fileString)
		err := r.client.Delete(context.TODO(), resourceSec)
		if err != nil && k8sErrors.IsNotFound(err) {
			return nil
		} else if err != nil {
			return errors.Errorf("delete secret %s/%s failed: %s", resourceSec.Namespace, resourceSec.Name, err.Error())
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) uninstallServiceAccount(instance *federatoraiv1alpha1.AlamedaService, resource *alamedaserviceparamter.Resource) error {
	for _, fileString := range resource.ServiceAccountList {
		resourceSA := componentConfig.NewServiceAccount(fileString)
		err := r.client.Delete(context.TODO(), resourceSA)
		if err != nil && k8sErrors.IsNotFound(err) {
			return nil
		} else if err != nil {
			return errors.Errorf("delete serviceAccount %s/%s failed: %s", resourceSA.Namespace, resourceSA.Name, err.Error())
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) uninstallClusterRole(instance *federatoraiv1alpha1.AlamedaService, resource *alamedaserviceparamter.Resource) error {
	for _, fileString := range resource.ClusterRoleList {
		resourceCR := componentConfig.NewClusterRole(fileString)
		err := r.client.Delete(context.TODO(), resourceCR)
		if err != nil && k8sErrors.IsNotFound(err) {
			return nil
		} else if err != nil {
			return errors.Errorf("delete clusterRole %s/%s failed: %s", resourceCR.Namespace, resourceCR.Name, err.Error())
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) uninstallClusterRoleBinding(instance *federatoraiv1alpha1.AlamedaService, resource *alamedaserviceparamter.Resource) error {
	for _, fileString := range resource.ClusterRoleBindingList {
		resourceCRB := componentConfig.NewClusterRoleBinding(fileString)
		err := r.client.Delete(context.TODO(), resourceCRB)
		if err != nil && k8sErrors.IsNotFound(err) {
			return nil
		} else if err != nil {
			return errors.Errorf("delete clusterRoleBinding %s/%s failed: %s", resourceCRB.Namespace, resourceCRB.Name, err.Error())
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) uninstallAlamedaScaler(instance *federatoraiv1alpha1.AlamedaService, resource *alamedaserviceparamter.Resource) error {
	for _, fileString := range resource.AlamedaScalerList {
		resourceScaler := componentConfig.NewAlamedaScaler(fileString)
		err := r.client.Delete(context.TODO(), resourceScaler)
		if err != nil && k8sErrors.IsNotFound(err) {
			return nil
		} else if err != nil {
			return errors.Errorf("delete resourceScaler %s/%s failed: %s", resourceScaler.Namespace, resourceScaler.Name, err.Error())
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) uninstallMutatingWebhookConfiguration(instance *federatoraiv1alpha1.AlamedaService, resource *alamedaserviceparamter.Resource) error {
	for _, fileString := range resource.MutatingWebhookConfigurationList {
		mutatingWebhookConfiguration, err := componentConfig.NewMutatingWebhookConfiguration(fileString)
		if err != nil {
			return errors.Wrap(err, "new MutatingWebhookConfiguration failed")
		}
		err = r.client.Delete(context.TODO(), mutatingWebhookConfiguration)
		if err != nil && k8sErrors.IsNotFound(err) {
			return nil
		} else if err != nil {
			return errors.Errorf("delete MutatingWebhookConfiguratio %s failed: %s", mutatingWebhookConfiguration.Name, err.Error())
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) uninstallValidatingWebhookConfiguration(instance *federatoraiv1alpha1.AlamedaService, resource *alamedaserviceparamter.Resource) error {
	for _, fileString := range resource.AlamedaScalerList {
		validatingWebhookConfiguration, err := componentConfig.NewValidatingWebhookConfiguration(fileString)
		if err != nil {
			return errors.Wrap(err, "new ValidatingWebhookConfiguration failed")
		}
		err = r.client.Delete(context.TODO(), validatingWebhookConfiguration)
		if err != nil && k8sErrors.IsNotFound(err) {
			return nil
		} else if err != nil {
			return errors.Errorf("delete ValidatingWebhookConfiguration %s failed: %s", validatingWebhookConfiguration.Name, err.Error())
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) uninstallAPIService(instance *federatoraiv1alpha1.AlamedaService, resource *alamedaserviceparamter.Resource) error {
	for _, fileString := range resource.APIServiceList {
		apiService, err := componentConfig.NewAPIService(fileString)
		if err != nil {
			return errors.Wrap(err, "new PIService failed")
		}
		err = r.client.Delete(context.TODO(), apiService)
		if err != nil && k8sErrors.IsNotFound(err) {
			return nil
		} else if err != nil {
			return errors.Errorf("delete PIService %s failed: %s", apiService.Name, err.Error())
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) uninstallScalerforAlameda(instance *federatoraiv1alpha1.AlamedaService, resource *alamedaserviceparamter.Resource) error {
	if err := r.uninstallAlamedaScaler(instance, resource); err != nil {
		return errors.Wrapf(err, "uninstall selfDriving scaler failed")
	}
	return nil
}

func (r *ReconcileAlamedaService) uninstallGUIComponent(instance *federatoraiv1alpha1.AlamedaService, resource *alamedaserviceparamter.Resource) error {
	if err := r.uninstallRoute(instance, resource); err != nil {
		return errors.Wrapf(err, "uninstall gui component failed")
	}
	if err := r.uninstallDeployment(instance, resource); err != nil {
		return errors.Wrapf(err, "uninstall gui component failed")
	}
	if err := r.uninstallService(instance, resource); err != nil {
		return errors.Wrapf(err, "uninstall gui component failed")
	}
	if err := r.uninstallConfigMap(instance, resource); err != nil {
		return errors.Wrapf(err, "uninstall gui component failed")
	}
	if err := r.uninstallServiceAccount(instance, resource); err != nil {
		return errors.Wrapf(err, "uninstall gui component failed")
	}
	if err := r.uninstallClusterRole(instance, resource); err != nil {
		return errors.Wrapf(err, "uninstall gui component failed")
	}
	if err := r.uninstallClusterRoleBinding(instance, resource); err != nil {
		return errors.Wrapf(err, "uninstall gui component failed")
	}
	return nil
}

func (r *ReconcileAlamedaService) uninstallExecutionComponent(instance *federatoraiv1alpha1.AlamedaService, resource *alamedaserviceparamter.Resource) error {
	if err := r.uninstallDeployment(instance, resource); err != nil {
		return errors.Wrapf(err, "uninstall execution component failed")
	}
	if err := r.uninstallService(instance, resource); err != nil {
		return errors.Wrapf(err, "uninstall execution component failed")
	}
	if err := r.uninstallSecret(instance, resource); err != nil {
		return errors.Wrapf(err, "uninstall gui component failed")
	}
	if err := r.uninstallServiceAccount(instance, resource); err != nil {
		return errors.Wrapf(err, "uninstall gui component failed")
	}
	if err := r.uninstallClusterRole(instance, resource); err != nil {
		return errors.Wrapf(err, "uninstall gui component failed")
	}
	if err := r.uninstallClusterRoleBinding(instance, resource); err != nil {
		return errors.Wrapf(err, "uninstall gui component failed")
	}
	return nil
}

func (r *ReconcileAlamedaService) uninstallFedemeterComponent(instance *federatoraiv1alpha1.AlamedaService, resource *alamedaserviceparamter.Resource) error {
	if err := r.uninstallIngress(instance, resource); err != nil {
		return errors.Wrapf(err, "uninstall Fedemeter component failed")
	}
	if err := r.uninstallDeployment(instance, resource); err != nil {
		return errors.Wrapf(err, "uninstall Fedemeter component failed")
	}
	if err := r.uninstallService(instance, resource); err != nil {
		return errors.Wrapf(err, "uninstall Fedemeter component failed")
	}
	if err := r.uninstallSecret(instance, resource); err != nil {
		return errors.Wrapf(err, "uninstall Fedemeter component failed")
	}
	if err := r.uninstallConfigMap(instance, resource); err != nil {
		return errors.Wrapf(err, "uninstall Fedemeter component failed")
	}
	if err := r.uninstallStatefulSet(instance, resource); err != nil {
		return errors.Wrapf(err, "uninstall Fedemeter component failed")
	}
	return nil
}

func (r *ReconcileAlamedaService) uninstallPersistentVolumeClaim(instance *federatoraiv1alpha1.AlamedaService, resource *alamedaserviceparamter.Resource) error {
	for _, fileString := range resource.PersistentVolumeClaimList {
		resourcePVC := componentConfig.NewPersistentVolumeClaim(fileString)
		foundPVC := &corev1.PersistentVolumeClaim{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: resourcePVC.Name, Namespace: resourcePVC.Namespace}, foundPVC)
		if err != nil && k8sErrors.IsNotFound(err) {
			continue
		} else if err != nil {
			return errors.Errorf("get PersistentVolumeClaim %s/%s failed: %s", resourcePVC.Namespace, resourcePVC.Name, err.Error())
		} else {
			err := r.client.Delete(context.TODO(), resourcePVC)
			if err != nil && k8sErrors.IsNotFound(err) {
				return nil
			} else if err != nil {
				return errors.Errorf("delete PersistentVolumeClaim %s/%s failed: %s", resourcePVC.Namespace, resourcePVC.Name, err.Error())
			}
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) uninstallDaemonSet(instance *federatoraiv1alpha1.AlamedaService, resource *alamedaserviceparamter.Resource) error {
	for _, fileString := range resource.DaemonSetList {
		resourceDaemonSet := componentConfig.NewDaemonSet(fileString)
		foundDaemonSet := &appsv1.DaemonSet{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: resourceDaemonSet.Name, Namespace: resourceDaemonSet.Namespace}, foundDaemonSet)
		if err != nil && k8sErrors.IsNotFound(err) {
			continue
		} else if err != nil {
			return errors.Errorf("get DaemonSet %s/%s failed: %s", resourceDaemonSet.Namespace, resourceDaemonSet.Name, err.Error())
		} else {
			err := r.client.Delete(context.TODO(), resourceDaemonSet)
			if err != nil && k8sErrors.IsNotFound(err) {
				return nil
			} else if err != nil {
				return errors.Errorf("delete DaemonSet %s/%s failed: %s", resourceDaemonSet.Namespace, resourceDaemonSet.Name, err.Error())
			}
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) uninstallPodSecurityPolicy(instance *federatoraiv1alpha1.AlamedaService, resource *alamedaserviceparamter.Resource) error {
	for _, fileString := range resource.PodSecurityPolicyList {
		resourcePodSecurityPolicy := componentConfig.NewPodSecurityPolicy(fileString)
		foundPodSecurityPolicy := &v1beta1.PodSecurityPolicy{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: resourcePodSecurityPolicy.Name, Namespace: resourcePodSecurityPolicy.Namespace}, foundPodSecurityPolicy)
		if err != nil && k8sErrors.IsNotFound(err) {
			continue
		} else if err != nil {
			return errors.Errorf("get PodSecurityPolicy %s/%s failed: %s", resourcePodSecurityPolicy.Namespace, resourcePodSecurityPolicy.Name, err.Error())
		} else {
			err := r.client.Delete(context.TODO(), resourcePodSecurityPolicy)
			if err != nil && k8sErrors.IsNotFound(err) {
				return nil
			} else if err != nil {
				return errors.Errorf("delete PodSecurityPolicy %s/%s failed: %s", resourcePodSecurityPolicy.Namespace, resourcePodSecurityPolicy.Name, err.Error())
			}
		}
	}
	return nil
}

func (r *ReconcileAlamedaService) uninstallResource(resource alamedaserviceparamter.Resource) error {

	if err := r.uninstallStatefulSet(nil, &resource); err != nil {
		return errors.Wrap(err, "uninstall StatefulSet failed")
	}

	if err := r.uninstallIngress(nil, &resource); err != nil {
		return errors.Wrap(err, "uninstall Ingress failed")
	}

	if err := r.uninstallRoute(nil, &resource); err != nil {
		return errors.Wrap(err, "uninstall Route failed")
	}

	if err := r.uninstallDeployment(nil, &resource); err != nil {
		return errors.Wrap(err, "uninstall Deployment failed")
	}

	if err := r.uninstallService(nil, &resource); err != nil {
		return errors.Wrap(err, "uninstall Service failed")
	}

	if err := r.uninstallConfigMap(nil, &resource); err != nil {
		return errors.Wrap(err, "uninstall ConfigMap failed")
	}

	if err := r.uninstallSecret(nil, &resource); err != nil {
		return errors.Wrap(err, "uninstall Secret failed")
	}

	if err := r.uninstallServiceAccount(nil, &resource); err != nil {
		return errors.Wrap(err, "uninstall ServiceAccount failed")
	}

	if err := r.uninstallClusterRole(nil, &resource); err != nil {
		return errors.Wrap(err, "uninstall ClusterRole failed")
	}

	if err := r.uninstallClusterRoleBinding(nil, &resource); err != nil {
		return errors.Wrap(err, "uninstall ClusterRoleBinding failed")
	}

	if err := r.uninstallAlamedaScaler(nil, &resource); err != nil {
		return errors.Wrap(err, "uninstall AlamedaScaler failed")
	}

	if err := r.uninstallPersistentVolumeClaim(nil, &resource); err != nil {
		return errors.Wrap(err, "uninstall PersistentVolumeClaim failed")
	}

	if err := r.uninstallDaemonSet(nil, &resource); err != nil {
		return errors.Wrap(err, "uninstall DaemonSet failed")
	}

	if err := r.uninstallPodSecurityPolicy(nil, &resource); err != nil {
		return errors.Wrap(err, "uninstall PodSecurityPolicy failed")
	}

	if err := r.uninstallMutatingWebhookConfiguration(nil, &resource); err != nil {
		return errors.Wrap(err, "uninstall MutatingWebhookConfiguration failed")
	}

	if err := r.uninstallValidatingWebhookConfiguration(nil, &resource); err != nil {
		return errors.Wrap(err, "uninstall ValidatingWebhookConfiguration failed")
	}

	if err := r.uninstallAPIService(nil, &resource); err != nil {
		return errors.Wrap(err, "uninstall APIService failed")
	}

	return nil
}

func (r *ReconcileAlamedaService) needToReconcile(alamedaService *federatoraiv1alpha1.AlamedaService) (bool, error) {
	lock, lockErr := r.getAlamedaServiceLock(alamedaService.Namespace, alamedaServiceLockName)
	if lockErr == nil {
		if lockIsOwnedByAlamedaService(lock, alamedaService) {
			return true, nil
		}
	} else if k8sErrors.IsNotFound(lockErr) {
		err := r.createAlamedaServiceLock(alamedaService)
		if err == nil {
			return true, nil
		} else if !k8sErrors.IsAlreadyExists(err) {
			return false, errors.Wrap(err, "check if needs to reconcile failed")
		}
	} else if lockErr != nil {
		return false, errors.Wrap(lockErr, "check if needs to reconcile failed")
	}
	return false, nil
}

func (r *ReconcileAlamedaService) getAlamedaServiceLock(ns, name string) (rbacv1.ClusterRole, error) {
	lock := rbacv1.ClusterRole{}
	err := r.client.Get(context.Background(), types.NamespacedName{Name: name}, &lock)
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			return lock, err
		}
		return lock, errors.Errorf("get AlamedaService lock failed: %s", err.Error())
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
		}
		return errors.Errorf("create AlamedaService lock failed: %s", err.Error())
	}
	return nil
}

func (r *ReconcileAlamedaService) updateAlamedaServiceActivation(alamedaService *federatoraiv1alpha1.AlamedaService, active bool) error {
	copyAlamedaService := &federatoraiv1alpha1.AlamedaService{}
	r.client.Get(context.TODO(), client.ObjectKey{Namespace: alamedaService.Namespace, Name: alamedaService.Name}, copyAlamedaService)
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

func (r *ReconcileAlamedaService) updateAlamedaService(alamedaService *federatoraiv1alpha1.AlamedaService, namespaceName client.ObjectKey, asp *alamedaserviceparamter.AlamedaServiceParamter) error {
	if err := r.updateAlamedaServiceStatus(alamedaService, namespaceName, asp); err != nil {
		return err
	}
	if err := r.updateAlamedaServiceAnnotations(alamedaService, namespaceName); err != nil {
		return err
	}
	return nil
}

func (r *ReconcileAlamedaService) updateAlamedaServiceStatus(alamedaService *federatoraiv1alpha1.AlamedaService, namespaceName client.ObjectKey, asp *alamedaserviceparamter.AlamedaServiceParamter) error {
	copyAlamedaService := alamedaService.DeepCopy()
	if err := r.client.Get(context.TODO(), namespaceName, copyAlamedaService); err != nil {
		return errors.Errorf("get AlamedaService failed: %s", err.Error())
	}
	r.InitAlamedaService(copyAlamedaService)
	copyAlamedaService.Status.CRDVersion = asp.CurrentCRDVersion
	if err := r.client.Update(context.Background(), copyAlamedaService); err != nil {
		return errors.Errorf("update AlamedaService Status failed: %s", err.Error())
	}
	log.Info("Update AlamedaService Status Successfully", "resource.Name", copyAlamedaService.Name)
	return nil
}

func (r *ReconcileAlamedaService) updateAlamedaServiceAnnotations(alamedaService *federatoraiv1alpha1.AlamedaService, namespaceName client.ObjectKey) error {
	copyAlamedaService := alamedaService.DeepCopy()
	if err := r.client.Get(context.TODO(), namespaceName, copyAlamedaService); err != nil {
		return errors.Errorf("get AlamedaService failed: %s", err.Error())
	}
	r.InitAlamedaService(copyAlamedaService)
	jsonSpec, err := copyAlamedaService.GetSpecAnnotationWithoutKeycode()
	if err != nil {
		return errors.Errorf("get AlamedaService spec annotation without keycode failed: %s", err.Error())
	}
	if copyAlamedaService.Annotations != nil {
		copyAlamedaService.Annotations["previousAlamedaServiceSpec"] = jsonSpec
	} else {
		annotations := make(map[string]string)
		annotations["previousAlamedaServiceSpec"] = jsonSpec
		copyAlamedaService.Annotations = annotations
	}
	if err := r.client.Update(context.Background(), copyAlamedaService); err != nil {
		return errors.Errorf("update AlamedaService Annotations failed: %s", err.Error())
	}
	log.Info("Update AlamedaService Annotations Successfully", "resource.Name", copyAlamedaService.Name)
	return nil
}

func (r *ReconcileAlamedaService) checkAlamedaServiceSpecIsChange(alamedaService *federatoraiv1alpha1.AlamedaService, namespaceName client.ObjectKey) (bool, error) {
	jsonSpec, err := alamedaService.GetSpecAnnotationWithoutKeycode()
	if err != nil {
		return false, errors.Errorf("get AlamedaService spec annotation without keycode failed: %s", err.Error())
	}
	currentAlamedaServiceSpec := jsonSpec
	previousAlamedaServiceSpec := alamedaService.Annotations["previousAlamedaServiceSpec"]
	if currentAlamedaServiceSpec == previousAlamedaServiceSpec {
		return false, nil
	}
	return true, nil
}

func (r *ReconcileAlamedaService) deleteDeploymentWhenModifyConfigMapOrService(dep *appsv1.Deployment) error {
	err := r.client.Delete(context.TODO(), dep)
	if err != nil {
		return err
	}
	return nil
}
