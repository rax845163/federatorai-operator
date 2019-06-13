package mock

import (
	"time"

	"github.com/containers-ai/alameda/admission-controller/pkg/recommendator/resource"
	core_v1 "k8s.io/api/core/v1"
	k8s_resource "k8s.io/apimachinery/pkg/api/resource"
)

var _ resource.ResourceRecommendator = &mockResourceRecommendator{}

type mockResourceRecommendator struct {
}

func NewMockResourceRecommendator() resource.ResourceRecommendator {
	return &mockResourceRecommendator{}
}

func (m *mockResourceRecommendator) ListControllerPodResourceRecommendations(req resource.ListControllerPodResourceRecommendationsRequest) ([]*resource.PodResourceRecommendation, error) {

	currentTime := time.Now()
	twoDayBefore := currentTime.AddDate(0, 0, -2)
	oneDayBefore := currentTime.AddDate(0, 0, -1)
	oneDayAfter := currentTime.AddDate(0, 0, 1)

	recommendations := []*resource.PodResourceRecommendation{
		&resource.PodResourceRecommendation{
			Namespace: "webapp",
			Name:      "nginx-1",
			ContainerResourceRecommendations: []*resource.ContainerResourceRecommendation{
				&resource.ContainerResourceRecommendation{
					Name: "nginx",
					Limits: core_v1.ResourceList{
						core_v1.ResourceCPU:    k8s_resource.MustParse("200m"),
						core_v1.ResourceMemory: k8s_resource.MustParse("256Mi"),
					},
					Requests: core_v1.ResourceList{
						core_v1.ResourceCPU:    k8s_resource.MustParse("100m"),
						core_v1.ResourceMemory: k8s_resource.MustParse("128Mi"),
					},
				},
			},
			ValidStartTime: twoDayBefore,
			ValidEndTime:   oneDayBefore,
		},
		&resource.PodResourceRecommendation{
			Namespace: "webapp",
			Name:      "nginx-2",
			ContainerResourceRecommendations: []*resource.ContainerResourceRecommendation{
				&resource.ContainerResourceRecommendation{
					Name: "nginx",
					Limits: core_v1.ResourceList{
						core_v1.ResourceCPU:    k8s_resource.MustParse("100m"),
						core_v1.ResourceMemory: k8s_resource.MustParse("128Mi"),
					},
					Requests: core_v1.ResourceList{
						core_v1.ResourceCPU:    k8s_resource.MustParse("50m"),
						core_v1.ResourceMemory: k8s_resource.MustParse("64Mi"),
					},
				},
			},
			ValidStartTime: oneDayBefore,
			ValidEndTime:   oneDayAfter,
		},
		&resource.PodResourceRecommendation{
			Namespace: "webapp",
			Name:      "nginx-3",
			ContainerResourceRecommendations: []*resource.ContainerResourceRecommendation{
				&resource.ContainerResourceRecommendation{
					Name: "nginx",
					Limits: core_v1.ResourceList{
						core_v1.ResourceCPU:    k8s_resource.MustParse("110m"),
						core_v1.ResourceMemory: k8s_resource.MustParse("130Mi"),
					},
					Requests: core_v1.ResourceList{
						core_v1.ResourceCPU:    k8s_resource.MustParse("60m"),
						core_v1.ResourceMemory: k8s_resource.MustParse("60Mi"),
					},
				},
			},
			ValidStartTime: oneDayBefore,
			ValidEndTime:   oneDayAfter,
		},
		&resource.PodResourceRecommendation{
			Namespace: "webapp",
			Name:      "nginx-3",
			ContainerResourceRecommendations: []*resource.ContainerResourceRecommendation{
				&resource.ContainerResourceRecommendation{
					Name: "nginx",
					Limits: core_v1.ResourceList{
						core_v1.ResourceCPU:    k8s_resource.MustParse("120m"),
						core_v1.ResourceMemory: k8s_resource.MustParse("140Mi"),
					},
					Requests: core_v1.ResourceList{
						core_v1.ResourceCPU:    k8s_resource.MustParse("70m"),
						core_v1.ResourceMemory: k8s_resource.MustParse("70Mi"),
					},
				},
			},
			ValidStartTime: oneDayBefore,
			ValidEndTime:   oneDayAfter,
		},
	}

	return recommendations, nil
}
