package datahub

import (
	"strings"

	"github.com/containers-ai/alameda/admission-controller/pkg/validator/controller"
	autoscaling_v1alpha1 "github.com/containers-ai/alameda/operator/pkg/apis/autoscaling/v1alpha1"
	"github.com/containers-ai/alameda/pkg/utils/log"
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	"github.com/pkg/errors"
	context "golang.org/x/net/context"
	"google.golang.org/genproto/googleapis/rpc/code"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	scope = log.RegisterScope("conrtoller-validator", "Datahub conrtoller validator", 0)
)

type validator struct {
	datahubServiceClient datahub_v1alpha1.DatahubServiceClient
	sigsK8SClient        client.Client
}

// NewControllerValidator returns controller validator which fetch controller information from containers-ai/alameda Datahub
func NewControllerValidator(datahubServiceClient datahub_v1alpha1.DatahubServiceClient, sigsK8SClient client.Client) controller.Validator {
	return &validator{
		datahubServiceClient: datahubServiceClient,
		sigsK8SClient:        sigsK8SClient,
	}
}

func (v *validator) IsControllerEnabledExecution(namespace, name, kind string) (bool, error) {

	datahubKind, exist := datahub_v1alpha1.Kind_value[strings.ToUpper(kind)]
	if !exist {
		return false, errors.Errorf("no matched datahub kind for kind: %s", kind)
	}

	ctx := buildDefaultRequestContext()
	req := &datahub_v1alpha1.ListControllersRequest{
		NamespacedName: &datahub_v1alpha1.NamespacedName{
			Namespace: namespace,
			Name:      name,
		},
	}
	scope.Debugf("query ListControllers to datahub, send request: %+v", req)
	resp, err := v.datahubServiceClient.ListControllers(ctx, req)
	scope.Debugf("query ListControllers to datahub, received response: %+v", resp)
	if err != nil {
		return false, errors.Errorf("query ListControllers to datahub failed: errMsg: %s", err.Error())
	}
	if resp.Status == nil {
		return false, errors.New("receive nil status from datahub")
	} else if resp.Status.Code != int32(code.Code_OK) {
		return false, errors.Errorf("status code not 0: receive status code: %d,message: %s", resp.Status.Code, resp.Status.Message)
	}

	controllers := resp.Controllers
	indices := getMatchedControllerIndices(controllers, namespace, name, datahub_v1alpha1.Kind(datahubKind))
	if len(indices) == 0 {
		return false, errors.Errorf("cannot find matched controller (%s/%s ,kind: %s) from datahub", namespace, name, kind)
	}
	controller := controllers[0]

	alamedaScalerIndices := getMatchedResourceIndicesWithKind(controller.OwnerInfo, datahub_v1alpha1.Kind_ALAMEDASCALER)
	if len(alamedaScalerIndices) == 0 {
		return false, errors.Errorf("cannot find matched AlamedaScaler to controller (%s/%s ,kind: %s) from datahub", namespace, name, kind)
	}
	alamedaScalerInfo := controller.OwnerInfo[alamedaScalerIndices[0]]
	alamedaScalerNamespacedName := alamedaScalerInfo.NamespacedName
	if alamedaScalerNamespacedName == nil {
		return false, errors.Errorf("getting AlamedaScaler with empty NamespacedName controller (%s/%s ,kind: %s) from datahub", namespace, name, kind)
	} else if alamedaScalerNamespacedName.Namespace == "" || alamedaScalerNamespacedName.Name == "" {
		return false, errors.Errorf("getting AlamedaScaler with empty NamespacedName controller (%s/%s ,kind: %s) from datahub", namespace, name, kind)
	}

	alamedaScaler := autoscaling_v1alpha1.AlamedaScaler{}
	err = v.sigsK8SClient.Get(
		ctx,
		client.ObjectKey{
			Namespace: alamedaScalerNamespacedName.Namespace,
			Name:      alamedaScalerNamespacedName.Name,
		},
		&alamedaScaler)
	if err != nil {
		return false, errors.Errorf("get AlamedaScaler from k8s failed: %s", err.Error())
	}
	scope.Debugf(`get monitoring AlamedaScaler for controller, controller:{
		namespace: %s,
		name: %s,
		kind: %s
	}, AlamedaScaler:{
		namespace: %s,
		name: %s
	}`, namespace, name, kind, alamedaScaler.Namespace, alamedaScaler.Name)
	return alamedaScaler.Spec.EnableExecution, nil
}

func getMatchedControllerIndices(controllers []*datahub_v1alpha1.Controller, namespace string, name string, kind datahub_v1alpha1.Kind) []int {

	controllersInfo := make([]*datahub_v1alpha1.ResourceInfo, len(controllers))
	for i, controller := range controllers {
		controllersInfo[i] = controller.ControllerInfo
	}

	return getMatchedResourceIndex(controllersInfo, namespace, name, kind)
}

func getMatchedResourceIndex(resourcesInfo []*datahub_v1alpha1.ResourceInfo, namespace string, name string, kind datahub_v1alpha1.Kind) []int {

	indices := make([]int, 0)
	for i, resourceInfo := range resourcesInfo {

		if resourceInfo == nil {
			continue
		} else if resourceInfo.NamespacedName == nil {
			continue
		} else if resourceInfo.NamespacedName.Namespace == namespace && resourceInfo.NamespacedName.Name == name && resourceInfo.Kind == kind {
			indices = append(indices, i)
		}
	}
	return indices
}

func getMatchedResourceIndicesWithKind(resourcesInfo []*datahub_v1alpha1.ResourceInfo, kind datahub_v1alpha1.Kind) []int {

	indices := make([]int, 0)
	for i, resourceInfo := range resourcesInfo {

		if resourceInfo == nil {
			continue
		} else if resourceInfo.Kind == kind {
			indices = append(indices, i)
		}
	}
	return indices
}

func buildDefaultRequestContext() context.Context {
	return context.TODO()
}
