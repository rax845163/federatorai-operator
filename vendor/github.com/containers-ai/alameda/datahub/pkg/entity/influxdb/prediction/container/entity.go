package container

import (
	"time"

	"github.com/containers-ai/alameda/datahub/pkg/dao/prediction"
	"github.com/containers-ai/alameda/datahub/pkg/metric"
	"github.com/containers-ai/alameda/datahub/pkg/utils"
)

type Field = string
type Tag = string
type MetricType = string

const (
	Time        Tag = "time"
	Namespace   Tag = "namespace"
	PodName     Tag = "pod_name"
	Name        Tag = "name"
	Metric      Tag = "metric"
	Granularity Tag = "granularity"
	Kind        Tag = "kind"

	Value Field = "value"
)

var (
	// Tags Tags' name in influxdb
	Tags = []Tag{Namespace, PodName, Name, Metric, Granularity, Kind}
	// Fields Fields' name in influxdb
	Fields = []Field{Value}
	// MetricTypeCPUUsage Enum of tag "metric"
	MetricTypeCPUUsage MetricType = "cpu_usage_seconds_percentage"
	// MetricTypeMemoryUsage Enum of tag "metric"
	MetricTypeMemoryUsage MetricType = "memory_usage_bytes"

	// LocalMetricTypeToPkgMetricType Convert local package metric type to package alameda.datahub.metric.NodeMetricType
	LocalMetricTypeToPkgMetricType = map[MetricType]metric.NodeMetricType{
		MetricTypeCPUUsage:    metric.TypeContainerCPUUsageSecondsPercentage,
		MetricTypeMemoryUsage: metric.TypeContainerMemoryUsageBytes,
	}

	// PkgMetricTypeToLocalMetricType Convert package alameda.datahub.metric.NodeMetricType to local package metric type
	PkgMetricTypeToLocalMetricType = map[metric.NodeMetricType]MetricType{
		metric.TypeContainerCPUUsageSecondsPercentage: MetricTypeCPUUsage,
		metric.TypeContainerMemoryUsageBytes:          MetricTypeMemoryUsage,
	}
)

// Entity Container prediction entity in influxDB
type Entity struct {
	Timestamp time.Time

	Namespace   *string
	PodName     *string
	Name        *string
	Metric      MetricType
	Value       *string
	Kind        string
	Granularity *string
}

// NewEntityFromMap Build entity from map
func NewEntityFromMap(data map[string]string) Entity {

	// TODO: log error
	tempTimestamp, _ := utils.ParseTime(data[Time])

	entity := Entity{
		Timestamp: tempTimestamp,
	}

	if namespace, exist := data[Namespace]; exist {
		entity.Namespace = &namespace
	}

	if podName, exist := data[PodName]; exist {
		entity.PodName = &podName
	}

	if name, exist := data[Name]; exist {
		entity.Name = &name
	}

	if metric, exist := data[Metric]; exist {
		entity.Metric = metric
	}

	if value, exist := data[Value]; exist {
		entity.Value = &value
	}

	if kind, exist := data[Kind]; exist {
		entity.Kind = kind
	}

	if granularity, exist := data[Granularity]; exist {
		entity.Granularity = &granularity
	}

	return entity
}

// ContainerPrediction Create container prediction base on entity
func (e Entity) ContainerPrediction() prediction.ContainerPrediction {

	var (
		samples             []metric.Sample
		containerPrediction prediction.ContainerPrediction
	)

	samples = append(samples, metric.Sample{Timestamp: e.Timestamp, Value: *e.Value})

	containerPrediction = prediction.ContainerPrediction{
		Namespace:        *e.Namespace,
		PodName:          *e.PodName,
		ContainerName:    *e.Name,
		PredictionsRaw:   map[metric.ContainerMetricType][]metric.Sample{},
		PredictionsUpper: map[metric.ContainerMetricType][]metric.Sample{},
		PredictionsLower: map[metric.ContainerMetricType][]metric.Sample{},
	}

	//metricType := LocalMetricTypeToPkgMetricType[*e.Metric]
	metricType := e.Metric

	if e.Kind == metric.ContainerMetricKindUpperbound {
		containerPrediction.PredictionsUpper[metricType] = samples
	} else if e.Kind == metric.ContainerMetricKindLowerbound {
		containerPrediction.PredictionsLower[metricType] = samples
	} else {
		containerPrediction.PredictionsRaw[metricType] = samples
	}

	return containerPrediction
}
