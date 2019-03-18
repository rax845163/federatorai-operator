package alamedaservice

import (
	"context"

	"github.com/sirupsen/logrus"

	federatoraiv1alpha1 "github.com/containers-ai/federatorai-operator/pkg/apis/federatorai/v1alpha1"
	"github.com/containers-ai/federatorai-operator/pkg/component"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_alamedaservice")

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
	return &ReconcileAlamedaService{client: mgr.GetClient(), scheme: mgr.GetScheme()}
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
	/*
		err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForObject{})
		if err != nil {
			return err
		}*/

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

	err = c.Watch(&source.Kind{Type: &rbacv1.ClusterRole{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &federatoraiv1alpha1.AlamedaService{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &rbacv1.ClusterRoleBinding{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &federatoraiv1alpha1.AlamedaService{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileAlamedaService{}

// ReconcileAlamedaService reconciles a AlamedaService object
type ReconcileAlamedaService struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
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
	//r.InstallComponentA(instance)
	// Define a new Pod object
	//pod := newPodForCR(instance)
	//Define a new Deployment object
	dep := newDeploymentForCR(instance)
	// Set AlamedaService instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, dep, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	found := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: dep.Name, Namespace: dep.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			return reconcile.Result{}, err
		}
		// Deployment created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}
	// Deployment already exists - don't requeue
	reqLogger.Info("Skip reconcile: Deployment already exists", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
	r.InstallClusterRoleBinding(instance)
	r.InstallClusterRole(instance)
	r.InstallServiceAccount(instance)
	r.InstallCinfigMap(instance)
	r.InstallPersistentVolumeClaim(instance)
	r.InstallService(instance)
	r.InstallDeployment(instance)
	return reconcile.Result{}, nil
}

func (r *ReconcileAlamedaService) InstallClusterRoleBinding(instance *federatoraiv1alpha1.AlamedaService) {
	file_location := [...]string{"../../assets/ClusterRoleBinding/alameda-datahubCRB.yaml",
		"../../assets/ClusterRoleBinding/alameda-operatorCRB.yaml",
		"../../assets/ClusterRoleBinding/alameda-evictionerCRB.yaml",
		"../../assets/ClusterRoleBinding/admission-controllerCRB.yaml"}
	for _, file_str := range file_location {
		ComponentA_crb := component.NewClusterRoleBinding(file_str)
		if err := controllerutil.SetControllerReference(instance, ComponentA_crb, r.scheme); err != nil {
			logrus.Fatalf("Fail ", file_str, " SetControllerReference")
		}
		found_ComponentA_crb := &rbacv1.ClusterRoleBinding{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: ComponentA_crb.Name}, found_ComponentA_crb)
		if err != nil && errors.IsNotFound(err) {
			logrus.Info("Creating a new ComponentA_crb... ", "ComponentA_crb.Namespace", ComponentA_crb.Namespace, "ComponentA_crb.Name", ComponentA_crb.Name)
			err = r.client.Create(context.TODO(), ComponentA_crb)
			if err != nil {
				logrus.Fatalf("Fail Creating ComponentA_crb: %v", err)
			}
			logrus.Info("successfully Creating ComponentA_crb")
		} else if err != nil {
			logrus.Fatalf("Not Found ComponentA_crb: %v", err)
		}
		logrus.Info("Found  ", ComponentA_crb.Name)
	}
	logrus.Info("Install ClusterRoleBinding OK...........................................................................................  ")
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
			logrus.Fatalf("Fail ComponentA_sa SetControllerReference")
		}
		found_ComponentA_cr := &rbacv1.ClusterRole{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: ComponentA_cr.Name}, found_ComponentA_cr)
		if err != nil && errors.IsNotFound(err) {
			logrus.Info("Creating a new ComponentA_cr... ", "ComponentA_cr.Namespace", ComponentA_cr.Namespace, "ComponentA_cr.Name", ComponentA_cr.Name)
			err = r.client.Create(context.TODO(), ComponentA_cr)
			if err != nil {
				logrus.Fatalf("Fail Creating ComponentA_cr: %v", err)
			}
			logrus.Info("successfully Creating ComponentA_cr")
		} else if err != nil {
			logrus.Fatalf("Not Found ComponentA_cr: %v", err)
		}
		logrus.Info("Found  ", ComponentA_cr.Name)
	}
	logrus.Info("Install ClusterRole OK...........................................................................................  ")
}

func (r *ReconcileAlamedaService) InstallServiceAccount(instance *federatoraiv1alpha1.AlamedaService) {
	file_location := [...]string{"../../assets/ServiceAccount/alameda-datahubSA.yaml",
		"../../assets/ServiceAccount/alameda-operatorSA.yaml",
		"../../assets/ServiceAccount/alameda-evictionerSA.yaml",
		"../../assets/ServiceAccount/admission-controllerSA.yaml"}
	for _, file_str := range file_location {
		ComponentA_sa := component.NewServiceAccount(file_str)
		if err := controllerutil.SetControllerReference(instance, ComponentA_sa, r.scheme); err != nil {
			logrus.Fatalf("Fail ComponentA_sa SetControllerReference")
		}
		found_ComponentA_sa := &corev1.ServiceAccount{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: ComponentA_sa.Name, Namespace: ComponentA_sa.Namespace}, found_ComponentA_sa)
		if err != nil && errors.IsNotFound(err) {
			logrus.Info("Creating a new ComponentA_sa... ", "ComponentA_sa.Namespace", ComponentA_sa.Namespace, "ComponentA_sa.Name", ComponentA_sa.Name)
			err = r.client.Create(context.TODO(), ComponentA_sa)
			if err != nil {
				logrus.Fatalf("Fail Creating ComponentA_sa: %v", err)
			}
			logrus.Info("successfully Creating ComponentA_sa")
		} else if err != nil {
			logrus.Fatalf("Not Found ComponentA_sa: %v", err)
		}
		logrus.Info("Found  ", ComponentA_sa.Name)
	}
	logrus.Info("Install ServiceAccount OK...........................................................................................  ")
}

func (r *ReconcileAlamedaService) InstallCinfigMap(instance *federatoraiv1alpha1.AlamedaService) {
	file_location := [...]string{"../../assets/ConfigMap/grafana-datasources.yaml"}
	for _, file_str := range file_location {
		ComponentA_cm := component.NewConfigMap(file_str)
		if err := controllerutil.SetControllerReference(instance, ComponentA_cm, r.scheme); err != nil {
			logrus.Fatalf("Fail ComponentA_sa SetControllerReference")
		}
		found_ComponentA_cm := &corev1.ConfigMap{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: ComponentA_cm.Name, Namespace: ComponentA_cm.Namespace}, found_ComponentA_cm)
		if err != nil && errors.IsNotFound(err) {
			logrus.Info("Creating a new ComponentA_cm... ", "ComponentA_cm.Namespace", ComponentA_cm.Namespace, "ComponentA_cm.Name", ComponentA_cm.Name)
			err = r.client.Create(context.TODO(), ComponentA_cm)
			if err != nil {
				logrus.Fatalf("Fail Creating ComponentA_cm: %v", err)
			}
			logrus.Info("successfully Creating ComponentA_cm")
		} else if err != nil {
			logrus.Fatalf("Not Found ComponentA_cm: %v", err)
		}
		logrus.Info("Found  ", ComponentA_cm.Name)
	}
	logrus.Info("Install ConfigMap OK...........................................................................................  ")
}
func (r *ReconcileAlamedaService) InstallPersistentVolumeClaim(instance *federatoraiv1alpha1.AlamedaService) {
	file_location := [...]string{"../../assets/PersistentVolumeClaim/my-alamedainfluxdbPVC.yaml",
		"../../assets/PersistentVolumeClaim/my-alamedagrafanaPVC.yaml"}
	for _, file_str := range file_location {
		ComponentA_pvc := component.NewPersistentVolumeClaim(file_str)
		if err := controllerutil.SetControllerReference(instance, ComponentA_pvc, r.scheme); err != nil {
			logrus.Fatalf("Fail ComponentA_sa SetControllerReference")
		}
		found_ComponentA_pvc := &corev1.PersistentVolumeClaim{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: ComponentA_pvc.Name, Namespace: ComponentA_pvc.Namespace}, found_ComponentA_pvc)
		if err != nil && errors.IsNotFound(err) {
			logrus.Info("Creating a new ComponentA_pvc... ", "ComponentA_pvc.Namespace", ComponentA_pvc.Namespace, "ComponentA_pvc.Name", ComponentA_pvc.Name)
			err = r.client.Create(context.TODO(), ComponentA_pvc)
			if err != nil {
				logrus.Fatalf("Fail Creating ComponentA_pvc: %v", err)
			}
			logrus.Info("successfully Creating ComponentA_pvc")
		} else if err != nil {
			logrus.Fatalf("Not Found ComponentA_pvc: %v", err)
		}
		logrus.Info("Found  ", ComponentA_pvc.Name)
	}
	logrus.Info("Install PersistentVolumeClaim OK...........................................................................................  ")
}

func (r *ReconcileAlamedaService) InstallService(instance *federatoraiv1alpha1.AlamedaService) {
	file_location := [...]string{"../../assets/Service/alameda-datahubSV.yaml",
		"../../assets/Service/admission-controllerSV.yaml",
		"../../assets/Service/alameda-influxdbSV.yaml",
		"../../assets/Service/alameda-grafanaSV.yaml"}
	for _, file_str := range file_location {
		ComponentA_sv := component.NewService(file_str)
		if err := controllerutil.SetControllerReference(instance, ComponentA_sv, r.scheme); err != nil {
			logrus.Fatalf("Fail ComponentA_sa SetControllerReference")
		}
		found_ComponentA_sv := &corev1.Service{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: ComponentA_sv.Name, Namespace: ComponentA_sv.Namespace}, found_ComponentA_sv)
		if err != nil && errors.IsNotFound(err) {
			logrus.Info("Creating a new ComponentA_sv... ", "ComponentA_sv.Namespace", ComponentA_sv.Namespace, "ComponentA_sv.Name", ComponentA_sv.Name)
			err = r.client.Create(context.TODO(), ComponentA_sv)
			if err != nil {
				logrus.Fatalf("Fail Creating ComponentA_sv: %v", err)
			}
			logrus.Info("successfully Creating ComponentA_sv")
		} else if err != nil {
			logrus.Fatalf("Not Found ComponentA_sv: %v", err)
		}
		logrus.Info("Found  ", ComponentA_sv.Name)
	}
	logrus.Info("Install Service OK...........................................................................................  ")
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
			logrus.Fatalf("Fail ComponentA_dep SetControllerReference")
		}
		found_ComponentA_dep := &appsv1.Deployment{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: ComponentA_dep.Name, Namespace: ComponentA_dep.Namespace}, found_ComponentA_dep)
		if err != nil && errors.IsNotFound(err) {
			logrus.Info("Creating a new ComponentA_dep... ", "ComponentA_dep.Namespace", ComponentA_dep.Namespace, "ComponentA_dep.Name", ComponentA_dep.Name)
			err = r.client.Create(context.TODO(), ComponentA_dep)
			if err != nil {
				logrus.Fatalf("Fail Creating ComponentA_dep: %v", err)
			}
			logrus.Info("successfully Creating ComponentA_dep")
		} else if err != nil {
			logrus.Fatalf("Not Found ComponentA_dep: %v", err)
		}
		logrus.Info("Found  ", ComponentA_dep.Name)
	}
	logrus.Info("Install Deployment OK...........................................................................................  ")
}

/*
	"../../assets/Deployment/alameda-datahubDM.yaml":                       alamedadatahubDMYaml,
	"../../assets/Deployment/alameda-operatorDM.yaml":                      alamedaoperatorDMYaml,
	"../../assets/Deployment/alameda-evictionerDM.yaml":                    alamedaevictionerDMYaml,
	"../../assets/Deployment/admission-controllerDM.yaml":                  admissioncontrollerDMYaml,
	"../../assets/Deployment/alameda-influxdbDM.yaml":                      alamedainfluxdbDMYaml,
	"../../assets/Deployment/alameda-grafanaDM.yaml":                       alamedagrafanaDMYaml,
*/
/*
func (r *ReconcileAlamedaService) InstallComponentA(instance *federatoraiv1alpha1.AlamedaService) {
	ComponentA_dep := component.NewComponentADeployment()
	if err := controllerutil.SetControllerReference(instance, ComponentA_dep, r.scheme); err != nil {
		logrus.Fatalf("Fail ComponentA_dep SetControllerReference")
	}
	found_ComponentA_dep := &appsv1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: ComponentA_dep.Name, Namespace: ComponentA_dep.Namespace}, found_ComponentA_dep)
	if err != nil && errors.IsNotFound(err) {
		logrus.Info("Creating a new ComponentA_dep... ", "ComponentA_dep.Namespace", ComponentA_dep.Namespace, "ComponentA_dep.Name", ComponentA_dep.Name)
		err = r.client.Create(context.TODO(), ComponentA_dep)
		if err != nil {
			logrus.Fatalf("Fail Creating ComponentA_dep: %v", err)
		}
		logrus.Info("successfully Creating ComponentA_dep")
	} else if err != nil {
		logrus.Fatalf("Not Found ComponentA_dep: %v", err)
	}
	logrus.Info("Found  ", ComponentA_dep.Name)

	ComponentA_sa := component.NewComponentAServiceAccount()
	if err := controllerutil.SetControllerReference(instance, ComponentA_sa, r.scheme); err != nil {
		logrus.Fatalf("Fail ComponentA_sa SetControllerReference")
	}
	found_ComponentA_sa := &corev1.ServiceAccount{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: ComponentA_sa.Name, Namespace: ComponentA_sa.Namespace}, found_ComponentA_sa)
	if err != nil && errors.IsNotFound(err) {
		logrus.Info("Creating a new ComponentA_sa... ", "ComponentA_sa.Namespace", ComponentA_sa.Namespace, "ComponentA_sa.Name", ComponentA_sa.Name)
		err = r.client.Create(context.TODO(), ComponentA_sa)
		if err != nil {
			logrus.Fatalf("Fail Creating ComponentA_sa: %v", err)
		}
		logrus.Info("successfully Creating ComponentA_sa")
	} else if err != nil {
		logrus.Fatalf("Not Found ComponentA_sa: %v", err)
	}
	logrus.Info("Found  ", ComponentA_sa.Name)

	ComponentA_cm := component.NewComponentAConfigMap()
	if err := controllerutil.SetControllerReference(instance, ComponentA_cm, r.scheme); err != nil {
		logrus.Fatalf("Fail ComponentA_sa SetControllerReference")
	}
	found_ComponentA_cm := &corev1.ConfigMap{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: ComponentA_cm.Name, Namespace: ComponentA_cm.Namespace}, found_ComponentA_cm)
	if err != nil && errors.IsNotFound(err) {
		logrus.Info("Creating a new ComponentA_cm... ", "ComponentA_cm.Namespace", ComponentA_cm.Namespace, "ComponentA_cm.Name", ComponentA_cm.Name)
		err = r.client.Create(context.TODO(), ComponentA_cm)
		if err != nil {
			logrus.Fatalf("Fail Creating ComponentA_cm: %v", err)
		}
		logrus.Info("successfully Creating ComponentA_cm")
	} else if err != nil {
		logrus.Fatalf("Not Found ComponentA_cm: %v", err)
	}
	logrus.Info("Found  ", ComponentA_cm.Name)

	ComponentA_sv := component.NewComponentAService()
	if err := controllerutil.SetControllerReference(instance, ComponentA_sv, r.scheme); err != nil {
		logrus.Fatalf("Fail ComponentA_sa SetControllerReference")
	}
	found_ComponentA_sv := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: ComponentA_sv.Name, Namespace: ComponentA_sv.Namespace}, found_ComponentA_sv)
	if err != nil && errors.IsNotFound(err) {
		logrus.Info("Creating a new ComponentA_sv... ", "ComponentA_sv.Namespace", ComponentA_sv.Namespace, "ComponentA_sv.Name", ComponentA_sv.Name)
		err = r.client.Create(context.TODO(), ComponentA_sv)
		if err != nil {
			logrus.Fatalf("Fail Creating ComponentA_sv: %v", err)
		}
		logrus.Info("successfully Creating ComponentA_sv")
	} else if err != nil {
		logrus.Fatalf("Not Found ComponentA_sv: %v", err)
	}
	logrus.Info("Found  ", ComponentA_sv.Name)

	ComponentA_pvc := component.NewComponentAPersistentVolumeClaim()
	if err := controllerutil.SetControllerReference(instance, ComponentA_pvc, r.scheme); err != nil {
		logrus.Fatalf("Fail ComponentA_sa SetControllerReference")
	}
	found_ComponentA_pvc := &corev1.PersistentVolumeClaim{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: ComponentA_pvc.Name, Namespace: ComponentA_pvc.Namespace}, found_ComponentA_pvc)
	if err != nil && errors.IsNotFound(err) {
		logrus.Info("Creating a new ComponentA_pvc... ", "ComponentA_pvc.Namespace", ComponentA_pvc.Namespace, "ComponentA_pvc.Name", ComponentA_pvc.Name)
		err = r.client.Create(context.TODO(), ComponentA_pvc)
		if err != nil {
			logrus.Fatalf("Fail Creating ComponentA_pvc: %v", err)
		}
		logrus.Info("successfully Creating ComponentA_pvc")
	} else if err != nil {
		logrus.Fatalf("Not Found ComponentA_pvc: %v", err)
	}
	logrus.Info("Found  ", ComponentA_pvc.Name)
	fmt.Println(".ComponentA_cr := component.NewComponentAClusterRole()")
	ComponentA_cr := component.NewComponentAClusterRole()
	if err := controllerutil.SetControllerReference(instance, ComponentA_cr, r.scheme); err != nil {
		logrus.Fatalf("Fail ComponentA_sa SetControllerReference")
	}
	found_ComponentA_cr := &rbacv1.ClusterRole{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: ComponentA_cr.Name}, found_ComponentA_cr)
	if err != nil && errors.IsNotFound(err) {
		logrus.Info("Creating a new ComponentA_cr... ", "ComponentA_cr.Namespace", ComponentA_cr.Namespace, "ComponentA_cr.Name", ComponentA_cr.Name)
		err = r.client.Create(context.TODO(), ComponentA_cr)
		if err != nil {
			logrus.Fatalf("Fail Creating ComponentA_cr: %v", err)
		}
		logrus.Info("successfully Creating ComponentA_cr")
	} else if err != nil {
		logrus.Fatalf("Not Found ComponentA_cr: %v", err)
	}
	logrus.Info("Found  ", ComponentA_cr.Name)

	ComponentA_crb := component.NewComponentAClusterRoleBinding()
	if err := controllerutil.SetControllerReference(instance, ComponentA_crb, r.scheme); err != nil {
		logrus.Fatalf("Fail ComponentA_sa SetControllerReference")
	}
	found_ComponentA_crb := &rbacv1.ClusterRoleBinding{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: ComponentA_crb.Name}, found_ComponentA_crb)
	if err != nil && errors.IsNotFound(err) {
		logrus.Info("Creating a new ComponentA_crb... ", "ComponentA_crb.Namespace", ComponentA_crb.Namespace, "ComponentA_crb.Name", ComponentA_crb.Name)
		err = r.client.Create(context.TODO(), ComponentA_crb)
		if err != nil {
			logrus.Fatalf("Fail Creating ComponentA_crb: %v", err)
		}
		logrus.Info("successfully Creating ComponentA_crb")
	} else if err != nil {
		logrus.Fatalf("Not Found ComponentA_crb: %v", err)
	}
	logrus.Info("Found  ", ComponentA_crb.Name)

}
*/
// newDeploymentForCR returns a busybox Deployment with the same name/namespace as the cr
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
						Image: "busybox",
						Name:  "busybox",
					}},
				},
			},
		},
	}
	return dep
}
func int32Ptr(i int32) *int32 { return &i }
