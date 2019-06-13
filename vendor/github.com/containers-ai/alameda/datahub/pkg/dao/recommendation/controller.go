package recommendation

import (
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
)

// ContainerOperation defines container measurement operation of recommendation database
type ControllerOperation interface {
	AddControllerRecommendations([]*datahub_v1alpha1.ControllerRecommendation) error
	ListControllerRecommendations(controllerNamespacedName *datahub_v1alpha1.NamespacedName, queryCondition *datahub_v1alpha1.QueryCondition) ([]*datahub_v1alpha1.ControllerRecommendation, error)
}
