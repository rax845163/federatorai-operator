package alamedaservice

import (
	"context"

	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/alamedaserviceparamter"
	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/enable"
	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/updateparamter"

	federatoraiv1alpha1 "github.com/containers-ai/federatorai-operator/pkg/apis/federatorai/v1alpha1"
	"github.com/containers-ai/federatorai-operator/pkg/component"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"

	"k8s.io/apimachinery/pkg/api/errors"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

var (
	_   reconcile.Reconciler = &ReconcileAlamedaService{}
	log                      = logf.Log.WithName("controller_alamedaservice")
	//name                             = "kroos-installnamespace"
	gracePeriod = int64(3)
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

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
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			//r.UninstallDeployment(instance)
			//r.UninstallService(instance)
			//r.UninstallConfigMap(instance)
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	asp := alamedaserviceparamter.NewAlamedaServiceParamter(instance)
	r.InstallConfigMap(instance, asp)
	r.InstallService(instance, asp)
	r.InstallDeployment(instance, asp)
	if !asp.EnableExecution || !asp.EnableGUI { // if EnableExecution Or EnableGUI has been changed to false
		if !asp.EnableExecution {
			log.Info("EnableExecution has been changed to false")
			r.UninstallExcutionComponent(instance)
		}
		if !asp.EnableGUI {
			log.Info("EnableGUI has been changed to false")
			r.UninstallGUIComponent(instance)
		}
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileAlamedaService) InstallConfigMap(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter) {
	FileLocation := []string{"ConfigMap/grafana-datasources.yaml"}
	FileLocation = enable.IgnoreGUIYAML(FileLocation, asp.Guicomponent)
	FileLocation = enable.IgnoreExcutionYAML(FileLocation, asp.Excutioncomponent)
	for _, FileStr := range FileLocation {
		Resource_cm := component.NewConfigMap(FileStr)
		if err := controllerutil.SetControllerReference(instance, Resource_cm, r.scheme); err != nil {
			log.Error(err, "Fail ResourceCM SetControllerReference")
		}
		Resource_cm.Namespace = instance.Namespace
		found_cm := &corev1.ConfigMap{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: Resource_cm.Name, Namespace: Resource_cm.Namespace}, found_cm)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Creating a new Resource ConfigMap... ", "ResourceCM.Name", Resource_cm.Name)
			err = r.client.Create(context.TODO(), Resource_cm)
			if err != nil {
				log.Error(err, "Fail Creating Resource ConfigMap", "ResourceCM.Name", Resource_cm.Name)
			} else {
				log.Info("Successfully Creating Resource ConfigMap", "ResourceCM.Name", Resource_cm.Name)
			}
		} else if err != nil {
			log.Error(err, "Not Found Resource ConfigMap", "ResourceCM.Name", Resource_cm.Name)
		}
	}
}
func (r *ReconcileAlamedaService) InstallPersistentVolumeClaim(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter) {
	FileLocation := []string{"PersistentVolumeClaim/my-alamedainfluxdbPVC.yaml",
		"PersistentVolumeClaim/my-alamedagrafanaPVC.yaml"}
	FileLocation = enable.IgnoreGUIYAML(FileLocation, asp.Guicomponent)
	FileLocation = enable.IgnoreExcutionYAML(FileLocation, asp.Excutioncomponent)
	for _, FileStr := range FileLocation {
		Resource_pvc := component.NewPersistentVolumeClaim(FileStr)
		if err := controllerutil.SetControllerReference(instance, Resource_pvc, r.scheme); err != nil {
			log.Error(err, "Fail ResourcePVC SetControllerReference")
		}
		Resource_pvc.Namespace = instance.Namespace
		found_pvc := &corev1.PersistentVolumeClaim{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: Resource_pvc.Name, Namespace: Resource_pvc.Namespace}, found_pvc)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Creating a new Resource PersistentVolumeClaim... ", "ResourcePVC.Name", Resource_pvc.Name)
			err = r.client.Create(context.TODO(), Resource_pvc)
			if err != nil {
				log.Error(err, "Fail Creating Resource PersistentVolumeClaim", "ResourcePVC.Name", Resource_pvc.Name)
			} else {
				log.Info("Successfully Creating Resource PersistentVolumeClaim", "ResourcePVC.Name", Resource_pvc.Name)
			}
		} else if err != nil {
			log.Error(err, "Not Found Resource PersistentVolumeClaim", "ResourcePVC.Name", Resource_pvc.Name)
		}
	}
}

func (r *ReconcileAlamedaService) InstallService(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter) {
	FileLocation := []string{"Service/alameda-datahubSV.yaml",
		"Service/admission-controllerSV.yaml",
		"Service/alameda-influxdbSV.yaml",
		"Service/alameda-grafanaSV.yaml"}
	FileLocation = enable.IgnoreGUIYAML(FileLocation, asp.Guicomponent)
	FileLocation = enable.IgnoreExcutionYAML(FileLocation, asp.Excutioncomponent)
	for _, FileStr := range FileLocation {
		Resource_sv := component.NewService(FileStr)
		if err := controllerutil.SetControllerReference(instance, Resource_sv, r.scheme); err != nil {
			log.Error(err, "Fail ResourceSV SetControllerReference")
		}
		Resource_sv.Namespace = instance.Namespace
		found_sv := &corev1.Service{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: Resource_sv.Name, Namespace: Resource_sv.Namespace}, found_sv)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Creating a new Resource Service... ", "ResourceSV.Name", Resource_sv.Name)
			err = r.client.Create(context.TODO(), Resource_sv)
			if err != nil {
				log.Error(err, "Fail Creating Resource Service", "ResourceSV.Name", Resource_sv.Name)
			} else {
				log.Info("Successfully Creating Resource Service", "ResourceSV.Name", Resource_sv.Name)
			}
		} else if err != nil {
			log.Error(err, "Not Found Resource Service", "ResourceSV.Name", Resource_sv.Name)
		}
	}
}

func (r *ReconcileAlamedaService) InstallDeployment(instance *federatoraiv1alpha1.AlamedaService, asp *alamedaserviceparamter.AlamedaServiceParamter) {
	FileLocation := []string{"Deployment/alameda-datahubDM.yaml",
		"Deployment/alameda-operatorDM.yaml",
		"Deployment/alameda-evictionerDM.yaml",
		"Deployment/admission-controllerDM.yaml",
		"Deployment/alameda-influxdbDM.yaml",
		"Deployment/alameda-grafanaDM.yaml"}
	FileLocation = enable.IgnoreGUIYAML(FileLocation, asp.Guicomponent)
	FileLocation = enable.IgnoreExcutionYAML(FileLocation, asp.Excutioncomponent)
	for _, FileStr := range FileLocation {
		Resource_dep := component.NewDeployment(FileStr)
		if err := controllerutil.SetControllerReference(instance, Resource_dep, r.scheme); err != nil {
			log.Error(err, "Fail ResourceDep SetControllerReference")

		}
		Resource_dep.Namespace = instance.Namespace
		found_dep := &appsv1.Deployment{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: Resource_dep.Name, Namespace: Resource_dep.Namespace}, found_dep)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Creating a new Resource Deployment... ", "ResourceDep.Name", Resource_dep.Name)
			Resource_dep = updateparamter.ProcessImageVersion(Resource_dep, asp.Version)
			Resource_dep = updateparamter.ProcessPrometheusService(Resource_dep, asp.PrometheusService)
			err = r.client.Create(context.TODO(), Resource_dep)
			if err != nil {
				log.Error(err, "Fail Creating Resource Deployment", "ResourceDep.Name", Resource_dep.Name)
			} else {
				log.Info("Successfully Creating Resource Deployment", "ResourceDep.Name", Resource_dep.Name)
			}
		} else if err != nil {
			log.Error(err, "Not Found Resource Deployment", "ResourceDep.Name", Resource_dep.Name)
		} else {
			update := updateparamter.MatchAlamedaServiceParamter(found_dep, asp.Version, asp.PrometheusService)
			if update {
				log.Info("Update Resource Deployment:", "ResourceDep.Name", found_dep.Name)
				found_dep = updateparamter.ProcessImageVersion(found_dep, asp.Version)
				found_dep = updateparamter.ProcessPrometheusService(found_dep, asp.PrometheusService)
				err = r.client.Update(context.TODO(), found_dep)
				if err != nil {
					log.Error(err, "Fail Update Resource Deployment", "ResourceDep.Name", found_dep.Name)
				}
				log.Info("Successfully Update Resource Deployment", "ResourceDep.Name", found_dep.Name)
			}
		}
	}
}
func (r *ReconcileAlamedaService) UninstallDeployment(instance *federatoraiv1alpha1.AlamedaService) {
	FileLocation := [...]string{"Deployment/alameda-datahubDM.yaml",
		"Deployment/alameda-operatorDM.yaml",
		"Deployment/alameda-evictionerDM.yaml",
		"Deployment/admission-controllerDM.yaml",
		"Deployment/alameda-influxdbDM.yaml",
		"Deployment/alameda-grafanaDM.yaml"}
	for _, FileStr := range FileLocation {
		Resource_dep := component.NewDeployment(FileStr)
		if err := controllerutil.SetControllerReference(instance, Resource_dep, r.scheme); err != nil {
			log.Error(err, "Fail ResourceDep SetControllerReference")
		}
		Resource_dep.Namespace = instance.Namespace
		found_dep := &appsv1.Deployment{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: Resource_dep.Name, Namespace: Resource_dep.Namespace}, found_dep)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Cluster IsNotFound Resource Deployment", "ResourceDep.Name", Resource_dep.Name)
		} else if err != nil {
			log.Error(err, "Not Found Resource Deployment", "ResourceDep.Name", Resource_dep)
		} else {
			err = r.client.Delete(context.TODO(), found_dep)
			if err != nil {
				log.Error(err, "Fail Delete Resource Deployment", "ResourceDep.Name", found_dep)
			}
		}
	}
}
func (r *ReconcileAlamedaService) UninstallService(instance *federatoraiv1alpha1.AlamedaService) {
	FileLocation := [...]string{"Service/alameda-datahubSV.yaml",
		"Service/admission-controllerSV.yaml",
		"Service/alameda-influxdbSV.yaml",
		"Service/alameda-grafanaSV.yaml"}
	for _, FileStr := range FileLocation {
		Resource_sv := component.NewService(FileStr)
		if err := controllerutil.SetControllerReference(instance, Resource_sv, r.scheme); err != nil {
			log.Error(err, "Fail ResourceSV SetControllerReference")
		}
		Resource_sv.Namespace = instance.Namespace
		found_sv := &corev1.Service{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: Resource_sv.Name, Namespace: Resource_sv.Namespace}, found_sv)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Cluster IsNotFound Resource Service", "ResourceSV.Name", Resource_sv.Name)
		} else if err != nil {
			log.Error(err, "Not Found Resource Service", "ResourceSV.Name", Resource_sv)
		} else {
			err = r.client.Delete(context.TODO(), found_sv)
			if err != nil {
				log.Error(err, "Fail Delete Resource Service", "ResourceSV.Name", found_sv)
			}
		}
	}
}
func (r *ReconcileAlamedaService) UninstallPersistentVolumeClaim(instance *federatoraiv1alpha1.AlamedaService) {
	FileLocation := [...]string{"PersistentVolumeClaim/my-alamedainfluxdbPVC.yaml",
		"PersistentVolumeClaim/my-alamedagrafanaPVC.yaml"}
	for _, FileStr := range FileLocation {
		Resource_pvc := component.NewPersistentVolumeClaim(FileStr)
		if err := controllerutil.SetControllerReference(instance, Resource_pvc, r.scheme); err != nil {
			log.Error(err, "Fail ResourcePVC SetControllerReference")
		}
		Resource_pvc.Namespace = instance.Namespace
		found_pvc := &corev1.PersistentVolumeClaim{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: Resource_pvc.Name, Namespace: Resource_pvc.Namespace}, found_pvc)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Cluster IsNotFound Resource PersistentVolumeClaim", "ResourcePVC.Name", Resource_pvc.Name)
		} else if err != nil {
			log.Error(err, "Not Found Resource PersistentVolumeClaim", "ResourcePVC.Name", Resource_pvc)
		} else {
			err = r.client.Delete(context.TODO(), found_pvc)
			if err != nil {
				log.Error(err, "Fail Delete Resource PersistentVolumeClaim", "ResourcePVC.Name", found_pvc)
			}
		}
	}
}
func (r *ReconcileAlamedaService) UninstallConfigMap(instance *federatoraiv1alpha1.AlamedaService) {
	FileLocation := [...]string{"ConfigMap/grafana-datasources.yaml"}
	for _, FileStr := range FileLocation {
		Resource_cm := component.NewConfigMap(FileStr)
		if err := controllerutil.SetControllerReference(instance, Resource_cm, r.scheme); err != nil {
			log.Error(err, "Fail ResourceCM SetControllerReference")
		}
		Resource_cm.Namespace = instance.Namespace
		found_cm := &corev1.ConfigMap{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: Resource_cm.Name, Namespace: Resource_cm.Namespace}, found_cm)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Cluster IsNotFound Resource ConfigMap", "ResourceCM.Name", Resource_cm.Name)
		} else if err != nil {
			log.Error(err, "Not Found Resource ConfigMap", "ResourceCM.Name", Resource_cm)
		} else {
			err = r.client.Delete(context.TODO(), found_cm)
			if err != nil {
				log.Error(err, "Fail Delete Resource ConfigMap", "ResourceCM.Name", found_cm)
			}
		}
	}
}

func (r *ReconcileAlamedaService) UninstallGUIComponent(instance *federatoraiv1alpha1.AlamedaService) {
	Resource_dep := component.NewDeployment("Deployment/alameda-grafanaDM.yaml")
	if err := controllerutil.SetControllerReference(instance, Resource_dep, r.scheme); err != nil {
		log.Error(err, "Fail ResourceDep SetControllerReference")
	}
	Resource_dep.Namespace = instance.Namespace
	found_dep := &appsv1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: Resource_dep.Name, Namespace: Resource_dep.Namespace}, found_dep)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Cluster IsNotFound Resource Deployment", "ResourceDep.Name", Resource_dep.Name)
	} else if err != nil {
		log.Error(err, "Not Found Resource Deployment", "ResourceDep.Name", Resource_dep)
	} else {
		err = r.client.Delete(context.TODO(), found_dep)
		if err != nil {
			log.Error(err, "Fail Delete Resource Deployment", "ResourceDep.Name", found_dep)
		}
	}
	Resource_sv := component.NewService("Service/alameda-grafanaSV.yaml")
	if err := controllerutil.SetControllerReference(instance, Resource_sv, r.scheme); err != nil {
		log.Error(err, "Fail ResourceSV SetControllerReference")
	}
	Resource_sv.Namespace = instance.Namespace
	found_sv := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: Resource_sv.Name, Namespace: Resource_sv.Namespace}, found_sv)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Cluster IsNotFound Resource Service", "ResourceSV.Name", Resource_sv.Name)
	} else if err != nil {
		log.Error(err, "Not Found Resource Service", "ResourceSV.Name", Resource_sv)
	} else {
		err = r.client.Delete(context.TODO(), found_sv)
		if err != nil {
			log.Error(err, "Fail Delete Resource Service", "ResourceSV.Name", found_sv)
		}
	}
	Resource_cm := component.NewConfigMap("ConfigMap/grafana-datasources.yaml")
	if err := controllerutil.SetControllerReference(instance, Resource_cm, r.scheme); err != nil {
		log.Error(err, "Fail ResourceCM SetControllerReference")
	}
	Resource_cm.Namespace = instance.Namespace
	found_cm := &corev1.ConfigMap{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: Resource_cm.Name, Namespace: Resource_cm.Namespace}, found_cm)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Cluster IsNotFound Resource ConfigMap", "ResourceCM.Name", Resource_cm.Name)
	} else if err != nil {
		log.Error(err, "Not Found Resource ConfigMap", "ResourceCM.Name", Resource_cm)
	} else {
		err = r.client.Delete(context.TODO(), found_cm)
		if err != nil {
			log.Error(err, "Fail Delete Resource ConfigMap", "ResourceCM.Name", found_cm)
		}
	}
}
func (r *ReconcileAlamedaService) UninstallExcutionComponent(instance *federatoraiv1alpha1.AlamedaService) {
	FileLocation := [...]string{"Deployment/admission-controllerDM.yaml",
		"Deployment/alameda-evictionerDM.yaml"}
	for _, FileStr := range FileLocation {
		Resource_dep := component.NewDeployment(FileStr)
		if err := controllerutil.SetControllerReference(instance, Resource_dep, r.scheme); err != nil {
			log.Error(err, "Fail ResourceDep SetControllerReference")
		}
		Resource_dep.Namespace = instance.Namespace
		found_dep := &appsv1.Deployment{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: Resource_dep.Name, Namespace: Resource_dep.Namespace}, found_dep)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Cluster IsNotFound Resource Deployment", "ResourceDep.Name", Resource_dep.Name)
		} else if err != nil {
			log.Error(err, "Not Found Resource Deployment", "ResourceDep.Name", Resource_dep)
		} else {
			err = r.client.Delete(context.TODO(), found_dep)
			if err != nil {
				log.Error(err, "Fail Delete Resource Deployment", "ResourceDep.Name", found_dep)
			}
		}
	}
	Resource_sv := component.NewService("Service/admission-controllerSV.yaml")
	if err := controllerutil.SetControllerReference(instance, Resource_sv, r.scheme); err != nil {
		log.Error(err, "Fail ResourceSV SetControllerReference")
	}
	Resource_sv.Namespace = instance.Namespace
	found_sv := &corev1.Service{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: Resource_sv.Name, Namespace: Resource_sv.Namespace}, found_sv)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Cluster IsNotFound Resource Service", "ResourceSV.Name", Resource_sv.Name)
	} else if err != nil {
		log.Error(err, "Not Found Resource Service", "ResourceSV.Name", Resource_sv)
	} else {
		err = r.client.Delete(context.TODO(), found_sv)
		if err != nil {
			log.Error(err, "Fail Delete Resource Service", "ResourceSV.Name", found_sv)
		}
	}
}
