package prediction

import (
	"fmt"

	prediction_dao "github.com/containers-ai/alameda/datahub/pkg/dao/prediction"
	container_entity "github.com/containers-ai/alameda/datahub/pkg/entity/influxdb/prediction/container"
	"github.com/containers-ai/alameda/datahub/pkg/metric"
	"github.com/containers-ai/alameda/datahub/pkg/repository/influxdb"
	datahub_utils "github.com/containers-ai/alameda/datahub/pkg/utils"
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	"github.com/golang/protobuf/ptypes"
	influxdb_client "github.com/influxdata/influxdb/client/v2"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

// ContainerRepository Repository to access containers' prediction data
type ContainerRepository struct {
	influxDB *influxdb.InfluxDBRepository
}

// NewContainerRepositoryWithConfig New container repository with influxDB configuration
func NewContainerRepositoryWithConfig(influxDBCfg influxdb.Config) *ContainerRepository {
	return &ContainerRepository{
		influxDB: &influxdb.InfluxDBRepository{
			Address:  influxDBCfg.Address,
			Username: influxDBCfg.Username,
			Password: influxDBCfg.Password,
		},
	}
}

func (r *ContainerRepository) CreateContainerPrediction(in *datahub_v1alpha1.CreatePodPredictionsRequest) error {

	points := make([]*influxdb_client.Point, 0)

	for _, podPrediction := range in.GetPodPredictions() {
		podNamespace := podPrediction.GetNamespacedName().GetNamespace()
		podName := podPrediction.GetNamespacedName().GetName()
		for _, containerPrediction := range podPrediction.GetContainerPredictions() {
			containerName := containerPrediction.GetName()

			r.appendMetricDataToPoints(metric.ContainerMetricKindRaw, containerPrediction.GetPredictedRawData(), &points, podNamespace, podName, containerName)
			r.appendMetricDataToPoints(metric.ContainerMetricKindUpperbound, containerPrediction.GetPredictedUpperboundData(), &points, podNamespace, podName, containerName)
			r.appendMetricDataToPoints(metric.ContainerMetricKindLowerbound, containerPrediction.GetPredictedLowerboundData(), &points, podNamespace, podName, containerName)
		}
	}

	err := r.influxDB.WritePoints(points, influxdb_client.BatchPointsConfig{
		Database: string(influxdb.Prediction),
	})
	if err != nil {
		return errors.Wrap(err, "create container prediction failed")
	}

	return nil
}

func (r *ContainerRepository) appendMetricDataToPoints(kind metric.ContainerMetricKind, metricDataList []*datahub_v1alpha1.MetricData, points *[]*influxdb_client.Point, podNamespace string, podName string, containerName string) error {
	for _, metricData := range metricDataList {
		metricType := ""
		switch metricData.GetMetricType() {
		case datahub_v1alpha1.MetricType_CPU_USAGE_SECONDS_PERCENTAGE:
			metricType = metric.TypeContainerCPUUsageSecondsPercentage
		case datahub_v1alpha1.MetricType_MEMORY_USAGE_BYTES:
			metricType = metric.TypeContainerMemoryUsageBytes
		}

		if metricType == "" {
			return errors.New("No corresponding metricType")
		}

		granularity := metricData.GetGranularity()
		if granularity == 0 {
			granularity = 30
		}

		for _, data := range metricData.GetData() {
			tempTimeSeconds := data.GetTime().Seconds
			value := data.GetNumValue()
			valueWithoutFraction := strings.Split(value, ".")[0]
			valueInFloat64, err := datahub_utils.StringToFloat64(valueWithoutFraction)
			if err != nil {
				return errors.Wrap(err, "new influxdb data point failed")
			}

			tags := map[string]string{
				container_entity.Namespace:   podNamespace,
				container_entity.PodName:     podName,
				container_entity.Name:        containerName,
				container_entity.Metric:      metricType,
				container_entity.Kind:        kind,
				container_entity.Granularity: strconv.FormatInt(granularity, 10),
			}
			fields := map[string]interface{}{
				container_entity.Value: valueInFloat64,
			}
			point, err := influxdb_client.NewPoint(string(Container), tags, fields, time.Unix(tempTimeSeconds, 0))
			if err != nil {
				return errors.Wrap(err, "new influxdb data point failed")
			}
			*points = append(*points, point)
		}
	}

	return nil
}

// ListContainerPredictionsByRequest list containers' prediction from influxDB
func (r *ContainerRepository) ListContainerPredictionsByRequest(request prediction_dao.ListPodPredictionsRequest) ([]*datahub_v1alpha1.PodPrediction, error) {
	whereClause := r.buildInfluxQLWhereClauseFromRequest(request)
	influxdbStatement := influxdb.Statement{
		Measurement: Container,
		WhereClause: whereClause,
		GroupByTags: []string{container_entity.Namespace, container_entity.PodName, container_entity.Name, container_entity.Metric, container_entity.Kind, container_entity.Granularity},
	}

	queryCondition := influxdb.QueryCondition{
		StartTime:      request.QueryCondition.StartTime,
		EndTime:        request.QueryCondition.EndTime,
		StepTime:       request.QueryCondition.StepTime,
		TimestampOrder: request.QueryCondition.TimestampOrder,
		Limit:          request.QueryCondition.Limit,
	}
	influxdbStatement.AppendTimeConditionIntoWhereClause(queryCondition)
	influxdbStatement.SetLimitClauseFromQueryCondition(queryCondition)
	influxdbStatement.SetOrderClauseFromQueryCondition(queryCondition)
	cmd := influxdbStatement.BuildQueryCmd()

	results, err := r.influxDB.QueryDB(cmd, string(influxdb.Prediction))
	if err != nil {
		return []*datahub_v1alpha1.PodPrediction{}, errors.Wrap(err, "list container prediction failed")
	}

	rows := influxdb.PackMap(results)
	podPredictions := r.getPodPredictionsFromInfluxRows(rows)

	return podPredictions, nil
}

func (r *ContainerRepository) getPodPredictionsFromInfluxRows(rows []*influxdb.InfluxDBRow) []*datahub_v1alpha1.PodPrediction {
	podMap := map[string]*datahub_v1alpha1.PodPrediction{}
	podContainerMap := map[string]*datahub_v1alpha1.ContainerPrediction{}
	podContainerKindMetricMap := map[string]*datahub_v1alpha1.MetricData{}
	podContainerKindMetricSampleMap := map[string][]*datahub_v1alpha1.Sample{}

	for _, row := range rows {
		namespace := row.Tags[container_entity.Namespace]
		podName := row.Tags[container_entity.PodName]
		name := row.Tags[container_entity.Name]
		metricType := row.Tags[container_entity.Metric]

		metricValue := datahub_v1alpha1.MetricType(datahub_v1alpha1.MetricType_value[metricType])
		switch metricType {
		case metric.TypeContainerCPUUsageSecondsPercentage:
			metricValue = datahub_v1alpha1.MetricType_CPU_USAGE_SECONDS_PERCENTAGE
		case metric.TypeContainerMemoryUsageBytes:
			metricValue = datahub_v1alpha1.MetricType_MEMORY_USAGE_BYTES
		}

		kind := metric.ContainerMetricKindRaw
		if val, ok := row.Tags[container_entity.Kind]; ok {
			if val != "" {
				kind = val
			}
		}

		granularity := int64(30)
		if val, ok := row.Tags[container_entity.Granularity]; ok {
			if val != "" {
				granularity, _ = strconv.ParseInt(val, 10, 64)
			}
		}

		podMap[namespace+"|"+podName] = &datahub_v1alpha1.PodPrediction{}
		podMap[namespace+"|"+podName].NamespacedName = &datahub_v1alpha1.NamespacedName{
			Namespace: namespace,
			Name:      podName,
		}

		podContainerMap[namespace+"|"+podName+"|"+name] = &datahub_v1alpha1.ContainerPrediction{}
		podContainerMap[namespace+"|"+podName+"|"+name].Name = name

		metricKey := namespace + "|" + podName + "|" + name + "|" + kind + "|" + metricType
		podContainerKindMetricMap[metricKey] = &datahub_v1alpha1.MetricData{}
		podContainerKindMetricMap[metricKey].MetricType = metricValue
		podContainerKindMetricMap[metricKey].Granularity = granularity

		for _, data := range row.Data {
			t, _ := time.Parse(time.RFC3339, data[container_entity.Time])
			value := data[container_entity.Value]

			googleTimestamp, _ := ptypes.TimestampProto(t)

			tempSample := &datahub_v1alpha1.Sample{
				Time:     googleTimestamp,
				NumValue: value,
			}
			podContainerKindMetricSampleMap[metricKey] = append(podContainerKindMetricSampleMap[metricKey], tempSample)
		}
	}

	for k := range podContainerKindMetricMap {
		namespace := strings.Split(k, "|")[0]
		podName := strings.Split(k, "|")[1]
		name := strings.Split(k, "|")[2]
		kind := strings.Split(k, "|")[3]
		metricType := strings.Split(k, "|")[4]

		containerKey := namespace + "|" + podName + "|" + name
		metricKey := namespace + "|" + podName + "|" + name + "|" + kind + "|" + metricType

		podContainerKindMetricMap[metricKey].Data = podContainerKindMetricSampleMap[metricKey]

		if kind == metric.ContainerMetricKindUpperbound {
			podContainerMap[containerKey].PredictedUpperboundData = append(podContainerMap[containerKey].PredictedUpperboundData, podContainerKindMetricMap[metricKey])
		} else if kind == metric.ContainerMetricKindLowerbound {
			podContainerMap[containerKey].PredictedLowerboundData = append(podContainerMap[containerKey].PredictedLowerboundData, podContainerKindMetricMap[metricKey])
		} else {
			podContainerMap[containerKey].PredictedRawData = append(podContainerMap[containerKey].PredictedRawData, podContainerKindMetricMap[metricKey])
		}
	}

	for k := range podContainerMap {
		namespace := strings.Split(k, "|")[0]
		podName := strings.Split(k, "|")[1]
		name := strings.Split(k, "|")[2]

		podKey := namespace + "|" + podName
		containerKey := namespace + "|" + podName + "|" + name

		podMap[podKey].ContainerPredictions = append(podMap[podKey].ContainerPredictions, podContainerMap[containerKey])
	}

	podList := make([]*datahub_v1alpha1.PodPrediction, 0)
	for k := range podMap {
		podList = append(podList, podMap[k])
	}

	return podList
}

func (r *ContainerRepository) buildInfluxQLWhereClauseFromRequest(request prediction_dao.ListPodPredictionsRequest) string {

	var (
		whereClause string
		conditions  string
	)

	if request.Namespace != "" {
		conditions += fmt.Sprintf(`"%s"='%s'`, container_entity.Namespace, request.Namespace)
	}
	if request.PodName != "" {
		if conditions != "" {
			conditions += fmt.Sprintf(` AND "%s"='%s'`, container_entity.PodName, request.PodName)
		} else {
			conditions += fmt.Sprintf(`"%s"='%s'`, container_entity.PodName, request.PodName)
		}
	}

	if conditions != "" {
		if request.Granularity == 30 {
			conditions += fmt.Sprintf(` AND ("%s"='' OR "%s"='%d')`, container_entity.Granularity, container_entity.Granularity, request.Granularity)
		} else {
			conditions += fmt.Sprintf(` AND "%s"='%d'`, container_entity.Granularity, request.Granularity)
		}
	} else {
		if request.Granularity == 30 {
			conditions += fmt.Sprintf(`("%s"='' OR "%s"='%d')`, container_entity.Granularity, container_entity.Granularity, request.Granularity)
		} else {
			conditions += fmt.Sprintf(`"%s"='%d'`, container_entity.Granularity, request.Granularity)
		}
	}

	if conditions != "" {
		whereClause = fmt.Sprintf("WHERE %s", conditions)
	}

	return whereClause
}
