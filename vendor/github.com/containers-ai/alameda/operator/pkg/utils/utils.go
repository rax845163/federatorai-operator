package utils

import (
	"fmt"
	"strings"

	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
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
	return &datahub_v1alpha1.TopController{}, fmt.Errorf("Resource link format is not correct.")
}
