package datahub

import (
	metric_dao "github.com/containers-ai/alameda/datahub/pkg/dao/metric"
	"github.com/containers-ai/alameda/datahub/pkg/dao/prediction"
	"github.com/containers-ai/alameda/datahub/pkg/metric"
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	"github.com/golang/protobuf/ptypes"
)

type daoPodMetricExtended struct {
	*metric_dao.PodMetric
}

func (p daoPodMetricExtended) datahubPodMetric() *datahub_v1alpha1.PodMetric {

	var (
		datahubPodMetric datahub_v1alpha1.PodMetric
	)

	datahubPodMetric = datahub_v1alpha1.PodMetric{
		NamespacedName: &datahub_v1alpha1.NamespacedName{
			Namespace: string(p.Namespace),
			Name:      string(p.PodName),
		},
	}

	for _, containerMetric := range *p.ContainersMetricMap {
		containerMetricExtended := daoContainerMetricExtended{containerMetric}
		datahubContainerMetric := containerMetricExtended.datahubContainerMetric()
		datahubPodMetric.ContainerMetrics = append(datahubPodMetric.ContainerMetrics, datahubContainerMetric)
	}

	return &datahubPodMetric
}

type daoContainerMetricExtended struct {
	*metric_dao.ContainerMetric
}

func (c daoContainerMetricExtended) datahubContainerMetric() *datahub_v1alpha1.ContainerMetric {

	var (
		metricDataChan  = make(chan datahub_v1alpha1.MetricData)
		numOfGoroutines = 0

		datahubContainerMetric datahub_v1alpha1.ContainerMetric
	)

	datahubContainerMetric = datahub_v1alpha1.ContainerMetric{
		Name: string(c.ContainerName),
	}

	for metricType, samples := range c.Metrics {
		if datahubMetricType, exist := metric.TypeToDatahubMetricType[metricType]; exist {
			numOfGoroutines++
			go produceDatahubMetricDataFromSamples(datahubMetricType, samples, metricDataChan)
		}
	}

	for i := 0; i < numOfGoroutines; i++ {
		receivedMetricData := <-metricDataChan
		datahubContainerMetric.MetricData = append(datahubContainerMetric.MetricData, &receivedMetricData)
	}

	return &datahubContainerMetric
}

type daoNodeMetricExtended struct {
	*metric_dao.NodeMetric
}

func (n daoNodeMetricExtended) datahubNodeMetric() *datahub_v1alpha1.NodeMetric {

	var (
		metricDataChan  = make(chan datahub_v1alpha1.MetricData)
		numOfGoroutines = 0

		datahubNodeMetric datahub_v1alpha1.NodeMetric
	)

	datahubNodeMetric = datahub_v1alpha1.NodeMetric{
		Name: n.NodeName,
	}

	for metricType, samples := range n.Metrics {
		if datahubMetricType, exist := metric.TypeToDatahubMetricType[metricType]; exist {
			numOfGoroutines++
			go produceDatahubMetricDataFromSamples(datahubMetricType, samples, metricDataChan)
		}
	}

	for i := 0; i < numOfGoroutines; i++ {
		receivedMetricData := <-metricDataChan
		datahubNodeMetric.MetricData = append(datahubNodeMetric.MetricData, &receivedMetricData)
	}

	return &datahubNodeMetric
}

type daoPtrPodPredictionExtended struct {
	*prediction.PodPrediction
}

func (p daoPtrPodPredictionExtended) datahubPodPrediction() *datahub_v1alpha1.PodPrediction {

	var (
		datahubPodPrediction datahub_v1alpha1.PodPrediction
	)

	datahubPodPrediction = datahub_v1alpha1.PodPrediction{
		NamespacedName: &datahub_v1alpha1.NamespacedName{
			Namespace: string(p.Namespace),
			Name:      string(p.PodName),
		},
	}

	for _, ptrContainerPrediction := range *p.ContainersPredictionMap {
		containerPredictionExtended := daoContainerPredictionExtended{ptrContainerPrediction}
		datahubContainerPrediction := containerPredictionExtended.datahubContainerPrediction()
		datahubPodPrediction.ContainerPredictions = append(datahubPodPrediction.ContainerPredictions, datahubContainerPrediction)
	}

	return &datahubPodPrediction
}

type daoContainerPredictionExtended struct {
	*prediction.ContainerPrediction
}

func (c daoContainerPredictionExtended) datahubContainerPrediction() *datahub_v1alpha1.ContainerPrediction {

	var (
		metricDataChan = make(chan datahub_v1alpha1.MetricData)
		numOfGoroutine = 0

		datahubContainerPrediction datahub_v1alpha1.ContainerPrediction
	)

	datahubContainerPrediction = datahub_v1alpha1.ContainerPrediction{
		Name: string(c.ContainerName),
	}

	for metricType, samples := range c.PredictionsRaw {
		if datahubMetricType, exist := metric.TypeToDatahubMetricType[metricType]; exist {
			numOfGoroutine++
			go produceDatahubMetricDataFromSamples(datahubMetricType, samples, metricDataChan)
		}
	}

	for i := 0; i < numOfGoroutine; i++ {
		receivedPredictionData := <-metricDataChan
		datahubContainerPrediction.PredictedRawData = append(datahubContainerPrediction.PredictedRawData, &receivedPredictionData)
	}

	return &datahubContainerPrediction
}

type daoPtrNodePredictionExtended struct {
	*prediction.NodePrediction
}

func (d daoPtrNodePredictionExtended) datahubNodePrediction() *datahub_v1alpha1.NodePrediction {

	var (
		metricDataChan = make(chan datahub_v1alpha1.MetricData)
		numOfGoroutine = 0

		datahubNodePrediction datahub_v1alpha1.NodePrediction
	)

	datahubNodePrediction = datahub_v1alpha1.NodePrediction{
		Name:        string(d.NodeName),
		IsScheduled: d.IsScheduled,
	}

	for metricType, samples := range d.Predictions {
		if datahubMetricType, exist := metric.TypeToDatahubMetricType[metricType]; exist {
			numOfGoroutine++
			go produceDatahubMetricDataFromSamples(datahubMetricType, samples, metricDataChan)
		}
	}

	for i := 0; i < numOfGoroutine; i++ {
		receivedPredictionData := <-metricDataChan
		datahubNodePrediction.PredictedRawData = append(datahubNodePrediction.PredictedRawData, &receivedPredictionData)
	}

	return &datahubNodePrediction
}

type daoPtrNodesPredictionMapExtended struct {
	*prediction.NodesPredictionMap
}

func (d daoPtrNodesPredictionMapExtended) datahubNodePredictions() []*datahub_v1alpha1.NodePrediction {

	var (
		datahubNodePredictions = make([]*datahub_v1alpha1.NodePrediction, 0)
	)

	for _, ptrIsScheduledNodePredictionMap := range *d.NodesPredictionMap {

		if ptrScheduledNodePrediction, exist := (*ptrIsScheduledNodePredictionMap)[true]; exist {

			scheduledNodePredictionExtended := daoPtrNodePredictionExtended{ptrScheduledNodePrediction}
			sechduledDatahubNodePrediction := scheduledNodePredictionExtended.datahubNodePrediction()
			datahubNodePredictions = append(datahubNodePredictions, sechduledDatahubNodePrediction)
		}

		if noneScheduledNodePrediction, exist := (*ptrIsScheduledNodePredictionMap)[false]; exist {

			noneScheduledNodePredictionExtended := daoPtrNodePredictionExtended{noneScheduledNodePrediction}
			noneSechduledDatahubNodePrediction := noneScheduledNodePredictionExtended.datahubNodePrediction()
			datahubNodePredictions = append(datahubNodePredictions, noneSechduledDatahubNodePrediction)
		}
	}

	return datahubNodePredictions
}

func produceDatahubMetricDataFromSamples(metricType datahub_v1alpha1.MetricType, samples []metric.Sample, MetricDataChan chan<- datahub_v1alpha1.MetricData) {

	var (
		datahubMetricData datahub_v1alpha1.MetricData
	)

	datahubMetricData = datahub_v1alpha1.MetricData{
		MetricType: metricType,
	}

	for _, sample := range samples {

		// TODO: Send error to caller
		googleTimestamp, err := ptypes.TimestampProto(sample.Timestamp)
		if err != nil {
			googleTimestamp = nil
		}

		datahubSample := datahub_v1alpha1.Sample{Time: googleTimestamp, NumValue: sample.Value}
		datahubMetricData.Data = append(datahubMetricData.Data, &datahubSample)
	}

	MetricDataChan <- datahubMetricData
}
