package recommendation

import (
	"fmt"
	"strings"
	"time"

	recommendation_entity "github.com/containers-ai/alameda/datahub/pkg/entity/influxdb/recommendation"
	"github.com/containers-ai/alameda/datahub/pkg/entity/influxdb/utils/enumconv"
	"github.com/containers-ai/alameda/datahub/pkg/repository/influxdb"
	"github.com/containers-ai/alameda/datahub/pkg/utils"
	"github.com/containers-ai/alameda/pkg/utils/log"
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	"github.com/golang/protobuf/ptypes/timestamp"
	influxdb_client "github.com/influxdata/influxdb/client/v2"
	"github.com/pkg/errors"
	"strconv"
)

var (
	scope = log.RegisterScope("recommendation_db_measurement", "recommendation DB measurement", 0)
)

// ContainerRepository is used to operate node measurement of recommendation database
type ContainerRepository struct {
	influxDB *influxdb.InfluxDBRepository
}

// IsTag checks the column is tag or not
func (containerRepository *ContainerRepository) IsTag(column string) bool {
	for _, tag := range recommendation_entity.ContainerTags {
		if column == string(tag) {
			return true
		}
	}
	return false
}

// NewContainerRepository creates the ContainerRepository instance
func NewContainerRepository(influxDBCfg *influxdb.Config) *ContainerRepository {
	return &ContainerRepository{
		influxDB: &influxdb.InfluxDBRepository{
			Address:  influxDBCfg.Address,
			Username: influxDBCfg.Username,
			Password: influxDBCfg.Password,
		},
	}
}

// CreateContainerRecommendations add containers information container measurement
func (c *ContainerRepository) CreateContainerRecommendations(in *datahub_v1alpha1.CreatePodRecommendationsRequest) error {
	podRecommendations := in.GetPodRecommendations()
	granularity := in.GetGranularity()
	if granularity == 0 {
		granularity = 30
	}

	points := make([]*influxdb_client.Point, 0)
	for _, podRecommendation := range podRecommendations {
		if podRecommendation.GetApplyRecommendationNow() {
			//TODO
		}

		podNS := podRecommendation.GetNamespacedName().GetNamespace()
		podName := podRecommendation.GetNamespacedName().GetName()
		containerRecommendations := podRecommendation.GetContainerRecommendations()
		topController := podRecommendation.GetTopController()

		podPolicy := podRecommendation.GetAssignPodPolicy().GetPolicy()
		podPolicyValue := ""
		switch podPolicy.(type) {
		case *datahub_v1alpha1.AssignPodPolicy_NodeName:
			podPolicyValue = podPolicy.(*datahub_v1alpha1.AssignPodPolicy_NodeName).NodeName
		case *datahub_v1alpha1.AssignPodPolicy_NodePriority:
			nodeList := podPolicy.(*datahub_v1alpha1.AssignPodPolicy_NodePriority).NodePriority.GetNodes()
			if len(nodeList) > 0 {
				podPolicyValue = nodeList[0]
			}
			podPolicyValue = podPolicy.(*datahub_v1alpha1.AssignPodPolicy_NodePriority).NodePriority.GetNodes()[0]
		case *datahub_v1alpha1.AssignPodPolicy_NodeSelector:
			nodeMap := podPolicy.(*datahub_v1alpha1.AssignPodPolicy_NodeSelector).NodeSelector.Selector
			for _, value := range nodeMap {
				podPolicyValue = value
				break
			}
		}

		for _, containerRecommendation := range containerRecommendations {
			tags := map[string]string{
				recommendation_entity.ContainerNamespace:   podNS,
				recommendation_entity.ContainerPodName:     podName,
				recommendation_entity.ContainerName:        containerRecommendation.GetName(),
				recommendation_entity.ContainerGranularity: strconv.FormatInt(granularity, 10),
			}
			fields := map[string]interface{}{
				//TODO
				//string(recommendation_entity.ContainerPolicy):            "",
				recommendation_entity.ContainerTopControllerName: topController.GetNamespacedName().GetName(),
				recommendation_entity.ContainerTopControllerKind: enumconv.KindDisp[(topController.GetKind())],
				recommendation_entity.ContainerPolicy:            podPolicyValue,
				recommendation_entity.ContainerPolicyTime:        podRecommendation.GetAssignPodPolicy().GetTime().GetSeconds(),
			}

			for _, metricData := range containerRecommendation.GetInitialLimitRecommendations() {
				for _, datum := range metricData.GetData() {
					newFields := map[string]interface{}{}
					for key, value := range fields {
						newFields[key] = value
					}
					newFields[recommendation_entity.ContainerStartTime] = datum.GetTime().GetSeconds()
					newFields[recommendation_entity.ContainerEndTime] = datum.GetEndTime().GetSeconds()

					switch metricData.GetMetricType() {
					case datahub_v1alpha1.MetricType_CPU_USAGE_SECONDS_PERCENTAGE:
						if numVal, err := utils.StringToFloat64(strings.Split(datum.GetNumValue(), ".")[0]); err == nil {
							newFields[recommendation_entity.ContainerInitialResourceLimitCPU] = numVal
						}
					case datahub_v1alpha1.MetricType_MEMORY_USAGE_BYTES:
						if numVal, err := utils.StringToFloat64(strings.Split(datum.GetNumValue(), ".")[0]); err == nil {
							newFields[recommendation_entity.ContainerInitialResourceLimitMemory] = numVal
						}
					}

					if pt, err := influxdb_client.NewPoint(string(Container), tags, newFields, time.Unix(datum.GetTime().GetSeconds(), 0)); err == nil {
						points = append(points, pt)
					} else {
						scope.Error(err.Error())
					}
				}
			}

			for _, metricData := range containerRecommendation.GetInitialRequestRecommendations() {
				for _, datum := range metricData.GetData() {
					newFields := map[string]interface{}{}
					for key, value := range fields {
						newFields[key] = value
					}
					newFields[recommendation_entity.ContainerStartTime] = datum.GetTime().GetSeconds()
					newFields[recommendation_entity.ContainerEndTime] = datum.GetEndTime().GetSeconds()

					switch metricData.GetMetricType() {
					case datahub_v1alpha1.MetricType_CPU_USAGE_SECONDS_PERCENTAGE:
						if numVal, err := utils.StringToFloat64(strings.Split(datum.GetNumValue(), ".")[0]); err == nil {
							newFields[recommendation_entity.ContainerInitialResourceRequestCPU] = numVal
						}
					case datahub_v1alpha1.MetricType_MEMORY_USAGE_BYTES:
						if numVal, err := utils.StringToFloat64(strings.Split(datum.GetNumValue(), ".")[0]); err == nil {
							newFields[recommendation_entity.ContainerInitialResourceRequestMemory] = numVal
						}
					}

					if pt, err := influxdb_client.NewPoint(string(Container), tags, newFields, time.Unix(datum.GetTime().GetSeconds(), 0)); err == nil {
						points = append(points, pt)
					} else {
						scope.Error(err.Error())
					}
				}
			}

			for _, metricData := range containerRecommendation.GetLimitRecommendations() {
				for _, datum := range metricData.GetData() {
					newFields := map[string]interface{}{}
					for key, value := range fields {
						newFields[key] = value
					}
					newFields[recommendation_entity.ContainerStartTime] = datum.GetTime().GetSeconds()
					newFields[recommendation_entity.ContainerEndTime] = datum.GetEndTime().GetSeconds()

					switch metricData.GetMetricType() {
					case datahub_v1alpha1.MetricType_CPU_USAGE_SECONDS_PERCENTAGE:
						if numVal, err := utils.StringToFloat64(strings.Split(datum.GetNumValue(), ".")[0]); err == nil {
							newFields[recommendation_entity.ContainerResourceLimitCPU] = numVal
						}
					case datahub_v1alpha1.MetricType_MEMORY_USAGE_BYTES:
						if numVal, err := utils.StringToFloat64(strings.Split(datum.GetNumValue(), ".")[0]); err == nil {
							newFields[recommendation_entity.ContainerResourceLimitMemory] = numVal
						}
					}

					if pt, err := influxdb_client.NewPoint(string(Container), tags, newFields, time.Unix(datum.GetTime().GetSeconds(), 0)); err == nil {
						points = append(points, pt)
					} else {
						scope.Error(err.Error())
					}
				}
			}

			for _, metricData := range containerRecommendation.GetRequestRecommendations() {
				for _, datum := range metricData.GetData() {
					newFields := map[string]interface{}{}
					for key, value := range fields {
						newFields[key] = value
					}
					newFields[recommendation_entity.ContainerStartTime] = datum.GetTime().GetSeconds()
					newFields[recommendation_entity.ContainerEndTime] = datum.GetEndTime().GetSeconds()

					switch metricData.GetMetricType() {
					case datahub_v1alpha1.MetricType_CPU_USAGE_SECONDS_PERCENTAGE:
						if numVal, err := utils.StringToFloat64(strings.Split(datum.GetNumValue(), ".")[0]); err == nil {
							newFields[recommendation_entity.ContainerResourceRequestCPU] = numVal
						}
					case datahub_v1alpha1.MetricType_MEMORY_USAGE_BYTES:
						if numVal, err := utils.StringToFloat64(strings.Split(datum.GetNumValue(), ".")[0]); err == nil {
							newFields[recommendation_entity.ContainerResourceRequestMemory] = numVal
						}
					}
					if pt, err := influxdb_client.NewPoint(string(Container),
						tags, newFields,
						time.Unix(datum.GetTime().GetSeconds(), 0)); err == nil {
						points = append(points, pt)
					} else {
						scope.Error(err.Error())
					}
				}
			}
		}
	}
	err := c.influxDB.WritePoints(points, influxdb_client.BatchPointsConfig{
		Database: string(influxdb.Recommendation),
	})

	if err != nil {
		return err
	}
	return nil
}

// ListContainerRecommendations list container recommendations
func (c *ContainerRepository) ListContainerRecommendations(in *datahub_v1alpha1.ListPodRecommendationsRequest) ([]*datahub_v1alpha1.PodRecommendation, error) {
	kind := in.GetKind()
	granularity := in.GetGranularity()

	podRecommendations := make([]*datahub_v1alpha1.PodRecommendation, 0)

	influxdbStatement := influxdb.StatementNew{
		Measurement:    Container,
		QueryCondition: in.GetQueryCondition(),
		GroupByTags:    []string{recommendation_entity.ContainerName, recommendation_entity.ContainerNamespace, recommendation_entity.ContainerPodName},
	}

	nameCol := ""
	switch kind {
	case datahub_v1alpha1.Kind_POD:
		nameCol = string(recommendation_entity.ContainerPodName)
	case datahub_v1alpha1.Kind_DEPLOYMENT:
		nameCol = string(recommendation_entity.ContainerTopControllerName)
	case datahub_v1alpha1.Kind_DEPLOYMENTCONFIG:
		nameCol = string(recommendation_entity.ContainerTopControllerName)
	default:
		return podRecommendations, errors.Errorf("no matching kind for Datahub Kind, received Kind: %s", datahub_v1alpha1.Kind_name[int32(kind)])
	}
	influxdbStatement.AppendWhereCondition(recommendation_entity.ContainerNamespace, "=", in.GetNamespacedName().GetNamespace())
	influxdbStatement.AppendWhereCondition(nameCol, "=", in.GetNamespacedName().GetName())

	influxdbStatement.AppendTimeConditionFromQueryCondition()

	if kind != datahub_v1alpha1.Kind_POD {
		kindConditionStr := fmt.Sprintf("\"%s\"='%s'", recommendation_entity.ContainerTopControllerKind, enumconv.KindDisp[kind])
		influxdbStatement.AppendWhereCondition(recommendation_entity.ContainerTopControllerKind, "=", kindConditionStr)
	}

	if granularity == 0 || granularity == 30 {
		tempCondition := fmt.Sprintf("(\"%s\"='' OR \"%s\"='30')", recommendation_entity.ContainerGranularity, recommendation_entity.ContainerGranularity)
		influxdbStatement.AppendWhereConditionDirect(tempCondition)
	} else {
		influxdbStatement.AppendWhereCondition(recommendation_entity.ContainerGranularity, "=", strconv.FormatInt(granularity, 10))
	}

	influxdbStatement.AppendOrderClauseFromQueryCondition()
	influxdbStatement.AppendLimitClauseFromQueryCondition()

	cmd := influxdbStatement.BuildQueryCmd()
	scope.Debugf(fmt.Sprintf("ListContainerRecommendations: %s", cmd))

	podRecommendations, err := c.queryRecommendationNew(cmd, granularity)
	if err != nil {
		return podRecommendations, err
	}

	return podRecommendations, nil
}

func (c *ContainerRepository) ListAvailablePodRecommendations(in *datahub_v1alpha1.ListPodRecommendationsRequest) ([]*datahub_v1alpha1.PodRecommendation, error) {
	kind := in.GetKind()
	granularity := in.GetGranularity()

	podRecommendations := make([]*datahub_v1alpha1.PodRecommendation, 0)

	influxdbStatement := influxdb.StatementNew{
		Measurement:    Container,
		QueryCondition: in.GetQueryCondition(),
		GroupByTags:    []string{recommendation_entity.ContainerName, recommendation_entity.ContainerNamespace, recommendation_entity.ContainerPodName},
	}

	nameCol := ""
	switch kind {
	case datahub_v1alpha1.Kind_POD:
		nameCol = string(recommendation_entity.ContainerPodName)
	case datahub_v1alpha1.Kind_DEPLOYMENT:
		nameCol = string(recommendation_entity.ContainerTopControllerName)
	case datahub_v1alpha1.Kind_DEPLOYMENTCONFIG:
		nameCol = string(recommendation_entity.ContainerTopControllerName)
	default:
		return podRecommendations, errors.Errorf("no matching kind for Datahub Kind, received Kind: %s", datahub_v1alpha1.Kind_name[int32(kind)])
	}
	influxdbStatement.AppendWhereCondition(recommendation_entity.ContainerNamespace, "=", in.GetNamespacedName().GetNamespace())
	influxdbStatement.AppendWhereCondition(nameCol, "=", in.GetNamespacedName().GetName())

	if granularity == 0 || granularity == 30 {
		tempCondition := fmt.Sprintf("(\"%s\"='' OR \"%s\"='30')", recommendation_entity.ContainerGranularity, recommendation_entity.ContainerGranularity)
		influxdbStatement.AppendWhereConditionDirect(tempCondition)
	} else {
		influxdbStatement.AppendWhereCondition(recommendation_entity.ContainerGranularity, "=", strconv.FormatInt(granularity, 10))
	}

	whereStrTime := ""
	applyTime := in.GetQueryCondition().GetTimeRange().GetApplyTime().GetSeconds()
	if applyTime > 0 {
		whereStrTime = fmt.Sprintf(" \"end_time\">=%d AND \"start_time\"<=%d", applyTime, applyTime)
	}
	influxdbStatement.AppendWhereConditionDirect(whereStrTime)

	influxdbStatement.AppendOrderClauseFromQueryCondition()
	influxdbStatement.AppendLimitClauseFromQueryCondition()

	cmd := influxdbStatement.BuildQueryCmd()
	scope.Debugf(fmt.Sprintf("ListContainerRecommendations: %s", cmd))

	podRecommendations, err := c.queryRecommendationNew(cmd, granularity)
	if err != nil {
		return podRecommendations, err
	}

	return podRecommendations, nil
}

func (c *ContainerRepository) queryRecommendationNew(cmd string, granularity int64) ([]*datahub_v1alpha1.PodRecommendation, error) {
	podRecommendations := make([]*datahub_v1alpha1.PodRecommendation, 0)

	results, err := c.influxDB.QueryDB(cmd, string(influxdb.Recommendation))
	if err != nil {
		return podRecommendations, err
	}

	rows := influxdb.PackMap(results)

	for _, row := range rows {
		for _, data := range row.Data {
			podRecommendation := &datahub_v1alpha1.PodRecommendation{}
			podRecommendation.NamespacedName = &datahub_v1alpha1.NamespacedName{
				Namespace: data[recommendation_entity.ContainerNamespace],
				Name:      data[recommendation_entity.ContainerPodName],
			}

			tempTopControllerKind := data[recommendation_entity.ContainerTopControllerKind]
			var topControllerKind datahub_v1alpha1.Kind
			if val, ok := enumconv.KindEnum[tempTopControllerKind]; ok {
				topControllerKind = val
			}

			podRecommendation.TopController = &datahub_v1alpha1.TopController{
				NamespacedName: &datahub_v1alpha1.NamespacedName{
					Namespace: data[recommendation_entity.ContainerNamespace],
					Name:      data[recommendation_entity.ContainerTopControllerName],
				},
				Kind: topControllerKind,
			}

			startTime, _ := strconv.ParseInt(data[recommendation_entity.ContainerStartTime], 10, 64)
			endTime, _ := strconv.ParseInt(data[recommendation_entity.ContainerEndTime], 10, 64)

			podRecommendation.StartTime = &timestamp.Timestamp{
				Seconds: startTime,
			}

			podRecommendation.EndTime = &timestamp.Timestamp{
				Seconds: endTime,
			}

			policyTime, _ := strconv.ParseInt(data[recommendation_entity.ContainerPolicyTime], 10, 64)
			podRecommendation.AssignPodPolicy = &datahub_v1alpha1.AssignPodPolicy{
				Time: &timestamp.Timestamp{
					Seconds: policyTime,
				},
				Policy: &datahub_v1alpha1.AssignPodPolicy_NodeName{
					NodeName: data[recommendation_entity.ContainerPolicy],
				},
			}

			containerRecommendation := &datahub_v1alpha1.ContainerRecommendation{}
			containerRecommendation.Name = data[recommendation_entity.ContainerName]

			metricTypeList := []datahub_v1alpha1.MetricType{datahub_v1alpha1.MetricType_CPU_USAGE_SECONDS_PERCENTAGE, datahub_v1alpha1.MetricType_MEMORY_USAGE_BYTES}
			sampleTime := &timestamp.Timestamp{
				Seconds: startTime,
			}
			sampleEndTime := &timestamp.Timestamp{
				Seconds: endTime,
			}

			//
			for _, metricType := range metricTypeList {
				metricDataList := make([]*datahub_v1alpha1.MetricData, 0)
				for a := 0; a < 4; a++ {
					sample := &datahub_v1alpha1.Sample{
						Time:    sampleTime,
						EndTime: sampleEndTime,
					}

					metricData := &datahub_v1alpha1.MetricData{
						MetricType:  metricType,
						Granularity: granularity,
					}
					metricData.Data = append(metricData.Data, sample)
					metricDataList = append(metricDataList, metricData)
				}

				containerRecommendation.LimitRecommendations = append(containerRecommendation.LimitRecommendations, metricDataList[0])
				containerRecommendation.RequestRecommendations = append(containerRecommendation.RequestRecommendations, metricDataList[1])
				containerRecommendation.InitialLimitRecommendations = append(containerRecommendation.InitialLimitRecommendations, metricDataList[2])
				containerRecommendation.InitialRequestRecommendations = append(containerRecommendation.InitialRequestRecommendations, metricDataList[3])
			}

			containerRecommendation.LimitRecommendations[0].Data[0].NumValue = data[recommendation_entity.ContainerResourceLimitCPU]
			containerRecommendation.LimitRecommendations[1].Data[0].NumValue = data[recommendation_entity.ContainerResourceLimitMemory]

			containerRecommendation.RequestRecommendations[0].Data[0].NumValue = data[recommendation_entity.ContainerResourceRequestCPU]
			containerRecommendation.RequestRecommendations[1].Data[0].NumValue = data[recommendation_entity.ContainerResourceRequestMemory]

			containerRecommendation.InitialLimitRecommendations[0].Data[0].NumValue = data[recommendation_entity.ContainerInitialResourceLimitCPU]
			containerRecommendation.InitialLimitRecommendations[1].Data[0].NumValue = data[recommendation_entity.ContainerInitialResourceLimitMemory]

			containerRecommendation.InitialRequestRecommendations[0].Data[0].NumValue = data[recommendation_entity.ContainerInitialResourceRequestCPU]
			containerRecommendation.InitialRequestRecommendations[1].Data[0].NumValue = data[recommendation_entity.ContainerInitialResourceRequestMemory]

			podRecommendation.ContainerRecommendations = append(podRecommendation.ContainerRecommendations, containerRecommendation)

			podRecommendations = append(podRecommendations, podRecommendation)
		}
	}

	return podRecommendations, nil
}
