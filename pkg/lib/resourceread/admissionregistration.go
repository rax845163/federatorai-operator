package resourceread

import (
	"github.com/pkg/errors"

	"k8s.io/api/admissionregistration/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	admissionregistrationScheme = runtime.NewScheme()
	admissionregistrationCodecs = serializer.NewCodecFactory(admissionregistrationScheme)
)

func init() {
	if err := v1beta1.AddToScheme(admissionregistrationScheme); err != nil {
		log.Error(err, "Failed to AddToScheme")
	}
}

func ReadMutatingWebhookConfiguration(objBytes []byte) (*v1beta1.MutatingWebhookConfiguration, error) {
	requiredObj, err := runtime.Decode(admissionregistrationCodecs.UniversalDecoder(v1beta1.SchemeGroupVersion), objBytes)
	if err != nil {
		return nil, errors.Errorf("decode MutatingWebhookConfiguration failed: %s", err.Error())
	}
	return requiredObj.(*v1beta1.MutatingWebhookConfiguration), nil
}

func ReadValidatingWebhookConfiguration(objBytes []byte) (*v1beta1.ValidatingWebhookConfiguration, error) {
	requiredObj, err := runtime.Decode(admissionregistrationCodecs.UniversalDecoder(v1beta1.SchemeGroupVersion), objBytes)
	if err != nil {
		return nil, errors.Errorf("decode ValidatingWebhookConfiguration failed: %s", err.Error())
	}
	return requiredObj.(*v1beta1.ValidatingWebhookConfiguration), nil
}
