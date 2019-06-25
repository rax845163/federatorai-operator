package recommendation

import (
	recommendation_entity "github.com/containers-ai/alameda/datahub/pkg/entity/influxdb/recommendation"
	"github.com/containers-ai/alameda/datahub/pkg/repository/influxdb"
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	influxdb_client "github.com/influxdata/influxdb/client/v2"
	"strconv"
	"time"
)

type ControllerRepository struct {
	influxDB *influxdb.InfluxDBRepository
}

func NewControllerRepository(influxDBCfg *influxdb.Config) *ControllerRepository {
	return &ControllerRepository{
		influxDB: &influxdb.InfluxDBRepository{
			Address:  influxDBCfg.Address,
			Username: influxDBCfg.Username,
			Password: influxDBCfg.Password,
		},
	}
}

func (c *ControllerRepository) CreateControllerRecommendations(controllerRecommendations []*datahub_v1alpha1.ControllerRecommendation) error {
	points := make([]*influxdb_client.Point, 0)
	for _, conrollerRecommendation := range controllerRecommendations {
		recommendedType := conrollerRecommendation.GetRecommendedType()

		if recommendedType == datahub_v1alpha1.ControllerRecommendedType_CRT_Primitive {
			recommendedSpec := conrollerRecommendation.GetRecommendedSpec()

			tags := map[string]string{
				recommendation_entity.ControllerNamespace: recommendedSpec.GetNamespacedName().GetNamespace(),
				recommendation_entity.ControllerName:      recommendedSpec.GetNamespacedName().GetName(),
				recommendation_entity.ControllerType:      datahub_v1alpha1.ControllerRecommendedType_CRT_Primitive.String(),
			}

			fields := map[string]interface{}{
				recommendation_entity.ControllerCurrentReplicas: recommendedSpec.GetCurrentReplicas(),
				recommendation_entity.ControllerDesiredReplicas: recommendedSpec.GetDesiredReplicas(),
				recommendation_entity.ControllerCreateTime:      recommendedSpec.GetCreateTime().GetSeconds(),
				recommendation_entity.ControllerKind:            recommendedSpec.GetKind().String(),

				recommendation_entity.ControllerCurrentCPURequest: recommendedSpec.GetCurrentCpuRequests(),
				recommendation_entity.ControllerCurrentMEMRequest: recommendedSpec.GetCurrentMemRequests(),
				recommendation_entity.ControllerCurrentCPULimit:   recommendedSpec.GetCurrentCpuLimits(),
				recommendation_entity.ControllerCurrentMEMLimit:   recommendedSpec.GetCurrentMemLimits(),
				recommendation_entity.ControllerDesiredCPULimit:   recommendedSpec.GetDesiredCpuLimits(),
				recommendation_entity.ControllerDesiredMEMLimit:   recommendedSpec.GetDesiredMemLimits(),
			}

			pt, err := influxdb_client.NewPoint(string(Controller), tags, fields, time.Unix(recommendedSpec.GetTime().GetSeconds(), 0))
			if err != nil {
				scope.Error(err.Error())
			}

			points = append(points, pt)

		} else if recommendedType == datahub_v1alpha1.ControllerRecommendedType_CRT_K8s {
			recommendedSpec := conrollerRecommendation.GetRecommendedSpecK8S()

			tags := map[string]string{
				recommendation_entity.ControllerNamespace: recommendedSpec.GetNamespacedName().GetNamespace(),
				recommendation_entity.ControllerName:      recommendedSpec.GetNamespacedName().GetName(),
				recommendation_entity.ControllerType:      datahub_v1alpha1.ControllerRecommendedType_CRT_K8s.String(),
			}

			fields := map[string]interface{}{
				recommendation_entity.ControllerCurrentReplicas: recommendedSpec.GetCurrentReplicas(),
				recommendation_entity.ControllerDesiredReplicas: recommendedSpec.GetDesiredReplicas(),
				recommendation_entity.ControllerCreateTime:      recommendedSpec.GetCreateTime().GetSeconds(),
				recommendation_entity.ControllerKind:            recommendedSpec.GetKind().String(),
			}

			pt, err := influxdb_client.NewPoint(string(Controller), tags, fields, time.Unix(recommendedSpec.GetTime().GetSeconds(), 0))
			if err != nil {
				scope.Error(err.Error())
			}

			points = append(points, pt)
		}
	}

	err := c.influxDB.WritePoints(points, influxdb_client.BatchPointsConfig{
		Database: string(influxdb.Recommendation),
	})

	if err != nil {
		scope.Error(err.Error())
		return err
	}

	return nil
}

func (c *ControllerRepository) ListControllerRecommendations(in *datahub_v1alpha1.ListControllerRecommendationsRequest) ([]*datahub_v1alpha1.ControllerRecommendation, error) {
	namespace := in.GetNamespacedName().GetNamespace()
	name := in.GetNamespacedName().GetName()
	recommendationType := in.GetRecommendedType()

	influxdbStatement := influxdb.StatementNew{
		Measurement:    Controller,
		QueryCondition: in.GetQueryCondition(),
	}

	influxdbStatement.AppendWhereCondition(recommendation_entity.ControllerNamespace, "=", namespace)
	influxdbStatement.AppendWhereCondition(recommendation_entity.ControllerName, "=", name)
	influxdbStatement.AppendTimeConditionFromQueryCondition()
	influxdbStatement.AppendLimitClauseFromQueryCondition()
	influxdbStatement.AppendOrderClauseFromQueryCondition()

	if recommendationType != datahub_v1alpha1.ControllerRecommendedType_CRT_Undefined {
		influxdbStatement.AppendWhereCondition(recommendation_entity.ControllerType, "=", recommendationType.String())
	}

	cmd := influxdbStatement.BuildQueryCmd()

	results, err := c.influxDB.QueryDB(cmd, string(influxdb.Recommendation))
	if err != nil {
		return make([]*datahub_v1alpha1.ControllerRecommendation, 0), err
	}

	influxdbRows := influxdb.PackMap(results)
	recommendations := c.getControllersRecommendationsFromInfluxRows(influxdbRows)

	return recommendations, nil
}

func (c *ControllerRepository) getControllersRecommendationsFromInfluxRows(rows []*influxdb.InfluxDBRow) []*datahub_v1alpha1.ControllerRecommendation {
	recommendations := make([]*datahub_v1alpha1.ControllerRecommendation, 0)
	for _, influxdbRow := range rows {
		for _, data := range influxdbRow.Data {
			currentReplicas, _ := strconv.ParseInt(data[recommendation_entity.ControllerCurrentReplicas], 10, 64)
			desiredReplicas, _ := strconv.ParseInt(data[recommendation_entity.ControllerDesiredReplicas], 10, 64)
			createTime, _ := strconv.ParseInt(data[recommendation_entity.ControllerCreateTime], 10, 64)

			t, _ := time.Parse(time.RFC3339, data[recommendation_entity.ControllerTime])
			tempTime, _ := ptypes.TimestampProto(t)

			currentCpuRequests, _ := strconv.ParseFloat(data[recommendation_entity.ControllerCurrentCPURequest], 32)
			currentMemRequests, _ := strconv.ParseFloat(data[recommendation_entity.ControllerCurrentMEMRequest], 32)
			currentCpuLimits, _ := strconv.ParseFloat(data[recommendation_entity.ControllerCurrentCPULimit], 32)
			currentMemLimits, _ := strconv.ParseFloat(data[recommendation_entity.ControllerCurrentMEMLimit], 32)
			desiredCpuLimits, _ := strconv.ParseFloat(data[recommendation_entity.ControllerDesiredCPULimit], 32)
			desiredMemLimits, _ := strconv.ParseFloat(data[recommendation_entity.ControllerDesiredMEMLimit], 32)

			var commendationType datahub_v1alpha1.ControllerRecommendedType
			if tempType, exist := data[recommendation_entity.ControllerType]; exist {
				if value, ok := datahub_v1alpha1.ControllerRecommendedType_value[tempType]; ok {
					commendationType = datahub_v1alpha1.ControllerRecommendedType(value)
				}
			}

			var commendationKind datahub_v1alpha1.Kind
			if tempKind, exist := data[recommendation_entity.ControllerKind]; exist {
				if value, ok := datahub_v1alpha1.Kind_value[tempKind]; ok {
					commendationKind = datahub_v1alpha1.Kind(value)
				}
			}

			if commendationType == datahub_v1alpha1.ControllerRecommendedType_CRT_Primitive {
				tempRecommendation := &datahub_v1alpha1.ControllerRecommendation{
					RecommendedType: commendationType,
					RecommendedSpec: &datahub_v1alpha1.ControllerRecommendedSpec{
						NamespacedName: &datahub_v1alpha1.NamespacedName{
							Namespace: data[string(recommendation_entity.ControllerNamespace)],
							Name:      data[string(recommendation_entity.ControllerName)],
						},
						CurrentReplicas: int32(currentReplicas),
						DesiredReplicas: int32(desiredReplicas),
						Time:            tempTime,
						CreateTime: &timestamp.Timestamp{
							Seconds: createTime,
						},
						Kind:               commendationKind,
						CurrentCpuRequests: float32(currentCpuRequests),
						CurrentMemRequests: float32(currentMemRequests),
						CurrentCpuLimits:   float32(currentCpuLimits),
						CurrentMemLimits:   float32(currentMemLimits),
						DesiredCpuLimits:   float32(desiredCpuLimits),
						DesiredMemLimits:   float32(desiredMemLimits),
					},
				}

				recommendations = append(recommendations, tempRecommendation)

			} else if commendationType == datahub_v1alpha1.ControllerRecommendedType_CRT_K8s {
				tempRecommendation := &datahub_v1alpha1.ControllerRecommendation{
					RecommendedType: commendationType,
					RecommendedSpecK8S: &datahub_v1alpha1.ControllerRecommendedSpecK8S{
						NamespacedName: &datahub_v1alpha1.NamespacedName{
							Namespace: data[string(recommendation_entity.ControllerNamespace)],
							Name:      data[string(recommendation_entity.ControllerName)],
						},
						CurrentReplicas: int32(currentReplicas),
						DesiredReplicas: int32(desiredReplicas),
						Time:            tempTime,
						CreateTime: &timestamp.Timestamp{
							Seconds: createTime,
						},
						Kind: commendationKind,
					},
				}

				recommendations = append(recommendations, tempRecommendation)
			}
		}
	}

	return recommendations
}
