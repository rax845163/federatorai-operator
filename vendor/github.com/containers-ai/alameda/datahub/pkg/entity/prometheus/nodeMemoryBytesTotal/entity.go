package nodeMemoryBytesTotal

import (
	metric_dao "github.com/containers-ai/alameda/datahub/pkg/dao/metric"
	"github.com/containers-ai/alameda/datahub/pkg/metric"
	"github.com/containers-ai/alameda/datahub/pkg/repository/prometheus"
)

const (
	// MetricName Metric name to query from prometheus
	MetricName = "node:node_memory_bytes_total:sum"
	// NodeLabel Node label name in the metric
	NodeLabel = "node"
)

// Entity Node total memory entity
type Entity struct {
	PrometheusEntity prometheus.Entity

	NodeName string
	Samples  []metric.Sample
}

// NewEntityFromPrometheusEntity New entity with field value assigned from prometheus entity
func NewEntityFromPrometheusEntity(e prometheus.Entity) Entity {

	var (
		samples []metric.Sample
	)

	samples = make([]metric.Sample, 0)

	for _, value := range e.Values {
		sample := metric.Sample{
			Timestamp: value.UnixTime,
			Value:     value.SampleValue,
		}
		samples = append(samples, sample)
	}

	return Entity{
		PrometheusEntity: e,
		NodeName:         e.Labels[NodeLabel],
		Samples:          samples,
	}
}

// NodeMetric Build NodeMetric base on entity properties
func (e *Entity) NodeMetric() metric_dao.NodeMetric {

	var (
		nodeMetric metric_dao.NodeMetric
	)

	nodeMetric = metric_dao.NodeMetric{
		NodeName: e.NodeName,
		Metrics: map[metric.NodeMetricType][]metric.Sample{
			metric.TypeNodeMemoryTotalBytes: e.Samples,
		},
	}

	return nodeMetric
}
