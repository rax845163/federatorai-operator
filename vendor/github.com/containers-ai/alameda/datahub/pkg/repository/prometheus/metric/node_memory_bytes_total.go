package metric

import (
	"fmt"

	"github.com/containers-ai/alameda/datahub/pkg/entity/prometheus/nodeMemoryBytesTotal"
	"github.com/containers-ai/alameda/datahub/pkg/repository/prometheus"
	"github.com/pkg/errors"
)

// NodeMemoryBytesTotalRepository Repository to access metric from prometheus
type NodeMemoryBytesTotalRepository struct {
	PrometheusConfig prometheus.Config
}

// NewNodeMemoryBytesTotalRepositoryWithConfig New node cpu utilization percentage repository with prometheus configuration
func NewNodeMemoryBytesTotalRepositoryWithConfig(cfg prometheus.Config) NodeMemoryBytesTotalRepository {
	return NodeMemoryBytesTotalRepository{PrometheusConfig: cfg}
}

func (n NodeMemoryBytesTotalRepository) ListMetricsByNodeName(nodeName string, options ...Option) ([]prometheus.Entity, error) {

	var (
		err error

		prometheusClient *prometheus.Prometheus

		nodeMemoryBytesTotalMetricName        string
		nodeMemoryBytesTotalQueryLabelsString string
		queryExpression                       string

		response prometheus.Response

		entities []prometheus.Entity
	)

	prometheusClient, err = prometheus.New(n.PrometheusConfig)
	if err != nil {
		return entities, errors.Wrap(err, "list node memory utilization by node name failed")
	}

	opt := buildDefaultOptions()
	for _, option := range options {
		option(&opt)
	}

	nodeMemoryBytesTotalMetricName = nodeMemoryBytesTotal.MetricName
	nodeMemoryBytesTotalQueryLabelsString = n.buildNodeMemoryBytesTotalQueryLabelsStringByNodeName(nodeName)

	if nodeMemoryBytesTotalQueryLabelsString != "" {
		queryExpression = fmt.Sprintf("%s{%s}", nodeMemoryBytesTotalMetricName, nodeMemoryBytesTotalQueryLabelsString)
	} else {
		queryExpression = fmt.Sprintf("%s", nodeMemoryBytesTotalMetricName)
	}

	response, err = prometheusClient.QueryRange(queryExpression, opt.startTime, opt.endTime, opt.stepTime)
	if err != nil {
		return entities, errors.Wrap(err, "list node memory bytes total by node name failed")
	} else if response.Status != prometheus.StatusSuccess {
		return entities, errors.Errorf("list node memory bytes total by node name failed: receive error response from prometheus: %s", response.Error)
	}

	entities, err = response.GetEntitis()
	if err != nil {
		return entities, errors.Wrap(err, "list node memory bytes total by node name failed")
	}

	return entities, nil
}

func (n NodeMemoryBytesTotalRepository) buildNodeMemoryBytesTotalQueryLabelsStringByNodeName(nodeName string) string {

	var (
		queryLabelsString = ""
	)

	if nodeName != "" {
		queryLabelsString += fmt.Sprintf(`%s = "%s"`, nodeMemoryBytesTotal.NodeLabel, nodeName)
	}

	return queryLabelsString
}
