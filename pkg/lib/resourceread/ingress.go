package resourceread

import (
	ingressv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	ingressScheme = runtime.NewScheme()
	ingressCodecs = serializer.NewCodecFactory(ingressScheme)
)

func init() {
	if err := ingressv1beta1.AddToScheme(ingressScheme); err != nil {
		log.Error(err, "Fail AddToScheme")
	}
}

func ReadIngressv1beta1(objBytes []byte) *ingressv1beta1.Ingress {
	requiredObj, err := runtime.Decode(ingressCodecs.UniversalDecoder(ingressv1beta1.SchemeGroupVersion), objBytes)
	if err != nil {
		log.Error(err, "Fail ReadIngressv1beta1")
	}
	return requiredObj.(*ingressv1beta1.Ingress)
}
