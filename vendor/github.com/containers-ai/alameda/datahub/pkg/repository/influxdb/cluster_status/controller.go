package clusterstatus

import (
	"fmt"
	"strconv"
	"time"

	controller_entity "github.com/containers-ai/alameda/datahub/pkg/entity/influxdb/cluster_status"
	"github.com/containers-ai/alameda/datahub/pkg/repository/influxdb"
	datahub_api "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	influxdb_client "github.com/influxdata/influxdb/client/v2"
	"strings"
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

func (c *ControllerRepository) CreateControllers(controllers []*datahub_api.Controller) error {
	points := make([]*influxdb_client.Point, 0)
	for _, controller := range controllers {
		controllerNamespace := controller.GetControllerInfo().GetNamespacedName().GetNamespace()
		controllerName := controller.GetControllerInfo().GetNamespacedName().GetName()
		controllerKind := controller.GetControllerInfo().GetKind().String()
		controllerExecution := controller.GetEnableRecommendationExecution()
		controllerPolicy := controller.GetPolicy().String()

		ownerNamespace := ""
		ownerName := ""
		ownerKind := ""

		if len(controller.GetOwnerInfo()) > 0 {
			ownerNamespace = controller.GetOwnerInfo()[0].GetNamespacedName().GetNamespace()
			ownerName = controller.GetOwnerInfo()[0].GetNamespacedName().GetName()
			ownerKind = controller.GetOwnerInfo()[0].GetKind().String()
		}

		tags := map[string]string{
			string(controller_entity.ControllerNamespace):      controllerNamespace,
			string(controller_entity.ControllerName):           controllerName,
			string(controller_entity.ControllerOwnerNamespace): ownerNamespace,
			string(controller_entity.ControllerOwnerName):      ownerName,
		}

		fields := map[string]interface{}{
			string(controller_entity.ControllerKind):            controllerKind,
			string(controller_entity.ControllerOwnerKind):       ownerKind,
			string(controller_entity.ControllerReplicas):        controller.GetReplicas(),
			string(controller_entity.ControllerEnableExecution): strconv.FormatBool(controllerExecution),
			string(controller_entity.ControllerPolicy):          controllerPolicy,
		}

		pt, err := influxdb_client.NewPoint(string(Controller), tags, fields, time.Unix(0, 0))
		if err != nil {
			scope.Error(err.Error())
		}
		points = append(points, pt)
	}

	err := c.influxDB.WritePoints(points, influxdb_client.BatchPointsConfig{
		Database: string(influxdb.ClusterStatus),
	})

	if err != nil {
		scope.Error(err.Error())
	}

	return nil
}

func (c *ControllerRepository) ListControllers(in *datahub_api.ListControllersRequest) ([]*datahub_api.Controller, error) {
	namespace := in.GetNamespacedName().GetNamespace()
	name := in.GetNamespacedName().GetName()

	whereStr := c.convertQueryCondition(namespace, name)

	influxdbStatement := influxdb.Statement{
		Measurement: Controller,
		WhereClause: whereStr,
		GroupByTags: []string{controller_entity.ControllerNamespace, controller_entity.ControllerName},
	}

	cmd := influxdbStatement.BuildQueryCmd()

	results, err := c.influxDB.QueryDB(cmd, string(influxdb.ClusterStatus))
	if err != nil {
		return make([]*datahub_api.Controller, 0), err
	}

	influxdbRows := influxdb.PackMap(results)

	controllerList := c.getControllersFromInfluxRows(influxdbRows)
	return controllerList, nil
}

func (c *ControllerRepository) DeleteControllers(in *datahub_api.DeleteControllersRequest) error {
	controllers := in.GetControllers()
	whereStr := ""

	for _, controller := range controllers {
		namespace := controller.GetControllerInfo().GetNamespacedName().GetNamespace()
		name := controller.GetControllerInfo().GetNamespacedName().GetName()
		whereStr += fmt.Sprintf(" (\"name\"='%s' AND \"namespace\"='%s') AND", name, namespace)
	}

	whereStr = strings.TrimSuffix(whereStr, "AND")

	if whereStr != "" {
		whereStr = "WHERE" + whereStr
	}
	cmd := fmt.Sprintf("DROP SERIES FROM %s %s", string(Controller), whereStr)

	_, err := c.influxDB.QueryDB(cmd, string(influxdb.ClusterStatus))
	if err != nil {
		return err
	}

	return nil
}

func (c *ControllerRepository) convertQueryCondition(namespace string, name string) string {
	ret := ""

	if namespace != "" {
		ret += fmt.Sprintf("\"namespace\"='%s' ", namespace)
	}

	if name != "" {
		ret += fmt.Sprintf("AND \"name\"='%s' ", name)
	}

	ret = strings.TrimPrefix(ret, "AND")
	if ret != "" {
		ret = "WHERE " + ret
	}
	return ret
}

func (c *ControllerRepository) getControllersFromInfluxRows(rows []*influxdb.InfluxDBRow) []*datahub_api.Controller {
	controllerList := make([]*datahub_api.Controller, 0)
	for _, row := range rows {
		namespace := row.Tags[controller_entity.ControllerNamespace]
		name := row.Tags[controller_entity.ControllerName]

		tempController := &datahub_api.Controller{
			ControllerInfo: &datahub_api.ResourceInfo{
				NamespacedName: &datahub_api.NamespacedName{
					Namespace: namespace,
					Name:      name,
				},
			},
		}

		ownerInfoList := make([]*datahub_api.ResourceInfo, 0)
		for _, data := range row.Data {
			ownerNamespace := data[controller_entity.ControllerOwnerNamespace]
			ownerName := data[controller_entity.ControllerOwnerName]
			tempOwnerKind := data[controller_entity.ControllerOwnerKind]
			var ownerKind datahub_api.Kind

			if val, found := datahub_api.Kind_value[tempOwnerKind]; found {
				ownerKind = datahub_api.Kind(val)
			}

			tempOwner := &datahub_api.ResourceInfo{
				NamespacedName: &datahub_api.NamespacedName{
					Namespace: ownerNamespace,
					Name:      ownerName,
				},
				Kind: ownerKind,
			}

			ownerInfoList = append(ownerInfoList, tempOwner)

			//------
			tempKind := data[controller_entity.ControllerKind]
			var kind datahub_api.Kind
			if val, found := datahub_api.Kind_value[tempKind]; found {
				kind = datahub_api.Kind(val)
			}
			tempController.ControllerInfo.Kind = kind

			tempReplicas, _ := strconv.ParseInt(data[string(controller_entity.ControllerReplicas)], 10, 32)
			tempController.Replicas = int32(tempReplicas)

			enableExecution, _ := strconv.ParseBool(data[controller_entity.ControllerEnableExecution])
			tempController.EnableRecommendationExecution = enableExecution

			tempPolicy := data[controller_entity.ControllerPolicy]
			var policy datahub_api.RecommendationPolicy
			if val, found := datahub_api.RecommendationPolicy_value[tempPolicy]; found {
				policy = datahub_api.RecommendationPolicy(val)
			}
			tempController.Policy = policy
		}

		tempController.OwnerInfo = ownerInfoList
		controllerList = append(controllerList, tempController)
	}

	return controllerList
}
