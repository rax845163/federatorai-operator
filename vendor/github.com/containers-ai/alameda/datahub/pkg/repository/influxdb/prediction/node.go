package prediction

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	prediction_dao "github.com/containers-ai/alameda/datahub/pkg/dao/prediction"
	node_entity "github.com/containers-ai/alameda/datahub/pkg/entity/influxdb/prediction/node"
	"github.com/containers-ai/alameda/datahub/pkg/repository/influxdb"
	datahub_utils "github.com/containers-ai/alameda/datahub/pkg/utils"
	influxdb_client "github.com/influxdata/influxdb/client/v2"
	"github.com/pkg/errors"

	"github.com/containers-ai/alameda/datahub/pkg/metric"
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	"github.com/golang/protobuf/ptypes"
)

type NodeRepository struct {
	influxDB *influxdb.InfluxDBRepository
}

func NewNodeRepositoryWithConfig(influxDBCfg influxdb.Config) *NodeRepository {
	return &NodeRepository{
		influxDB: &influxdb.InfluxDBRepository{
			Address:  influxDBCfg.Address,
			Username: influxDBCfg.Username,
			Password: influxDBCfg.Password,
		},
	}
}

func (r *NodeRepository) CreateNodePrediction(in *datahub_v1alpha1.CreateNodePredictionsRequest) error {

	points := make([]*influxdb_client.Point, 0)

	for _, nodePrediction := range in.GetNodePredictions() {
		nodeName := nodePrediction.GetName()
		isScheduled := nodePrediction.GetIsScheduled()

		r.appendMetricDataToPoints(metric.NodeMetricKindRaw, nodePrediction.GetPredictedRawData(), &points, nodeName, isScheduled)
		r.appendMetricDataToPoints(metric.NodeMetricKindUpperbound, nodePrediction.GetPredictedUpperboundData(), &points, nodeName, isScheduled)
		r.appendMetricDataToPoints(metric.NodeMetricKindLowerbound, nodePrediction.GetPredictedLowerboundData(), &points, nodeName, isScheduled)
	}

	err := r.influxDB.WritePoints(points, influxdb_client.BatchPointsConfig{
		Database: string(influxdb.Prediction),
	})
	if err != nil {
		return errors.Wrap(err, "create node prediction failed")
	}

	return nil
}

func (r *NodeRepository) appendMetricDataToPoints(kind metric.ContainerMetricKind, metricDataList []*datahub_v1alpha1.MetricData, points *[]*influxdb_client.Point, nodeName string, isScheduled bool) error {
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

		granularity := ""
		if metricData.GetGranularity() != 0 && metricData.GetGranularity() != 30 {
			granularity = strconv.FormatInt(metricData.GetGranularity(), 10)
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
				node_entity.Name:        nodeName,
				node_entity.IsScheduled: strconv.FormatBool(isScheduled),
				node_entity.Metric:      metricType,
				node_entity.Kind:        kind,
				node_entity.Granularity: granularity,
			}
			fields := map[string]interface{}{
				node_entity.Value: valueInFloat64,
			}
			point, err := influxdb_client.NewPoint(string(Node), tags, fields, time.Unix(tempTimeSeconds, 0))
			if err != nil {
				return errors.Wrap(err, "new influxdb data point failed")
			}
			*points = append(*points, point)
		}
	}

	return nil
}

func (r *NodeRepository) ListNodePredictionsByRequest(request prediction_dao.ListNodePredictionsRequest) ([]*datahub_v1alpha1.NodePrediction, error) {
	whereClause := r.buildInfluxQLWhereClauseFromRequest(request)
	influxdbStatement := influxdb.Statement{
		Measurement: Node,
		WhereClause: whereClause,
		//GroupByTags: []string{node_entity.Name, node_entity.Metric, node_entity.IsScheduled, node_entity.Kind, node_entity.Granularity},
		GroupByTags: []string{node_entity.Name, node_entity.Metric, node_entity.IsScheduled, node_entity.Kind},
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
		return []*datahub_v1alpha1.NodePrediction{}, errors.Wrap(err, "list node prediction failed")
	}

	rows := influxdb.PackMap(results)
	nodePredictions := r.getNodePredictionsFromInfluxRows(rows)

	return nodePredictions, nil
}

func (r *NodeRepository) getNodePredictionsFromInfluxRows(rows []*influxdb.InfluxDBRow) []*datahub_v1alpha1.NodePrediction {
	nodeMap := map[string]*datahub_v1alpha1.NodePrediction{}
	nodeMetricKindMap := map[string]*datahub_v1alpha1.MetricData{}
	nodeMetricKindSampleMap := map[string][]*datahub_v1alpha1.Sample{}

	for _, row := range rows {
		name := row.Tags[node_entity.Name]
		metricType := row.Tags[node_entity.Metric]
		isScheduled := row.Tags[node_entity.IsScheduled]

		metricValue := datahub_v1alpha1.MetricType(datahub_v1alpha1.MetricType_value[metricType])
		switch metricType {
		case metric.TypeContainerCPUUsageSecondsPercentage:
			metricValue = datahub_v1alpha1.MetricType_CPU_USAGE_SECONDS_PERCENTAGE
		case metric.TypeContainerMemoryUsageBytes:
			metricValue = datahub_v1alpha1.MetricType_MEMORY_USAGE_BYTES
		}

		kind := metric.NodeMetricKindRaw
		if val, ok := row.Tags[node_entity.Kind]; ok {
			if val != "" {
				kind = val
			}
		}

		granularity := int64(30)
		if val, ok := row.Tags[node_entity.Granularity]; ok {
			if val != "" {
				granularity, _ = strconv.ParseInt(val, 10, 64)
			}
		}

		nodeKey := name + "|" + isScheduled
		nodeMap[nodeKey] = &datahub_v1alpha1.NodePrediction{}
		nodeMap[nodeKey].Name = name
		nodeMap[nodeKey].IsScheduled, _ = strconv.ParseBool(isScheduled)

		metricKey := nodeKey + "|" + kind + "|" + metricType
		nodeMetricKindMap[metricKey] = &datahub_v1alpha1.MetricData{}
		nodeMetricKindMap[metricKey].MetricType = metricValue
		nodeMetricKindMap[metricKey].Granularity = granularity

		for _, data := range row.Data {
			t, _ := time.Parse(time.RFC3339, data[node_entity.Time])
			value := data[node_entity.Value]

			googleTimestamp, _ := ptypes.TimestampProto(t)

			tempSample := &datahub_v1alpha1.Sample{
				Time:     googleTimestamp,
				NumValue: value,
			}
			nodeMetricKindSampleMap[metricKey] = append(nodeMetricKindSampleMap[metricKey], tempSample)
		}
	}

	for k := range nodeMetricKindSampleMap {
		name := strings.Split(k, "|")[0]
		isScheduled := strings.Split(k, "|")[1]
		kind := strings.Split(k, "|")[2]
		metricType := strings.Split(k, "|")[3]

		nodeKey := name + "|" + isScheduled
		metricKey := nodeKey + "|" + kind + "|" + metricType

		nodeMetricKindMap[metricKey].Data = nodeMetricKindSampleMap[metricKey]

		if kind == metric.NodeMetricKindUpperbound {
			nodeMap[nodeKey].PredictedUpperboundData = append(nodeMap[nodeKey].PredictedUpperboundData, nodeMetricKindMap[metricKey])
		} else if kind == metric.NodeMetricKindLowerbound {
			nodeMap[nodeKey].PredictedLowerboundData = append(nodeMap[nodeKey].PredictedLowerboundData, nodeMetricKindMap[metricKey])
		} else {
			nodeMap[nodeKey].PredictedRawData = append(nodeMap[nodeKey].PredictedRawData, nodeMetricKindMap[metricKey])
		}
	}

	nodeList := make([]*datahub_v1alpha1.NodePrediction, 0)
	for k := range nodeMap {
		nodeList = append(nodeList, nodeMap[k])
	}

	return nodeList
}

func (r *NodeRepository) buildInfluxQLWhereClauseFromRequest(request prediction_dao.ListNodePredictionsRequest) string {

	var (
		whereClause string
		conditions  string
	)

	for _, nodeName := range request.NodeNames {
		conditions += fmt.Sprintf(`"%s" = '%s' or `, node_entity.Name, nodeName)
	}

	conditions = strings.TrimSuffix(conditions, "or ")

	if conditions != "" {
		if request.Granularity == 30 {
			conditions += fmt.Sprintf(` AND ("%s"='' OR "%s"='%d')`, node_entity.Granularity, node_entity.Granularity, request.Granularity)
		} else {
			conditions += fmt.Sprintf(` AND "%s"='%d'`, node_entity.Granularity, request.Granularity)
		}
	} else {
		if request.Granularity == 30 {
			conditions += fmt.Sprintf(`("%s"='' OR "%s"='%d')`, node_entity.Granularity, node_entity.Granularity, request.Granularity)
		} else {
			conditions += fmt.Sprintf(`"%s"='%d'`, node_entity.Granularity, request.Granularity)
		}
	}

	if conditions != "" {
		whereClause = fmt.Sprintf("where %s", conditions)
	}

	return whereClause
}
