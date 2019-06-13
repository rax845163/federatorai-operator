package recommendation

import (
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
)

// ContainerOperation defines container measurement operation of recommendation database
type ContainerOperation interface {
	AddPodRecommendations(in *datahub_v1alpha1.CreatePodRecommendationsRequest) error
	ListPodRecommendations(in *datahub_v1alpha1.ListPodRecommendationsRequest) ([]*datahub_v1alpha1.PodRecommendation, error)
	ListAvailablePodRecommendations(*datahub_v1alpha1.ListPodRecommendationsRequest) ([]*datahub_v1alpha1.PodRecommendation, error)
}
