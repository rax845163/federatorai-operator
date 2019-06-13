package resourceread

import (
	autoscaling_v1alpha1 "github.com/containers-ai/alameda/operator/pkg/apis/autoscaling/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	scalerScheme = runtime.NewScheme()
	scalerCodecs = serializer.NewCodecFactory(scalerScheme)
)

func init() {
	if err := autoscaling_v1alpha1.AddToScheme(scalerScheme); err != nil {
		log.Error(err, "Fail AddToScheme")
	}
}

func ReadScalerV1(objBytes []byte) *autoscaling_v1alpha1.AlamedaScaler {
	requiredObj, err := runtime.Decode(scalerCodecs.UniversalDecoder(autoscaling_v1alpha1.SchemeGroupVersion), objBytes)
	if err != nil {
		log.Error(err, "Fail to ReadScalerV1OrDie")

	}
	return requiredObj.(*autoscaling_v1alpha1.AlamedaScaler)
}
