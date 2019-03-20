package resourceread

import (
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	apiExtensionsScheme = runtime.NewScheme()
	apiExtensionsCodecs = serializer.NewCodecFactory(apiExtensionsScheme)
)

func init() {
	if err := apiextv1beta1.AddToScheme(apiExtensionsScheme); err != nil {
		log.Error(err, "Failed to AddToScheme")
	}
}

func ReadCustomResourceDefinitionV1Beta1OrDie(objBytes []byte) *apiextv1beta1.CustomResourceDefinition {
	requiredObj, err := runtime.Decode(apiExtensionsCodecs.UniversalDecoder(apiextv1beta1.SchemeGroupVersion), objBytes)
	if err != nil {
		log.Error(err, "Failed to ReadCustomResourceDefinitionV1Beta1OrDie")
	}
	return requiredObj.(*apiextv1beta1.CustomResourceDefinition)
}
