package alamedaservice

import (
	"context"
	"fmt"

	federatoraiv1alpha1 "github.com/containers-ai/federatorai-operator/pkg/apis/federatorai/v1alpha1"
	"github.com/containers-ai/federatorai-operator/pkg/component"
	"github.com/containers-ai/federatorai-operator/pkg/lib/resourceapply"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"

	"k8s.io/apimachinery/pkg/api/errors"
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

var (
	_           reconcile.Reconciler = &ReconcileAlamedaService{}
	log                              = logf.Log.WithName("controller_alamedaservice")
	name                             = "kroos-installnamespace"
	gracePeriod                      = int64(3)
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
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	fmt.Printf("AlamedaService.Spec: %v\n", instance.Spec)
	if instance.Spec.AlmedaInstallOrUninstall {
		//r.CreateNameSpace()
		r.RegisterTestsCRD()
		r.InstallClusterRoleBinding(instance)
		r.InstallClusterRole(instance)
		r.InstallServiceAccount(instance)
		r.InstallConfigMap(instance)
		r.InstallPersistentVolumeClaim(instance)
		r.InstallService(instance)
		r.InstallDeployment(instance)
	} else {
		r.UninstallDeployment(instance)
		r.UninstallService(instance)
		r.UninstallPersistentVolumeClaim(instance)
		r.UninstallConfigMap(instance)
		r.UninstallServiceAccount(instance)
		r.UninstallClusterRole(instance)
		r.UninstallClusterRoleBinding(instance)
		r.DeleteRegisterTestsCRD()
		//r.DeleteNameSpace()
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileAlamedaService) CreateNameSpace() {
	_, err := r.kubeClient.Core().Namespaces().Get(name, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating NameSpace", "NameSpace.Name", name)
		_, err = r.kubeClient.Core().Namespaces().Create(&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		})

		if err != nil {
			log.Error(err, "failed to create namespace with name", "NameSpace.Name", name)
		}
		log.Info("Successfully Creating NameSpace", "NameSpace.Name", name)
	} else {
		log.Info("Found NameSpace", "NameSpace.Name", name)
	}
}
func (r *ReconcileAlamedaService) RegisterTestsCRD() {
	file_location := [...]string{"../../manifests/TestCrds.yaml"} //"../../assets/CustomResourceDefinition/alamedarecommendationsCRD.yaml",
	//"../../assets/CustomResourceDefinition/alamedascalersCRD.yaml",

	for _, file_str := range file_location {
		crd := component.RegistryCustomResourceDefinition(file_str)
		_, _, _ = resourceapply.ApplyCustomResourceDefinition(r.apiextclient.ApiextensionsV1beta1(), crd)
	}
}
func (r *ReconcileAlamedaService) InstallClusterRoleBinding(instance *federatoraiv1alpha1.AlamedaService) {
	file_location := [...]string{"../../assets/ClusterRoleBinding/alameda-datahubCRB.yaml",
		"../../assets/ClusterRoleBinding/alameda-operatorCRB.yaml",
		"../../assets/ClusterRoleBinding/alameda-evictionerCRB.yaml",
		"../../assets/ClusterRoleBinding/admission-controllerCRB.yaml"}
	for _, file_str := range file_location {
		ComponentA_crb := component.NewClusterRoleBinding(file_str)
		if err := controllerutil.SetControllerReference(instance, ComponentA_crb, r.scheme); err != nil {
			log.Error(err, "Fail ResourceCRB SetControllerReference")
		}
		found_ComponentA_crb := &rbacv1.ClusterRoleBinding{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: ComponentA_crb.Name}, found_ComponentA_crb)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Creating a new Resource ClusterRoleBinding... ", "ResourceCRB.Name", ComponentA_crb.Name)
			err = r.client.Create(context.TODO(), ComponentA_crb)
			if err != nil {
				log.Error(err, "Fail Creating Resource ClusterRoleBinding", "ResourceCRB.Name", ComponentA_crb.Name)
			}
			log.Info("Successfully Creating Resource ClusterRoleBinding", "ResourceCRB.Name", ComponentA_crb.Name)
		} else if err != nil {
			log.Error(err, "Not Found Resource ClusterRoleBinding", "ResourceCRB.Name", ComponentA_crb.Name)
		}
	}
	log.Info("Install ClusterRoleBinding OK")
}

func (r *ReconcileAlamedaService) InstallClusterRole(instance *federatoraiv1alpha1.AlamedaService) {
	file_location := [...]string{"../../assets/ClusterRole/alameda-datahubCR.yaml",
		"../../assets/ClusterRole/alameda-operatorCR.yaml",
		"../../assets/ClusterRole/alameda-evictionerCR.yaml",
		"../../assets/ClusterRole/admission-controllerCR.yaml",
		"../../assets/ClusterRole/aggregate-alameda-admin-edit-alamedaCR.yaml"}
	for _, file_str := range file_location {
		ComponentA_cr := component.NewClusterRole(file_str)
		if err := controllerutil.SetControllerReference(instance, ComponentA_cr, r.scheme); err != nil {
			log.Error(err, "Fail ResourceCR SetControllerReference")
		}
		found_ComponentA_cr := &rbacv1.ClusterRole{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: ComponentA_cr.Name}, found_ComponentA_cr)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Creating a new Resource ClusterRole... ", "ResourceCR.Name", ComponentA_cr.Name)
			err = r.client.Create(context.TODO(), ComponentA_cr)
			if err != nil {
				log.Error(err, "Fail Creating Resource ClusterRole", "ResourceCR.Name", ComponentA_cr.Name)
			}
			log.Info("Successfully Creating Resource ClusterRole", "ResourceCR.Name", ComponentA_cr.Name)
		} else if err != nil {
			log.Error(err, "Not Found Resource ClusterRole", "ResourceCR.Name", ComponentA_cr.Name)
		}
	}
	log.Info("Install ClusterRole OK")
}

func (r *ReconcileAlamedaService) InstallServiceAccount(instance *federatoraiv1alpha1.AlamedaService) {
	file_location := [...]string{"../../assets/ServiceAccount/alameda-datahubSA.yaml",
		"../../assets/ServiceAccount/alameda-operatorSA.yaml",
		"../../assets/ServiceAccount/alameda-evictionerSA.yaml",
		"../../assets/ServiceAccount/admission-controllerSA.yaml"}
	for _, file_str := range file_location {
		ComponentA_sa := component.NewServiceAccount(file_str)
		if err := controllerutil.SetControllerReference(instance, ComponentA_sa, r.scheme); err != nil {
			log.Error(err, "Fail ResourceSA SetControllerReference")
		}
		found_ComponentA_sa := &corev1.ServiceAccount{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: ComponentA_sa.Name, Namespace: ComponentA_sa.Namespace}, found_ComponentA_sa)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Creating a new Resource ServiceAccount... ", "ResourceSA.Name", ComponentA_sa.Name)
			err = r.client.Create(context.TODO(), ComponentA_sa)
			if err != nil {
				log.Error(err, "Fail Creating Resource ServiceAccount", "ResourceSA.Name", ComponentA_sa.Name)
			}
			log.Info("Successfully Creating Resource ServiceAccount", "ResourceSA.Name", ComponentA_sa.Name)
		} else if err != nil {
			log.Error(err, "Not Found Resource ServiceAccount", "ResourceSA.Name", ComponentA_sa.Name)
		}
	}
	log.Info("Install ServiceAccount OK")
}

func (r *ReconcileAlamedaService) InstallConfigMap(instance *federatoraiv1alpha1.AlamedaService) {
	file_location := [...]string{"../../assets/ConfigMap/grafana-datasources.yaml"}
	for _, file_str := range file_location {
		ComponentA_cm := component.NewConfigMap(file_str)
		if err := controllerutil.SetControllerReference(instance, ComponentA_cm, r.scheme); err != nil {
			log.Error(err, "Fail ResourceCM SetControllerReference")
		}
		found_ComponentA_cm := &corev1.ConfigMap{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: ComponentA_cm.Name, Namespace: ComponentA_cm.Namespace}, found_ComponentA_cm)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Creating a new Resource ConfigMap... ", "ResourceCM.Name", ComponentA_cm.Name)
			err = r.client.Create(context.TODO(), ComponentA_cm)
			if err != nil {
				log.Error(err, "Fail Creating Resource ConfigMap", "ResourceCM.Name", ComponentA_cm.Name)
			}
			log.Info("Successfully Creating Resource ConfigMap", "ResourceCM.Name", ComponentA_cm.Name)
		} else if err != nil {
			log.Error(err, "Not Found Resource ConfigMap", "ResourceCM.Name", ComponentA_cm.Name)
		}
	}
	log.Info("Install ConfigMap OK")
}
func (r *ReconcileAlamedaService) InstallPersistentVolumeClaim(instance *federatoraiv1alpha1.AlamedaService) {
	file_location := [...]string{"../../assets/PersistentVolumeClaim/my-alamedainfluxdbPVC.yaml",
		"../../assets/PersistentVolumeClaim/my-alamedagrafanaPVC.yaml"}
	for _, file_str := range file_location {
		ComponentA_pvc := component.NewPersistentVolumeClaim(file_str)
		if err := controllerutil.SetControllerReference(instance, ComponentA_pvc, r.scheme); err != nil {
			log.Error(err, "Fail ResourcePVC SetControllerReference")
		}
		found_ComponentA_pvc := &corev1.PersistentVolumeClaim{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: ComponentA_pvc.Name, Namespace: ComponentA_pvc.Namespace}, found_ComponentA_pvc)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Creating a new Resource PersistentVolumeClaim... ", "ResourcePVC.Name", ComponentA_pvc.Name)
			err = r.client.Create(context.TODO(), ComponentA_pvc)
			if err != nil {
				log.Error(err, "Fail Creating Resource PersistentVolumeClaim", "ResourcePVC.Name", ComponentA_pvc.Name)
			}
			log.Info("Successfully Creating Resource PersistentVolumeClaim", "ResourcePVC.Name", ComponentA_pvc.Name)
		} else if err != nil {
			log.Error(err, "Not Found Resource PersistentVolumeClaim", "ResourcePVC.Name", ComponentA_pvc.Name)
		}
	}
	log.Info("Install PersistentVolumeClaim OK")
}

func (r *ReconcileAlamedaService) InstallService(instance *federatoraiv1alpha1.AlamedaService) {
	file_location := [...]string{"../../assets/Service/alameda-datahubSV.yaml",
		"../../assets/Service/admission-controllerSV.yaml",
		"../../assets/Service/alameda-influxdbSV.yaml",
		"../../assets/Service/alameda-grafanaSV.yaml"}
	for _, file_str := range file_location {
		ComponentA_sv := component.NewService(file_str)
		if err := controllerutil.SetControllerReference(instance, ComponentA_sv, r.scheme); err != nil {
			log.Error(err, "Fail ResourceSV SetControllerReference")

		}
		found_ComponentA_sv := &corev1.Service{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: ComponentA_sv.Name, Namespace: ComponentA_sv.Namespace}, found_ComponentA_sv)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Creating a new Resource Service... ", "ResourceSV.Name", ComponentA_sv.Name)
			err = r.client.Create(context.TODO(), ComponentA_sv)
			if err != nil {
				log.Error(err, "Fail Creating Resource Service", "ResourceSV.Name", ComponentA_sv.Name)
			}
			log.Info("Successfully Creating Resource Service", "ResourceSV.Name", ComponentA_sv.Name)
		} else if err != nil {
			log.Error(err, "Not Found Resource Service", "ResourceSV.Name", ComponentA_sv.Name)
		}
	}
	log.Info("Install Service OK")

}

func (r *ReconcileAlamedaService) InstallDeployment(instance *federatoraiv1alpha1.AlamedaService) {
	file_location := [...]string{"../../assets/Deployment/alameda-datahubDM.yaml",
		"../../assets/Deployment/alameda-operatorDM.yaml",
		"../../assets/Deployment/alameda-evictionerDM.yaml",
		"../../assets/Deployment/admission-controllerDM.yaml",
		"../../assets/Deployment/alameda-influxdbDM.yaml",
		"../../assets/Deployment/alameda-grafanaDM.yaml"}
	for _, file_str := range file_location {
		ComponentA_dep := component.NewDeployment(file_str)
		if err := controllerutil.SetControllerReference(instance, ComponentA_dep, r.scheme); err != nil {
			log.Error(err, "Fail ResourceDep SetControllerReference")

		}
		found_ComponentA_dep := &appsv1.Deployment{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: ComponentA_dep.Name, Namespace: ComponentA_dep.Namespace}, found_ComponentA_dep)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Creating a new Resource Deployment... ", "ResourceDep.Name", ComponentA_dep.Name)
			err = r.client.Create(context.TODO(), ComponentA_dep)
			if err != nil {
				log.Error(err, "Fail Creating Resource Deployment", "ResourceDep.Name", ComponentA_dep.Name)
			}
			log.Info("Successfully Creating Resource Deployment", "ResourceDep.Name", ComponentA_dep.Name)
		} else if err != nil {
			log.Error(err, "Not Found Resource Deployment", "ResourceDep.Name", ComponentA_dep.Name)
		}
	}
	log.Info("Install Deployment OK")
}
func (r *ReconcileAlamedaService) UninstallDeployment(instance *federatoraiv1alpha1.AlamedaService) {
	file_location := [...]string{"../../assets/Deployment/alameda-datahubDM.yaml",
		"../../assets/Deployment/alameda-operatorDM.yaml",
		"../../assets/Deployment/alameda-evictionerDM.yaml",
		"../../assets/Deployment/admission-controllerDM.yaml",
		"../../assets/Deployment/alameda-influxdbDM.yaml",
		"../../assets/Deployment/alameda-grafanaDM.yaml"}
	for _, file_str := range file_location {
		ComponentA_dep := component.NewDeployment(file_str)
		if err := controllerutil.SetControllerReference(instance, ComponentA_dep, r.scheme); err != nil {
			log.Error(err, "Fail ResourceDep SetControllerReference")
		}
		found_ComponentA_dep := &appsv1.Deployment{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: ComponentA_dep.Name, Namespace: ComponentA_dep.Namespace}, found_ComponentA_dep)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Cluster IsNotFound Resource Deployment", "ResourceDep.Name", ComponentA_dep.Name)
		} else if err != nil {
			log.Error(err, "Not Found Resource Deployment", "ResourceDep.Name", ComponentA_dep)
		} else {
			err = r.client.Delete(context.TODO(), ComponentA_dep)
			if err != nil {
				log.Error(err, "Fail Delete Resource Deployment", "ResourceDep.Name", ComponentA_dep)
			}
		}
	}
}
func (r *ReconcileAlamedaService) UninstallService(instance *federatoraiv1alpha1.AlamedaService) {
	file_location := [...]string{"../../assets/Service/alameda-datahubSV.yaml",
		"../../assets/Service/admission-controllerSV.yaml",
		"../../assets/Service/alameda-influxdbSV.yaml",
		"../../assets/Service/alameda-grafanaSV.yaml"}
	for _, file_str := range file_location {
		ComponentA_sv := component.NewService(file_str)
		if err := controllerutil.SetControllerReference(instance, ComponentA_sv, r.scheme); err != nil {
			log.Error(err, "Fail ResourceSV SetControllerReference")
		}
		found_ComponentA_sv := &corev1.Service{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: ComponentA_sv.Name, Namespace: ComponentA_sv.Namespace}, found_ComponentA_sv)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Cluster IsNotFound Resource Service", "ResourceSV.Name", ComponentA_sv.Name)
		} else if err != nil {
			log.Error(err, "Not Found Resource Service", "ResourceSV.Name", ComponentA_sv)
		} else {
			err = r.client.Delete(context.TODO(), ComponentA_sv)
			if err != nil {
				log.Error(err, "Fail Delete Resource Service", "ResourceSV.Name", ComponentA_sv)
			}
		}
	}
}
func (r *ReconcileAlamedaService) UninstallPersistentVolumeClaim(instance *federatoraiv1alpha1.AlamedaService) {
	file_location := [...]string{"../../assets/PersistentVolumeClaim/my-alamedainfluxdbPVC.yaml",
		"../../assets/PersistentVolumeClaim/my-alamedagrafanaPVC.yaml"}
	for _, file_str := range file_location {
		ComponentA_pvc := component.NewPersistentVolumeClaim(file_str)
		if err := controllerutil.SetControllerReference(instance, ComponentA_pvc, r.scheme); err != nil {
			log.Error(err, "Fail ResourcePVC SetControllerReference")
		}
		found_ComponentA_pvc := &corev1.PersistentVolumeClaim{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: ComponentA_pvc.Name, Namespace: ComponentA_pvc.Namespace}, found_ComponentA_pvc)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Cluster IsNotFound Resource PersistentVolumeClaim", "ResourcePVC.Name", ComponentA_pvc.Name)
		} else if err != nil {
			log.Error(err, "Not Found Resource PersistentVolumeClaim", "ResourcePVC.Name", ComponentA_pvc)
		} else {
			err = r.client.Delete(context.TODO(), ComponentA_pvc)
			if err != nil {
				log.Error(err, "Fail Delete Resource PersistentVolumeClaim", "ResourcePVC.Name", ComponentA_pvc)
			}
		}
	}
}
func (r *ReconcileAlamedaService) UninstallConfigMap(instance *federatoraiv1alpha1.AlamedaService) {
	file_location := [...]string{"../../assets/ConfigMap/grafana-datasources.yaml"}
	for _, file_str := range file_location {
		ComponentA_cm := component.NewConfigMap(file_str)
		if err := controllerutil.SetControllerReference(instance, ComponentA_cm, r.scheme); err != nil {
			log.Error(err, "Fail ResourceCM SetControllerReference")
		}
		found_ComponentA_cm := &corev1.ConfigMap{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: ComponentA_cm.Name, Namespace: ComponentA_cm.Namespace}, found_ComponentA_cm)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Cluster IsNotFound Resource ConfigMap", "ResourceCM.Name", ComponentA_cm.Name)
		} else if err != nil {
			log.Error(err, "Not Found Resource ConfigMap", "ResourceCM.Name", ComponentA_cm)
		} else {
			err = r.client.Delete(context.TODO(), ComponentA_cm)
			if err != nil {
				log.Error(err, "Fail Delete Resource ConfigMap", "ResourceCM.Name", ComponentA_cm)
			}
		}
	}
}
func (r *ReconcileAlamedaService) UninstallServiceAccount(instance *federatoraiv1alpha1.AlamedaService) {
	file_location := [...]string{"../../assets/ServiceAccount/alameda-datahubSA.yaml",
		"../../assets/ServiceAccount/alameda-operatorSA.yaml",
		"../../assets/ServiceAccount/alameda-evictionerSA.yaml",
		"../../assets/ServiceAccount/admission-controllerSA.yaml"}
	for _, file_str := range file_location {
		ComponentA_sa := component.NewServiceAccount(file_str)
		if err := controllerutil.SetControllerReference(instance, ComponentA_sa, r.scheme); err != nil {
			log.Error(err, "Fail ResourceSA SetControllerReference")
		}
		found_ComponentA_sa := &corev1.ServiceAccount{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: ComponentA_sa.Name, Namespace: ComponentA_sa.Namespace}, found_ComponentA_sa)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Cluster IsNotFound Resource ServiceAccount", "ResourceSA.Name", ComponentA_sa.Name)
		} else if err != nil {
			log.Error(err, "Not Found Resource ServiceAccount", "ResourceSA.Name", ComponentA_sa)
		} else {
			err = r.client.Delete(context.TODO(), ComponentA_sa)
			if err != nil {
				log.Error(err, "Fail Delete Resource ServiceAccount", "ResourceSA.Name", ComponentA_sa)
			}
		}
	}
}
func (r *ReconcileAlamedaService) UninstallClusterRole(instance *federatoraiv1alpha1.AlamedaService) {
	file_location := [...]string{"../../assets/ClusterRole/alameda-datahubCR.yaml",
		"../../assets/ClusterRole/alameda-operatorCR.yaml",
		"../../assets/ClusterRole/alameda-evictionerCR.yaml",
		"../../assets/ClusterRole/admission-controllerCR.yaml",
		"../../assets/ClusterRole/aggregate-alameda-admin-edit-alamedaCR.yaml"}
	for _, file_str := range file_location {
		ComponentA_cr := component.NewClusterRole(file_str)
		if err := controllerutil.SetControllerReference(instance, ComponentA_cr, r.scheme); err != nil {
			log.Error(err, "Fail ResourceCR SetControllerReference")
		}
		found_ComponentA_cr := &rbacv1.ClusterRole{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: ComponentA_cr.Name}, found_ComponentA_cr)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Cluster IsNotFound Resource ClusterRole", "ResourceCR.Name", ComponentA_cr.Name)
		} else if err != nil {
			log.Error(err, "Not Found Resource ClusterRole", "ResourceCR.Name", ComponentA_cr)
		} else {
			err = r.client.Delete(context.TODO(), ComponentA_cr)
			if err != nil {
				log.Error(err, "Fail Delete Resource ClusterRole", "ResourceCR.Name", ComponentA_cr)
			}
		}
	}
}
func (r *ReconcileAlamedaService) UninstallClusterRoleBinding(instance *federatoraiv1alpha1.AlamedaService) {
	file_location := [...]string{"../../assets/ClusterRoleBinding/alameda-datahubCRB.yaml",
		"../../assets/ClusterRoleBinding/alameda-operatorCRB.yaml",
		"../../assets/ClusterRoleBinding/alameda-evictionerCRB.yaml",
		"../../assets/ClusterRoleBinding/admission-controllerCRB.yaml"}
	for _, file_str := range file_location {
		ComponentA_crb := component.NewClusterRoleBinding(file_str)
		if err := controllerutil.SetControllerReference(instance, ComponentA_crb, r.scheme); err != nil {
			log.Error(err, "Fail ResourceCRB SetControllerReference")
		}
		found_ComponentA_crb := &rbacv1.ClusterRoleBinding{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: ComponentA_crb.Name}, found_ComponentA_crb)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Cluster IsNotFound Resource ClusterRoleBinding", "ResourceCRB.Name", ComponentA_crb.Name)
		} else if err != nil {
			log.Error(err, "Not Found Resource ClusterRoleBinding", "ResourceCRB.Name", ComponentA_crb)
		} else {
			err = r.client.Delete(context.TODO(), ComponentA_crb)
			if err != nil {
				log.Error(err, "Fail Delete Resource ClusterRoleBinding", "ResourceCRB.Name", ComponentA_crb)
			}
		}
	}
}
func (r *ReconcileAlamedaService) DeleteRegisterTestsCRD() {
	file_location := [...]string{"../../manifests/TestCrds.yaml"} //"../../assets/CustomResourceDefinition/alamedarecommendationsCRD.yaml",
	//"../../assets/CustomResourceDefinition/alamedascalersCRD.yaml",

	for _, file_str := range file_location {
		crd := component.RegistryCustomResourceDefinition(file_str)
		_, _, _ = resourceapply.DeleteCustomResourceDefinition(r.apiextclient.ApiextensionsV1beta1(), crd)
	}
}
func (r *ReconcileAlamedaService) DeleteNameSpace() {
	_, err := r.kubeClient.Core().Namespaces().Get(name, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		log.Info("Cluster IsNotFound Resource Namespaces", "NameSpace.Name", name)
	} else if err != nil {
		log.Error(err, "Not Found Resource Namespaces", "NameSpace.Name", name)
	} else {
		err = r.kubeClient.Core().Namespaces().Delete(name, &metav1.DeleteOptions{})
		if err != nil {
			log.Error(err, "Fail Delete Resource Namespaces", "NameSpace.Name", name)
		}
	}
}

// newDeploymentForCR returns a busybox Deployment with the same name/namespace as the cr
/*
func newDeploymentForCR(cr *federatoraiv1alpha1.AlamedaService) *appsv1.Deployment {
	labels := map[string]string{
		"app": cr.Name,
	}

	dep := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-deployment",
			Namespace: cr.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: "nginx",
						Name:  "nginx",

					}},
				},
			},
		},
	}
	return dep
}*/

func int32Ptr(i int32) *int32 { return &i }
