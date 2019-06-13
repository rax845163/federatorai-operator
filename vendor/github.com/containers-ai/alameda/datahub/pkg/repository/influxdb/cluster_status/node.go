package clusterstatus

import (
	"fmt"
	"strings"

	cluster_status_dao "github.com/containers-ai/alameda/datahub/pkg/dao/cluster_status"
	cluster_status_entity "github.com/containers-ai/alameda/datahub/pkg/entity/influxdb/cluster_status"
	"github.com/containers-ai/alameda/datahub/pkg/repository/influxdb"
	datahub_api "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	influxdb_client "github.com/influxdata/influxdb/client/v2"
	"github.com/pkg/errors"
)

type NodeRepository struct {
	influxDB *influxdb.InfluxDBRepository
}

func (nodeRepository *NodeRepository) IsTag(column string) bool {
	for _, tag := range cluster_status_entity.NodeTags {
		if column == string(tag) {
			return true
		}
	}
	return false
}

func NewNodeRepository(influxDBCfg *influxdb.Config) *NodeRepository {
	return &NodeRepository{
		influxDB: &influxdb.InfluxDBRepository{
			Address:  influxDBCfg.Address,
			Username: influxDBCfg.Username,
			Password: influxDBCfg.Password,
		},
	}
}

// AddAlamedaNodes add node information to database
func (nodeRepository *NodeRepository) AddAlamedaNodes(alamedaNodes []*datahub_v1alpha1.Node) error {
	points := []*influxdb_client.Point{}
	for _, alamedaNode := range alamedaNodes {
		isInCluster := true
		startTime := alamedaNode.StartTime.GetSeconds()
		entity := cluster_status_entity.NodeEntity{
			Time:        influxdb.ZeroTime,
			Name:        &alamedaNode.Name,
			IsInCluster: &isInCluster,
			CreatedTime: &startTime,
		}
		if nodeCapacity := alamedaNode.GetCapacity(); nodeCapacity != nil {
			entity.CPUCores = &nodeCapacity.CpuCores
			entity.MemoryBytes = &nodeCapacity.MemoryBytes
		}
		if pt, err := entity.InfluxDBPoint(string(Node)); err == nil {
			points = append(points, pt)
		} else {
			scope.Error(err.Error())
		}
	}
	err := nodeRepository.influxDB.WritePoints(points, influxdb_client.BatchPointsConfig{
		Database: string(influxdb.ClusterStatus),
	})
	if err != nil {
		return errors.Wrapf(err, "add alameda nodes failed: %s", err.Error())
	}
	return nil
}

func (nodeRepository *NodeRepository) RemoveAlamedaNodes(alamedaNodes []*datahub_v1alpha1.Node) error {
	hasErr := false
	errMsg := ""
	for _, alamedaNode := range alamedaNodes {
		cmd := fmt.Sprintf("DROP SERIES FROM %s WHERE \"%s\"='%s'",
			string(Node), string(cluster_status_entity.NodeName), alamedaNode.Name)
		_, err := nodeRepository.influxDB.QueryDB(cmd, string(influxdb.ClusterStatus))
		if err != nil {
			hasErr = true
			errMsg += errMsg + err.Error()
		}
	}
	if hasErr {
		return fmt.Errorf(errMsg)
	}
	return nil
}

func (nodeRepository *NodeRepository) ListAlamedaNodes(timeRange *datahub_api.TimeRange) ([]*cluster_status_entity.NodeEntity, error) {

	nodeEntities := []*cluster_status_entity.NodeEntity{}
	nodeCreatePeriodCondition := nodeRepository.getNodeCreatePeriodCondition(timeRange)

	cmd := fmt.Sprintf("SELECT * FROM %s WHERE \"%s\"=%s %s",
		string(Node), string(cluster_status_entity.NodeInCluster), "true", nodeCreatePeriodCondition)

	scope.Infof(fmt.Sprintf("Query nodes in cluster: %s", cmd))
	results, err := nodeRepository.influxDB.QueryDB(cmd, string(influxdb.ClusterStatus))
	if err != nil {
		return nodeEntities, errors.Wrap(err, "list alameda nodes from influxdb failed")
	}

	if len(results) == 1 && len(results[0].Series) > 0 {

		influxdbRows := influxdb.PackMap(results)
		for _, influxdbRow := range influxdbRows {
			for _, data := range influxdbRow.Data {
				nodeEntity := cluster_status_entity.NewNodeEntityFromMap(data)
				nodeEntities = append(nodeEntities, &nodeEntity)
			}
		}
	}

	return nodeEntities, nil
}

func (nodeRepository *NodeRepository) ListNodes(request cluster_status_dao.ListNodesRequest) ([]*cluster_status_entity.NodeEntity, error) {

	nodeEntities := []*cluster_status_entity.NodeEntity{}

	whereClause := nodeRepository.buildInfluxQLWhereClauseFromRequest(request)
	cmd := fmt.Sprintf("SELECT * FROM %s %s", string(Node), whereClause)
	scope.Debug(fmt.Sprintf("Query nodes in cluster: %s", cmd))
	results, err := nodeRepository.influxDB.QueryDB(cmd, string(influxdb.ClusterStatus))
	if err != nil {
		return nodeEntities, errors.Wrap(err, "list nodes from influxdb failed")
	}

	influxdbRows := influxdb.PackMap(results)
	for _, influxdbRow := range influxdbRows {
		for _, data := range influxdbRow.Data {
			nodeEntity := cluster_status_entity.NewNodeEntityFromMap(data)
			nodeEntities = append(nodeEntities, &nodeEntity)
		}
	}

	return nodeEntities, nil
}

func (nodeRepository *NodeRepository) buildInfluxQLWhereClauseFromRequest(request cluster_status_dao.ListNodesRequest) string {

	var (
		whereClause string
		conditions  string
	)

	conditions += fmt.Sprintf("\"%s\" = %t", cluster_status_entity.NodeInCluster, request.InCluster)

	statementFilteringNodes := ""
	for _, nodeName := range request.NodeNames {
		statementFilteringNodes += fmt.Sprintf(`"%s" = '%s' OR `, cluster_status_entity.NodeName, nodeName)
	}
	statementFilteringNodes = strings.TrimSuffix(statementFilteringNodes, "OR ")
	if statementFilteringNodes != "" {
		conditions = fmt.Sprintf("(%s) AND (%s)", conditions, statementFilteringNodes)
	}

	whereClause = fmt.Sprintf("WHERE %s", conditions)

	return whereClause
}

func (nodeRepository *NodeRepository) getNodeCreatePeriodCondition(timeRange *datahub_api.TimeRange) string {
	if timeRange == nil {
		return ""
	}

	var start int64 = 0
	var end int64 = 0

	if timeRange.StartTime != nil {
		start = timeRange.StartTime.Seconds
	}

	if timeRange.EndTime != nil {
		end = timeRange.EndTime.Seconds
	}

	if start == 0 && end == 0 {
		return ""
	} else if start == 0 && end != 0 {
		period := fmt.Sprintf(`AND "create_time" < %d`, end)
		return period
	} else if start != 0 && end == 0 {
		period := fmt.Sprintf(`AND "create_time" >= %d`, start)
		return period
	} else if start != 0 && end != 0 {
		period := fmt.Sprintf(`AND "create_time" >= %d AND "create_time" < %d`, start, end)
		return period
	}

	return ""
}
