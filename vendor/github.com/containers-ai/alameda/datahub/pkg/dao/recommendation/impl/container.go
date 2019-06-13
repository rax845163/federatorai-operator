package impl

import (
	influxdb_repository "github.com/containers-ai/alameda/datahub/pkg/repository/influxdb"

	influxdb_repository_recommendation "github.com/containers-ai/alameda/datahub/pkg/repository/influxdb/recommendation"
	"github.com/containers-ai/alameda/pkg/utils/log"
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
)

var (
	scope = log.RegisterScope("recommendation_dao_implement", "recommended dao implement", 0)
)

// Container Implements ContainerOperation interface
type Container struct {
	InfluxDBConfig influxdb_repository.Config
}

// AddPodRecommendations add pod recommendations to database
func (container *Container) AddPodRecommendations(in *datahub_v1alpha1.CreatePodRecommendationsRequest) error {
	containerRepository := influxdb_repository_recommendation.NewContainerRepository(&container.InfluxDBConfig)
	return containerRepository.CreateContainerRecommendations(in)
}

// ListPodRecommendations list pod recommendations
func (container *Container) ListPodRecommendations(in *datahub_v1alpha1.ListPodRecommendationsRequest) ([]*datahub_v1alpha1.PodRecommendation, error) {
	containerRepository := influxdb_repository_recommendation.NewContainerRepository(&container.InfluxDBConfig)
	return containerRepository.ListContainerRecommendations(in)
}

func (container *Container) ListAvailablePodRecommendations(in *datahub_v1alpha1.ListPodRecommendationsRequest) ([]*datahub_v1alpha1.PodRecommendation, error) {
	containerRepository := influxdb_repository_recommendation.NewContainerRepository(&container.InfluxDBConfig)
	return containerRepository.ListAvailablePodRecommendations(in)
}
