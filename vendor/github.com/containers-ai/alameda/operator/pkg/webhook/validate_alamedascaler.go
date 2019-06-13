package webhook

import (
	"context"
	"net/http"

	autoscalingv1alpha1 "github.com/containers-ai/alameda/operator/pkg/apis/autoscaling/v1alpha1"
	"github.com/containers-ai/alameda/pkg/utils"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	admissiontypes "sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

type alamedaScalerLabeler struct {
	client  client.Client
	decoder admissiontypes.Decoder
}

var _ admission.Handler = &alamedaScalerLabeler{}

func (labeler *alamedaScalerLabeler) Handle(ctx context.Context, req admissiontypes.Request) admissiontypes.Response {
	alamedaScaler := &autoscalingv1alpha1.AlamedaScaler{}

	err := labeler.decoder.Decode(req, alamedaScaler)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}
	scope.Debugf("AlamedaScaler received to validate as following %s", utils.InterfaceToString(alamedaScaler))
	res, err := labeler.validateAlamedaScalersFn(ctx, alamedaScaler)
	if err != nil {
		return admission.ValidationResponse(res, err.Error())
	}
	return admission.ValidationResponse(res, "")
}

var _ inject.Decoder = &alamedaScalerLabeler{}

// InjectDecoder injects the decoder into the alamedaScalerLabeler
func (labeler *alamedaScalerLabeler) InjectDecoder(d admissiontypes.Decoder) error {
	labeler.decoder = d
	return nil
}

var _ inject.Client = &alamedaScalerLabeler{}

// InjectClient injects the client into the alamedaScalerLabeler
func (labeler *alamedaScalerLabeler) InjectClient(c client.Client) error {
	labeler.client = c
	return nil
}

// validateAlamedaScalersFn validate the given alamedaScalerLabeler
func (labeler *alamedaScalerLabeler) validateAlamedaScalersFn(ctx context.Context, alamedaScaler *autoscalingv1alpha1.AlamedaScaler) (bool, error) {
	return isScalerValid(&labeler.client, &validatingObject{
		namespace:           alamedaScaler.GetNamespace(),
		name:                alamedaScaler.GetName(),
		kind:                alamedaScaler.GetObjectKind().GroupVersionKind().Kind,
		selectorMatchLabels: alamedaScaler.Spec.Selector.MatchLabels,
	})
}

func GetAlamedaScalerHandler() *alamedaScalerLabeler {
	return &alamedaScalerLabeler{}
}
