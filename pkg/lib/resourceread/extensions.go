package resourceread

import (
	//"github.com/pkg/errors"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	extensionScheme = runtime.NewScheme()
	extensionCodecs = serializer.NewCodecFactory(extensionScheme)
)

func init() {
	if err := v1beta1.AddToScheme(extensionScheme); err != nil {
		log.Error(err, "Fail AddToScheme")
	}
}

func ReadPodSecurityPolicyV1beta1(objBytes []byte) *v1beta1.PodSecurityPolicy {
	requiredObj, err := runtime.Decode(extensionCodecs.UniversalDecoder(v1beta1.SchemeGroupVersion), objBytes)
	if err != nil {
		log.Error(err, "Fail ReadPodSecurityPolicyV1beta1")
	}
	return requiredObj.(*v1beta1.PodSecurityPolicy)
}
