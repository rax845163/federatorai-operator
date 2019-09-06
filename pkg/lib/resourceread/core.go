package resourceread

import (
	"github.com/pkg/errors"
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

func ReadSecretV1(objBytes []byte) (*corev1.Secret, error) {
	requiredObj, err := runtime.Decode(coreCodecs.UniversalDecoder(corev1.SchemeGroupVersion), objBytes)
	if err != nil {
		return nil, errors.Errorf("failed to decode core v1 secret: %s", err.Error())
	}
	secret, ok := requiredObj.(*corev1.Secret)
	if !ok {
		return nil, errors.Errorf("failed to convert to core v1 secret")
	}
	return secret, nil
}

func ReadPersistentVolumeClaimV1(objBytes []byte) *corev1.PersistentVolumeClaim {
	requiredObj, err := runtime.Decode(coreCodecs.UniversalDecoder(corev1.SchemeGroupVersion), objBytes)
	if err != nil {
		log.Error(err, "Fail ReadPersistentVolumeClaimV1")
	}
	return requiredObj.(*corev1.PersistentVolumeClaim)
}
