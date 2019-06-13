package metric

import (
	"fmt"
	"time"

	"github.com/containers-ai/alameda/datahub/pkg/entity/prometheus/nodeCPUUsagePercentage"
	"github.com/containers-ai/alameda/datahub/pkg/repository/prometheus"
	"github.com/pkg/errors"
)

// NodeCPUUsagePercentageRepository Repository to access metric node:node_cpu_utilisation:avg1m from prometheus
type NodeCPUUsagePercentageRepository struct {
	PrometheusConfig prometheus.Config
}

// NewNodeCPUUsagePercentageRepositoryWithConfig New node cpu usage percentage repository with prometheus configuration
func NewNodeCPUUsagePercentageRepositoryWithConfig(cfg prometheus.Config) NodeCPUUsagePercentageRepository {
	return NodeCPUUsagePercentageRepository{PrometheusConfig: cfg}
}

// ListMetricsByPodNamespacedName Provide metrics from response of querying request contain namespace, pod_name and default labels
func (n NodeCPUUsagePercentageRepository) ListMetricsByNodeName(nodeName string, options ...Option) ([]prometheus.Entity, error) {

	var (
		err error

		prometheusClient *prometheus.Prometheus

		metricName        string
		queryLabelsString string
		queryExpression   string

		response prometheus.Response

		entities []prometheus.Entity
	)

	prometheusClient, err = prometheus.New(n.PrometheusConfig)
	if err != nil {
		return entities, errors.Wrap(err, "list node cpu usage metrics by node name failed")
	}

	opt := buildDefaultOptions()
	for _, option := range options {
		option(&opt)
	}

	metricName = nodeCPUUsagePercentage.MetricName
	queryLabelsString = n.buildQueryLabelsStringByNodeName(nodeName)

	if queryLabelsString != "" {
		queryExpression = fmt.Sprintf("%s{%s}", metricName, queryLabelsString)
	} else {
		queryExpression = fmt.Sprintf("%s", metricName)
	}

	stepTimeInSeconds := int64(opt.stepTime.Nanoseconds() / int64(time.Second))
	queryExpression, err = wrapQueryExpressionWithAggregationOverTimeFunction(queryExpression, opt.aggregateOverTimeFunc, stepTimeInSeconds)
	if err != nil {
		return entities, errors.Wrap(err, "list node cpu usage metrics by node name failed")
	}

	response, err = prometheusClient.QueryRange(queryExpression, opt.startTime, opt.endTime, opt.stepTime)
	if err != nil {
		return entities, errors.Wrap(err, "list node cpu usage metrics by node name failed")
	} else if response.Status != prometheus.StatusSuccess {
		return entities, errors.Errorf("list node cpu usage metrics by node name failed: receive error response from prometheus: %s", response.Error)
	}

	entities, err = response.GetEntitis()
	if err != nil {
		return entities, errors.Wrap(err, "list node cpu usage metrics by node name failed")
	}

	return entities, nil
}

func (n NodeCPUUsagePercentageRepository) buildQueryLabelsStringByNodeName(nodeName string) string {

	var (
		queryLabelsString = ""
	)

	if nodeName != "" {
		queryLabelsString += fmt.Sprintf(`%s = "%s"`, nodeCPUUsagePercentage.NodeLabel, nodeName)
	}

	return queryLabelsString
}
