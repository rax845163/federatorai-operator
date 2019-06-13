package clusterstatus

import (
	datahub_api "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
)

// Node provides node measurement operations
type NodeOperation interface {
	RegisterAlamedaNodes([]*datahub_api.Node) error
	DeregisterAlamedaNodes([]*datahub_api.Node) error
	ListAlamedaNodes(timeRange *datahub_api.TimeRange) ([]*datahub_api.Node, error)
	ListNodes(ListNodesRequest) ([]*datahub_api.Node, error)
}

type ListNodesRequest struct {
	NodeNames []string
	InCluster bool
}
