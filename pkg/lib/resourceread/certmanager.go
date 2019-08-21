package resourceread

import (
	"github.com/pkg/errors"

	v1alpha1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	certmanagerScheme = runtime.NewScheme()
	certmanagerCodecs = serializer.NewCodecFactory(certmanagerScheme)
)

func init() {
	if err := v1alpha1.AddToScheme(certmanagerScheme); err != nil {
		log.Error(err, "Failed to AddToScheme")
	}
}

func ReadCertificate(objBytes []byte) (*v1alpha1.Certificate, error) {
	requiredObj, err := runtime.Decode(certmanagerCodecs.UniversalDecoder(v1alpha1.SchemeGroupVersion), objBytes)
	if err != nil {
		return nil, errors.Errorf("decode Certificate failed: %s", err.Error())
	}
	return requiredObj.(*v1alpha1.Certificate), nil
}

func ReadIssuer(objBytes []byte) (*v1alpha1.Issuer, error) {
	requiredObj, err := runtime.Decode(certmanagerCodecs.UniversalDecoder(v1alpha1.SchemeGroupVersion), objBytes)
	if err != nil {
		return nil, errors.Errorf("decode Issuer failed: %s", err.Error())
	}
	return requiredObj.(*v1alpha1.Issuer), nil
}
