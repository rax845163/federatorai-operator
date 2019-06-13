package impl

import (
	cluster_status_dao "github.com/containers-ai/alameda/datahub/pkg/dao/cluster_status"
	influxdb_repository "github.com/containers-ai/alameda/datahub/pkg/repository/influxdb"
	influxdb_repository_cluster_status "github.com/containers-ai/alameda/datahub/pkg/repository/influxdb/cluster_status"
	datahub_api "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	"github.com/pkg/errors"
)

// Implement Node interface
type Node struct {
	InfluxDBConfig influxdb_repository.Config
}

func (node *Node) RegisterAlamedaNodes(alamedaNodes []*datahub_v1alpha1.Node) error {
	nodeRepository := influxdb_repository_cluster_status.NewNodeRepository(&node.InfluxDBConfig)
	return nodeRepository.AddAlamedaNodes(alamedaNodes)
}

func (node *Node) DeregisterAlamedaNodes(alamedaNodes []*datahub_v1alpha1.Node) error {
	nodeRepository := influxdb_repository_cluster_status.NewNodeRepository(&node.InfluxDBConfig)
	return nodeRepository.RemoveAlamedaNodes(alamedaNodes)
}

func (node *Node) ListAlamedaNodes(timeRange *datahub_api.TimeRange) ([]*datahub_v1alpha1.Node, error) {
	alamedaNodes := []*datahub_v1alpha1.Node{}
	nodeRepository := influxdb_repository_cluster_status.NewNodeRepository(&node.InfluxDBConfig)
	entities, err := nodeRepository.ListAlamedaNodes(timeRange)
	if err != nil {
		return alamedaNodes, errors.Wrap(err, "list alameda nodes failed")
	}
	for _, entity := range entities {
		alamedaNodes = append(alamedaNodes, entity.BuildDatahubNode())
	}
	return alamedaNodes, nil
}

func (node *Node) ListNodes(request cluster_status_dao.ListNodesRequest) ([]*datahub_v1alpha1.Node, error) {
	nodes := []*datahub_v1alpha1.Node{}
	nodeRepository := influxdb_repository_cluster_status.NewNodeRepository(&node.InfluxDBConfig)
	entities, err := nodeRepository.ListNodes(request)
	if err != nil {
		return nodes, errors.Wrap(err, "list nodes failed")
	}
	for _, entity := range entities {
		nodes = append(nodes, entity.BuildDatahubNode())
	}
	return nodes, nil
}
