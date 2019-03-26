package resourceread

import (
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var (
	appsScheme = runtime.NewScheme()
	appsCodecs = serializer.NewCodecFactory(appsScheme)
	log        = logf.Log.WithName("controller_alamedaservice")
)

func init() {
	if err := appsv1.AddToScheme(appsScheme); err != nil {
		log.Error(err, "Fail AddToScheme")
	}
}

func ReadDeploymentV1(objBytes []byte) *appsv1.Deployment {
	requiredObj, err := runtime.Decode(appsCodecs.UniversalDecoder(appsv1.SchemeGroupVersion), objBytes)
	if err != nil {
		log.Error(err, "Fail to ReadDeploymentV1OrDie")

	}
	return requiredObj.(*appsv1.Deployment)
}
