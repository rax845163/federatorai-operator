package metric

import (
	"fmt"
	"time"

	"github.com/containers-ai/alameda/datahub/pkg/entity/prometheus/nodeMemoryUtilization"
	"github.com/containers-ai/alameda/datahub/pkg/repository/prometheus"
	"github.com/pkg/errors"
)

// NodeMemoryUtilizationRepository Repository to access metric from prometheus
type NodeMemoryUtilizationRepository struct {
	PrometheusConfig prometheus.Config
}

// NewNodeMemoryUtilizationRepositoryWithConfig New node cpu utilization percentage repository with prometheus configuration
func NewNodeMemoryUtilizationRepositoryWithConfig(cfg prometheus.Config) NodeMemoryUtilizationRepository {
	return NodeMemoryUtilizationRepository{PrometheusConfig: cfg}
}

// ListMetricsByNodeName Provide metrics from response of querying request contain namespace, pod_name and default labels
func (n NodeMemoryUtilizationRepository) ListMetricsByNodeName(nodeName string, options ...Option) ([]prometheus.Entity, error) {

	var (
		err error

		prometheusClient *prometheus.Prometheus

		nodeMemoryUtilizationMetricName   string
		nodeMemoryUtilizationLabelsString string
		queryExpression                   string

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

	nodeMemoryUtilizationMetricName = nodeMemoryUtilization.MetricName
	nodeMemoryUtilizationLabelsString = n.buildNodeMemoryUtilizationQueryLabelsStringByNodeName(nodeName)

	if nodeMemoryUtilizationLabelsString != "" {
		queryExpression = fmt.Sprintf("%s{%s}", nodeMemoryUtilizationMetricName, nodeMemoryUtilizationLabelsString)
	} else {
		queryExpression = fmt.Sprintf("%s", nodeMemoryUtilizationMetricName)
	}

	stepTimeInSeconds := int64(opt.stepTime.Nanoseconds() / int64(time.Second))
	queryExpression, err = wrapQueryExpressionWithAggregationOverTimeFunction(queryExpression, opt.aggregateOverTimeFunc, stepTimeInSeconds)
	if err != nil {
		return entities, errors.Wrap(err, "list node memory utilization by node name failed")
	}

	response, err = prometheusClient.QueryRange(queryExpression, opt.startTime, opt.endTime, opt.stepTime)
	if err != nil {
		return entities, errors.Wrap(err, "list node memory utilization by node name failed")
	} else if response.Status != prometheus.StatusSuccess {
		return entities, errors.Errorf("list node memory utilization by node name failed: receive error response from prometheus: %s", response.Error)
	}

	entities, err = response.GetEntitis()
	if err != nil {
		return entities, errors.Wrap(err, "list node memory utilization by node name failed")
	}

	return entities, nil
}

func (n NodeMemoryUtilizationRepository) buildNodeMemoryUtilizationQueryLabelsStringByNodeName(nodeName string) string {

	var (
		queryLabelsString = ""
	)

	if nodeName != "" {
		queryLabelsString += fmt.Sprintf(`%s = "%s"`, nodeMemoryUtilization.NodeLabel, nodeName)
	}

	return queryLabelsString
}
