package metric

import (
	"fmt"
	"time"

	"github.com/containers-ai/alameda/datahub/pkg/entity/prometheus/nodeMemoryBytesTotal"
	"github.com/containers-ai/alameda/datahub/pkg/entity/prometheus/nodeMemoryUtilization"
	"github.com/containers-ai/alameda/datahub/pkg/repository/prometheus"
	"github.com/containers-ai/alameda/pkg/utils/log"
	"github.com/pkg/errors"
)

var (
	scope = log.RegisterScope("node memory usage bytes", "node memory usage bytes", 0)
)

// NodeMemoryUsageBytesRepository Repository to access metric from prometheus
type NodeMemoryUsageBytesRepository struct {
	PrometheusConfig prometheus.Config
}

// NewNodeMemoryUsageBytesRepositoryWithConfig New node cpu usage percentage repository with prometheus configuration
func NewNodeMemoryUsageBytesRepositoryWithConfig(cfg prometheus.Config) NodeMemoryUsageBytesRepository {
	return NodeMemoryUsageBytesRepository{PrometheusConfig: cfg}
}

// ListMetricsByNodeName Provide metrics from response of querying request contain namespace, pod_name and default labels
func (n NodeMemoryUsageBytesRepository) ListMetricsByNodeName(nodeName string, options ...Option) ([]prometheus.Entity, error) {

	var (
		err error

		prometheusClient *prometheus.Prometheus

		nodeMemoryBytesTotalQueryExpression    string
		nodeMemoryBytesTotalMetricName         string
		nodeMemoryBytesTotalQueryLabelsString  string
		nodeMemoryUtilizationQueryExpression   string
		nodeMemoryUtilizationMetricName        string
		nodeMemoryUtilizationQueryLabelsString string
		queryExpression                        string

		response prometheus.Response

		entities []prometheus.Entity
	)

	prometheusClient, err = prometheus.New(n.PrometheusConfig)
	if err != nil {
		return entities, errors.Wrap(err, "list node memory usage by node name failed")
	}

	opt := buildDefaultOptions()
	for _, option := range options {
		option(&opt)
	}
	stepTimeInSeconds := int64(opt.stepTime.Nanoseconds() / int64(time.Second))

	nodeMemoryBytesTotalMetricName = nodeMemoryBytesTotal.MetricName
	nodeMemoryBytesTotalQueryLabelsString = n.buildNodeMemoryBytesTotalQueryLabelsStringByNodeName(nodeName)
	if nodeMemoryBytesTotalQueryLabelsString != "" {
		nodeMemoryBytesTotalQueryExpression = fmt.Sprintf("%s{%s}", nodeMemoryBytesTotalMetricName, nodeMemoryBytesTotalQueryLabelsString)
	} else {
		nodeMemoryBytesTotalQueryExpression = fmt.Sprintf("%s", nodeMemoryBytesTotalMetricName)
	}
	nodeMemoryBytesTotalQueryExpression, err = wrapQueryExpressionWithAggregationOverTimeFunction(nodeMemoryBytesTotalQueryExpression, opt.aggregateOverTimeFunc, stepTimeInSeconds)
	if err != nil {
		return entities, errors.Wrap(err, "list node memory usage metrics by node name failed")
	}

	nodeMemoryUtilizationMetricName = nodeMemoryUtilization.MetricName
	nodeMemoryUtilizationQueryLabelsString = n.buildNodeMemoryUtilizationQueryLabelsStringByNodeName(nodeName)
	if nodeMemoryUtilizationQueryLabelsString != "" {
		nodeMemoryUtilizationQueryExpression = fmt.Sprintf("%s{%s}", nodeMemoryUtilizationMetricName, nodeMemoryUtilizationQueryLabelsString)
	} else {
		nodeMemoryUtilizationQueryExpression = fmt.Sprintf("%s", nodeMemoryUtilizationMetricName)
	}
	nodeMemoryUtilizationQueryExpression, err = wrapQueryExpressionWithAggregationOverTimeFunction(nodeMemoryUtilizationQueryExpression, opt.aggregateOverTimeFunc, stepTimeInSeconds)
	if err != nil {
		return entities, errors.Wrap(err, "list node memory usage metrics by node name failed")
	}

	queryExpression = fmt.Sprintf("%s * %s", nodeMemoryBytesTotalQueryExpression, nodeMemoryUtilizationQueryExpression)

	response, err = prometheusClient.QueryRange(queryExpression, opt.startTime, opt.endTime, opt.stepTime)
	if err != nil {
		return entities, errors.Wrap(err, "list node memory bytes total by node name failed")
	} else if response.Status != prometheus.StatusSuccess {
		return entities, errors.Errorf("list node memory bytes total by node name failed: receive error response from prometheus: %s", response.Error)
	}

	entities, err = response.GetEntitis()
	if err != nil {
		return entities, errors.Wrap(err, "list node memory usage by node name failed")
	}

	return entities, nil
}

func (n NodeMemoryUsageBytesRepository) buildNodeMemoryBytesTotalQueryLabelsStringByNodeName(nodeName string) string {

	var (
		queryLabelsString = ""
	)

	if nodeName != "" {
		queryLabelsString += fmt.Sprintf(`%s = "%s"`, nodeMemoryBytesTotal.NodeLabel, nodeName)
	}

	return queryLabelsString
}

func (n NodeMemoryUsageBytesRepository) buildNodeMemoryUtilizationQueryLabelsStringByNodeName(nodeName string) string {

	var (
		queryLabelsString = ""
	)

	if nodeName != "" {
		queryLabelsString += fmt.Sprintf(`%s = "%s"`, nodeMemoryUtilization.NodeLabel, nodeName)
	}

	return queryLabelsString
}
