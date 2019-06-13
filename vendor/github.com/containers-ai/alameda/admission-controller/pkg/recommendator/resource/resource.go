package resource

import (
	"time"

	core_v1 "k8s.io/api/core/v1"
)

type ContainerResourceRecommendation struct {
	Name     string
	Limits   core_v1.ResourceList
	Requests core_v1.ResourceList
}

type PodResourceRecommendation struct {
	Namespace                        string
	Name                             string
	TopControllerKind                string
	TopControllerName                string
	ContainerResourceRecommendations []*ContainerResourceRecommendation
	ValidStartTime                   time.Time
	ValidEndTime                     time.Time
}

type ResourceRecommendator interface {
	ListControllerPodResourceRecommendations(ListControllerPodResourceRecommendationsRequest) ([]*PodResourceRecommendation, error)
}

type ListControllerPodResourceRecommendationsRequest struct {
	Namespace string
	Name      string
	Kind      string
	Time      *time.Time
}
