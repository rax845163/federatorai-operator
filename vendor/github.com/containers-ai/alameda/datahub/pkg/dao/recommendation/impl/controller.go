package impl

import (
	influxdb_repository "github.com/containers-ai/alameda/datahub/pkg/repository/influxdb"

	influxdb_repository_recommendation "github.com/containers-ai/alameda/datahub/pkg/repository/influxdb/recommendation"
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
)

type Controller struct {
	InfluxDBConfig influxdb_repository.Config
}

func (c *Controller) AddControllerRecommendations(controllerRecommendations []*datahub_v1alpha1.ControllerRecommendation) error {
	controllerRepository := influxdb_repository_recommendation.NewControllerRepository(&c.InfluxDBConfig)
	return controllerRepository.CreateControllerRecommendations(controllerRecommendations)
}

func (c *Controller) ListControllerRecommendations(in *datahub_v1alpha1.ListControllerRecommendationsRequest) ([]*datahub_v1alpha1.ControllerRecommendation, error) {
	controllerRepository := influxdb_repository_recommendation.NewControllerRepository(&c.InfluxDBConfig)
	return controllerRepository.ListControllerRecommendations(in)
}
