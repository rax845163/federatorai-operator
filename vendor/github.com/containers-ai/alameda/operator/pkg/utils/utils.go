package utils

import (
	"fmt"
	"strings"

	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"

	openshift_api_apps_v1 "github.com/openshift/api/apps"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// GetNamespacedNameKey returns string "namespaced/name"
func GetNamespacedNameKey(namespace, name string) string {
	return fmt.Sprintf("%s/%s", namespace, name)
}

//ParseResourceLinkForTopController parses resourcelink string to get top controller information
func ParseResourceLinkForTopController(resourceLink string) (*datahub_v1alpha1.TopController, error) {
	res := strings.Split(resourceLink, "/")
	if len(res) >= 5 {
		kind := datahub_v1alpha1.Kind_POD
		switch res[3] {
		case "deployments":
			kind = datahub_v1alpha1.Kind_DEPLOYMENT
		case "deploymentconfigs":
			kind = datahub_v1alpha1.Kind_DEPLOYMENTCONFIG
		case "statefulsets":
			kind = datahub_v1alpha1.Kind_STATEFULSET
		default:
			kind = datahub_v1alpha1.Kind_POD
		}
		return &datahub_v1alpha1.TopController{
			NamespacedName: &datahub_v1alpha1.NamespacedName{
				Namespace: res[2],
				Name:      res[4],
			},
			Kind: kind,
		}, nil
	}
	return &datahub_v1alpha1.TopController{}, fmt.Errorf("resource link format is not correct")
}

var (
	hasOpenshiftAPIAppsV1 *bool
)

// ServerHasOpenshiftAPIAppsV1 returns true if the api-server has apiGroup named in "apps.openshift.io"
func ServerHasOpenshiftAPIAppsV1() (bool, error) {

	if hasOpenshiftAPIAppsV1 == nil {
		if exist, err := serverHasAPIGroup(openshift_api_apps_v1.GroupName); err != nil {
			return false, err
		} else {
			hasOpenshiftAPIAppsV1 = &exist
		}
	}

	return *hasOpenshiftAPIAppsV1, nil
}

func serverHasAPIGroup(apiGroupName string) (bool, error) {

	config, err := config.GetConfig()
	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return false, err
	}
	apiGroups, err := k8sClient.ServerGroups()
	if err != nil {
		return false, err
	}
	for _, apiGroup := range apiGroups.Groups {
		if apiGroup.Name == apiGroupName {
			return true, nil
		}
	}
	return false, nil
}
