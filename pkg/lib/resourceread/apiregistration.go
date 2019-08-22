package resourceread

import (
	"github.com/pkg/errors"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	apiregistrationv1beta1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1beta1"
)

var (
	apiregistrationScheme = runtime.NewScheme()
	apiregistrationCodecs = serializer.NewCodecFactory(apiregistrationScheme)
)

func init() {
	if err := apiregistrationv1beta1.AddToScheme(apiregistrationScheme); err != nil {
		log.Error(err, "Failed to AddToScheme")
	}
}

func ReadAPIService(objBytes []byte) (*apiregistrationv1beta1.APIService, error) {
	requiredObj, err := runtime.Decode(apiregistrationCodecs.UniversalDecoder(apiregistrationv1beta1.SchemeGroupVersion), objBytes)
	if err != nil {
		return nil, errors.Errorf("decode APIService failed: %s", err.Error())
	}
	return requiredObj.(*apiregistrationv1beta1.APIService), nil
}
