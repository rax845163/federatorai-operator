package kubernetes

import (
	autoscaling_v1alpha1 "github.com/containers-ai/alameda/operator/pkg/apis/autoscaling/v1alpha1"
	openshift_apps "github.com/openshift/api/apps"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	apps_v1 "k8s.io/api/apps/v1"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var Scheme = runtime.NewScheme()

func init() {
	addToScheme(Scheme)
}

func addToScheme(scheme *runtime.Scheme) {
	core_v1.AddToScheme(scheme)
	apps_v1.AddToScheme(scheme)
	admissionregistrationv1beta1.AddToScheme(scheme)
	autoscaling_v1alpha1.AddToScheme(scheme)
	err := openshift_apps.Install(scheme)
	if err != nil {
		panic(err)
	}
}
