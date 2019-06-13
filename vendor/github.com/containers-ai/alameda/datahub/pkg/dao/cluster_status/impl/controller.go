package impl

import (
	influxdb_repository "github.com/containers-ai/alameda/datahub/pkg/repository/influxdb"
	influxdb_repository_cluster_status "github.com/containers-ai/alameda/datahub/pkg/repository/influxdb/cluster_status"
	datahub_api "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
)

type Controller struct {
	InfluxDBConfig influxdb_repository.Config
}

func (c *Controller) CreateControllers(controllers []*datahub_api.Controller) error {
	controllerRepository := influxdb_repository_cluster_status.NewControllerRepository(&c.InfluxDBConfig)
	return controllerRepository.CreateControllers(controllers)
}

func (c *Controller) ListControllers(in *datahub_api.ListControllersRequest) ([]*datahub_api.Controller, error) {
	controllerRepository := influxdb_repository_cluster_status.NewControllerRepository(&c.InfluxDBConfig)
	return controllerRepository.ListControllers(in)
}

func (c *Controller) DeleteControllers(in *datahub_api.DeleteControllersRequest) error {
	controllerRepository := influxdb_repository_cluster_status.NewControllerRepository(&c.InfluxDBConfig)
	return controllerRepository.DeleteControllers(in)
}
