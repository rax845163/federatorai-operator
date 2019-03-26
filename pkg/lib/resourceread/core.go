package resourceread

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	coreScheme = runtime.NewScheme()
	coreCodecs = serializer.NewCodecFactory(coreScheme)
)

func init() {
	if err := corev1.AddToScheme(coreScheme); err != nil {
		log.Error(err, "Fail AddToScheme")
	}
}

func ReadConfigMapV1(objBytes []byte) *corev1.ConfigMap {
	requiredObj, err := runtime.Decode(coreCodecs.UniversalDecoder(corev1.SchemeGroupVersion), objBytes)
	if err != nil {
		log.Error(err, "Fail ReadConfigMapV1")
	}
	return requiredObj.(*corev1.ConfigMap)
}

func ReadServiceAccountV1(objBytes []byte) *corev1.ServiceAccount {
	requiredObj, err := runtime.Decode(coreCodecs.UniversalDecoder(corev1.SchemeGroupVersion), objBytes)
	if err != nil {
		log.Error(err, "Fail ReadServiceAccountV1")
	}
	return requiredObj.(*corev1.ServiceAccount)
}

func ReadServiceV1(objBytes []byte) *corev1.Service {
	requiredObj, err := runtime.Decode(coreCodecs.UniversalDecoder(corev1.SchemeGroupVersion), objBytes)
	if err != nil {
		log.Error(err, "Fail ReadServiceV1")
	}
	return requiredObj.(*corev1.Service)
}

func ReadSecretV1(objBytes []byte) *corev1.Secret {
	requiredObj, err := runtime.Decode(coreCodecs.UniversalDecoder(corev1.SchemeGroupVersion), objBytes)
	if err != nil {
		log.Error(err, "Fail ReadSecretV1")
	}
	return requiredObj.(*corev1.Secret)
}

func ReadPersistentVolumeClaimV1(objBytes []byte) *corev1.PersistentVolumeClaim {
	requiredObj, err := runtime.Decode(coreCodecs.UniversalDecoder(corev1.SchemeGroupVersion), objBytes)
	if err != nil {
		log.Error(err, "Fail ReadPersistentVolumeClaimV1")
	}
	return requiredObj.(*corev1.PersistentVolumeClaim)
}
