package webhook

import (
	"fmt"

	"github.com/containers-ai/alameda/operator/pkg/utils/resources"
	"github.com/containers-ai/alameda/pkg/utils"
	"github.com/containers-ai/alameda/pkg/utils/kubernetes"
	openshift_apps_v1 "github.com/openshift/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type validatingObject struct {
	namespace           string
	name                string
	kind                string
	labels              map[string]string
	selectorMatchLabels map[string]string
}

func isTopControllerValid(client *client.Client, topCtl *validatingObject) (bool, error) {
	listResources := resources.NewListResources(*client)
	// TODO: may use ListAllAlamedaScaler if alamedascaler supports selectNamespace option
	scalers, err := listResources.ListNamespaceAlamedaScaler(topCtl.namespace)
	scope.Debugf("%v alamedascaler in namespace %s to check deplicated selection of %s %s", len(scalers), topCtl.namespace, topCtl.kind, topCtl.name)
	if err != nil {
		return false, err
	}
	matchedScalerList := []*validatingObject{}
	for _, scaler := range scalers {
		if isLabelSelected(scaler.Spec.Selector.MatchLabels, topCtl.labels) {
			matchedScalerList = append(matchedScalerList, &validatingObject{
				name:      scaler.GetName(),
				namespace: scaler.GetNamespace(),
			})
		}
	}
	if len(matchedScalerList) > 1 {
		matchedNamesapcedNames := fmt.Sprintf("%s/%s", matchedScalerList[0].namespace, matchedScalerList[0].name)
		for idx, matched := range matchedScalerList {
			if idx > 0 {
				matchedNamesapcedNames = fmt.Sprintf("%s, %s/%s", matchedNamesapcedNames, matched.namespace, matched.name)
			}
		}

		return false, fmt.Errorf("%s (%s/%s) is selected by more than 1 alamedascaler (%s)", topCtl.kind, topCtl.namespace, topCtl.name, matchedNamesapcedNames)
	}
	return true, nil
}

func getSelectedDeploymentConfigs(listResources *resources.ListResources, namespace string, selectorMatchLabels map[string]string) ([]openshift_apps_v1.DeploymentConfig, error) {
	okdCluster, err := kubernetes.IsOKDCluster()
	if err != nil {
		scope.Errorf(err.Error())
	}
	if !okdCluster {
		return []openshift_apps_v1.DeploymentConfig{}, nil
	}
	// TODO: may use ListDeploymentConfigsByLabels if alamedascaler supports selectNamespace option
	return listResources.ListDeploymentConfigsByNamespaceLabels(namespace, selectorMatchLabels)
}

func isScalerValid(client *client.Client, scalerObj *validatingObject) (bool, error) {
	listResources := resources.NewListResources(*client)
	// TODO: may use ListAllAlamedaScaler if alamedascaler supports selectNamespace option
	scalers, err := listResources.ListNamespaceAlamedaScaler(scalerObj.namespace)
	if err != nil {
		return false, err
	}
	// TODO: may use ListDeploymentsByLabels if alamedascaler supports selectNamespace option
	selectedDeployments, err := listResources.ListDeploymentsByNamespaceLabels(scalerObj.namespace, scalerObj.selectorMatchLabels)
	if err != nil {
		return false, err
	}

	selectedDeploymentConfigs, err := getSelectedDeploymentConfigs(listResources, scalerObj.namespace, scalerObj.selectorMatchLabels)
	if err != nil {
		return false, err
	}

	for _, scaler := range scalers {
		if scaler.GetNamespace() == scalerObj.namespace && scaler.GetName() == scalerObj.name {
			continue
		}

		for _, selectedDeployment := range selectedDeployments {
			if _, ok := scaler.Status.AlamedaController.Deployments[fmt.Sprintf("%s/%s", selectedDeployment.GetNamespace(), selectedDeployment.GetName())]; ok {
				return false, fmt.Errorf("Deployment %s/%s selected by scaler %s/%s is already selected by scaler %s/%s",
					selectedDeployment.GetNamespace(), selectedDeployment.GetName(),
					scalerObj.namespace, scalerObj.name, scaler.GetNamespace(), scaler.GetName())
			}
		}
		for _, selectedDeploymentConfig := range selectedDeploymentConfigs {
			if _, ok := scaler.Status.AlamedaController.DeploymentConfigs[fmt.Sprintf("%s/%s", selectedDeploymentConfig.GetNamespace(), selectedDeploymentConfig.GetName())]; ok {
				return false, fmt.Errorf("DeploymentConfig %s/%s selected by scaler %s/%s is already selected by scaler %s/%s",
					selectedDeploymentConfig.GetNamespace(), selectedDeploymentConfig.GetName(),
					scalerObj.namespace, scalerObj.name, scaler.GetNamespace(), scaler.GetName())
			}
		}
	}
	return true, nil
}

func isLabelSelected(selector, label map[string]string) bool {
	isSelected := true
	scope.Debugf("Check label is selected by selector.")
	scope.Debugf("Selector is %s.", utils.InterfaceToString(selector))
	scope.Debugf("Label is %s.", utils.InterfaceToString(label))
	for selKey, selVal := range selector {
		if _, ok := label[selKey]; !ok {
			isSelected = false
			break
		}
		if label[selKey] != selVal {
			isSelected = false
			break
		}
	}
	if isSelected {
		scope.Debugf("Label is matched by selector.")
	} else {
		scope.Debugf("Label is not matched by selector.")
	}
	return isSelected
}
