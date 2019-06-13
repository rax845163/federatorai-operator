package containerCPUUsagePercentage

import (
	"fmt"
	metric_dao "github.com/containers-ai/alameda/datahub/pkg/dao/metric"
	"github.com/containers-ai/alameda/datahub/pkg/metric"
	"github.com/containers-ai/alameda/datahub/pkg/repository/prometheus"
	"github.com/containers-ai/alameda/pkg/utils/log"
	"math"
	"strconv"
)

const (
	// MetricName Metric name to query from prometheus
	MetricName = "namespace_pod_name_container_name:container_cpu_usage_seconds_total:sum_rate"
	// NamespaceLabel Namespace label name in the metric
	NamespaceLabel = "namespace"
	// PodLabelName pod label name in the metric
	PodLabelName = "pod_name"
	// ContainerLabel container label name in the metric
	ContainerLabel = "container_name"
)

var (
	scope = log.RegisterScope("prometheus", "metrics repository", 0)
)

// Entity Container cpu usage percentage entity
type Entity struct {
	PrometheusEntity prometheus.Entity

	Namespace     string
	PodName       string
	ContainerName string
	Samples       []metric.Sample
}

// NewEntityFromPrometheusEntity New entity with field value assigned from prometheus entity
func NewEntityFromPrometheusEntity(e prometheus.Entity) Entity {

	var (
		samples []metric.Sample
	)

	samples = make([]metric.Sample, 0)

	for _, value := range e.Values {
		v := "0"
		if s, err := strconv.ParseFloat(value.SampleValue, 64); err == nil {
			v = fmt.Sprintf("%d", int(math.Ceil(s*1000)))
			if v == "0" {
				v = "1"
			}
		} else {
			scope.Errorf("containerCPUUsagePercentage.NewEntityFromPrometheusEntity: %s", err.Error())
		}
		sample := metric.Sample{
			Timestamp: value.UnixTime,
			Value:     v,
		}
		samples = append(samples, sample)
	}

	return Entity{
		PrometheusEntity: e,
		Namespace:        e.Labels[NamespaceLabel],
		PodName:          e.Labels[PodLabelName],
		ContainerName:    e.Labels[ContainerLabel],
		Samples:          samples,
	}
}

// ContainerMetric Build ContainerMetric base on entity properties
func (e *Entity) ContainerMetric() metric_dao.ContainerMetric {

	var (
		containerMetric metric_dao.ContainerMetric
	)

	containerMetric = metric_dao.ContainerMetric{
		Namespace:     e.Namespace,
		PodName:       e.PodName,
		ContainerName: e.ContainerName,
		Metrics: map[metric.ContainerMetricType][]metric.Sample{
			metric.TypeContainerCPUUsageSecondsPercentage: e.Samples,
		},
	}

	return containerMetric
}
