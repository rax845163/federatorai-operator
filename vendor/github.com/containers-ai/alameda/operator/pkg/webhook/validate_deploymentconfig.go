package webhook

import (
	"context"
	"net/http"

	"github.com/containers-ai/alameda/pkg/utils"
	osappsapi "github.com/openshift/api/apps/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	admissiontypes "sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

type deploymentConfigLabeler struct {
	client  client.Client
	decoder admissiontypes.Decoder
}

var _ admission.Handler = &deploymentConfigLabeler{}

func (labeler *deploymentConfigLabeler) Handle(ctx context.Context, req admissiontypes.Request) admissiontypes.Response {
	deploymentConfig := &osappsapi.DeploymentConfig{}

	err := labeler.decoder.Decode(req, deploymentConfig)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}
	scope.Debugf("DeploymentConfig received to validate as following %s", utils.InterfaceToString(deploymentConfig))
	res, err := labeler.validateDeploymentConfigsFn(ctx, deploymentConfig)
	if err != nil {
		return admission.ValidationResponse(res, err.Error())
	}
	return admission.ValidationResponse(res, "")
}

var _ inject.Decoder = &deploymentLabeler{}

// InjectDecoder injects the decoder into the deploymentConfigLabeler
func (labeler *deploymentConfigLabeler) InjectDecoder(d admissiontypes.Decoder) error {
	labeler.decoder = d
	return nil
}

// deploymentConfigLabeler implements inject.Client.
var _ inject.Client = &deploymentConfigLabeler{}

// InjectClient injects the client into the deploymentConfigLabeler
func (labeler *deploymentConfigLabeler) InjectClient(c client.Client) error {
	labeler.client = c
	return nil
}

// validateDeploymentsFn validate the given deploymentConfig
func (labeler *deploymentConfigLabeler) validateDeploymentConfigsFn(ctx context.Context, deploymentConfig *osappsapi.DeploymentConfig) (bool, error) {
	return isTopControllerValid(&labeler.client, &validatingObject{
		namespace: deploymentConfig.GetNamespace(),
		name:      deploymentConfig.GetName(),
		kind:      deploymentConfig.GetObjectKind().GroupVersionKind().Kind,
		labels:    deploymentConfig.GetLabels(),
	})
}

func GetDeploymentConfigHandler() *deploymentConfigLabeler {
	return &deploymentConfigLabeler{}
}
