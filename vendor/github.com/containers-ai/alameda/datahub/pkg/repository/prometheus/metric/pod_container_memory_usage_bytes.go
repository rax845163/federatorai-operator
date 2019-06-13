package metric

import (
	"fmt"
	"time"

	"github.com/containers-ai/alameda/datahub/pkg/entity/prometheus/containerMemoryUsageBytes"
	"github.com/containers-ai/alameda/datahub/pkg/repository/prometheus"
	"github.com/pkg/errors"
)

// PodContainerMemoryUsageBytesRepository Repository to access metric container_memory_usage_bytes from prometheus
type PodContainerMemoryUsageBytesRepository struct {
	PrometheusConfig prometheus.Config
}

// NewPodContainerMemoryUsageBytesRepositoryWithConfig New pod container memory usage bytes repository with prometheus configuration
func NewPodContainerMemoryUsageBytesRepositoryWithConfig(cfg prometheus.Config) PodContainerMemoryUsageBytesRepository {
	return PodContainerMemoryUsageBytesRepository{PrometheusConfig: cfg}
}

// ListMetricsByPodNamespacedName Provide metrics from response of querying request contain namespace, pod_name and default labels
func (c PodContainerMemoryUsageBytesRepository) ListMetricsByPodNamespacedName(namespace string, podName string, options ...Option) ([]prometheus.Entity, error) {

	var (
		err error

		prometheusClient *prometheus.Prometheus

		metricName        string
		queryLabelsString string
		queryExpression   string

		response prometheus.Response

		entities []prometheus.Entity
	)

	prometheusClient, err = prometheus.New(c.PrometheusConfig)
	if err != nil {
		return entities, errors.Wrap(err, "list pod container memory usage metrics failed")
	}

	opt := buildDefaultOptions()
	for _, option := range options {
		option(&opt)
	}

	metricName = containerMemoryUsageBytes.MetricName
	queryLabelsString = c.buildQueryLabelsStringByNamespaceAndPodName(namespace, podName)

	if queryLabelsString != "" {
		queryExpression = fmt.Sprintf("%s{%s}", metricName, queryLabelsString)
	} else {
		queryExpression = fmt.Sprintf("%s", metricName)
	}

	stepTimeInSeconds := int64(opt.stepTime.Nanoseconds() / int64(time.Second))
	queryExpression, err = wrapQueryExpressionWithAggregationOverTimeFunction(queryExpression, opt.aggregateOverTimeFunc, stepTimeInSeconds)
	if err != nil {
		return entities, errors.Wrap(err, "list pod container memory usage metric by namespaced name failed")
	}

	response, err = prometheusClient.QueryRange(queryExpression, opt.startTime, opt.endTime, opt.stepTime)
	if err != nil {
		return entities, errors.Wrap(err, "list pod container memory usage metrics failed")
	} else if response.Status != prometheus.StatusSuccess {
		return entities, errors.Errorf("list pod container memory usage metrics failed: receive error response from prometheus: %s", response.Error)
	}

	entities, err = response.GetEntitis()
	if err != nil {
		return entities, errors.Wrap(err, "list pod container memory usage metrics failed")
	}

	return entities, nil
}

func (c PodContainerMemoryUsageBytesRepository) buildDefaultQueryLabelsString() string {

	var queryLabelsString = ""

	queryLabelsString += fmt.Sprintf(`%s != "" ,`, containerMemoryUsageBytes.PodLabelName)
	queryLabelsString += fmt.Sprintf(`%s != "" ,`, containerMemoryUsageBytes.ContainerLabel)
	queryLabelsString += fmt.Sprintf(`%s != "POD"`, containerMemoryUsageBytes.ContainerLabel)

	return queryLabelsString
}

func (c PodContainerMemoryUsageBytesRepository) buildQueryLabelsStringByNamespaceAndPodName(namespace string, podName string) string {

	var (
		queryLabelsString = c.buildDefaultQueryLabelsString()
	)

	if namespace != "" {
		queryLabelsString += fmt.Sprintf(`,%s = "%s"`, containerMemoryUsageBytes.NamespaceLabel, namespace)
	}

	if podName != "" {
		queryLabelsString += fmt.Sprintf(`,%s = "%s"`, containerMemoryUsageBytes.PodLabelName, podName)
	}

	return queryLabelsString
}
