package clusterstatus

import (
	"strconv"
	"time"

	"github.com/containers-ai/alameda/datahub/pkg/utils"
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	influxdb_client "github.com/influxdata/influxdb/client/v2"
)

type nodeField = string
type nodeTag = string

const (
	// NodeTime is the time node information is inserted to databse
	NodeTime nodeTag = "time"
	// NodeName is the name of node
	NodeName nodeTag = "name"

	// NodeGroup is node group name
	NodeGroup nodeField = "group"
	// NodeInCluster is the state node is in cluster or not
	NodeInCluster nodeField = "in_cluster"
	// NodeCPUCores is the amount of cores in node
	NodeCPUCores nodeField = "node_cpu_cores"
	// NodeMemoryBytes is the amount of momory bytes in node
	NodeMemoryBytes nodeField = "node_memory_bytes"

	NodeCreateTime nodeField = "create_time"
)

var (
	// NodeTags list tags of node measurement
	NodeTags = []nodeTag{NodeTime, NodeName}
	// NodeFields list fields of node measurement
	NodeFields = []nodeField{NodeGroup, NodeInCluster}
)

// NodeEntity is entity in database
type NodeEntity struct {
	Time        time.Time
	Name        *string
	NodeGroup   *string
	IsInCluster *bool
	CPUCores    *int64
	MemoryBytes *int64
	CreatedTime *int64
}

// NewNodeEntityFromMap Build entity from map
func NewNodeEntityFromMap(data map[string]string) NodeEntity {

	// TODO: log error
	tempTimestamp, _ := utils.ParseTime(data[ContainerTime])

	entity := NodeEntity{
		Time: tempTimestamp,
	}

	if name, exist := data[NodeName]; exist {
		entity.Name = &name
	}
	if nodeGroup, exist := data[NodeGroup]; exist {
		entity.NodeGroup = &nodeGroup
	}
	if isInCluster, exist := data[NodeInCluster]; exist {
		value, _ := strconv.ParseBool(isInCluster)
		entity.IsInCluster = &value
	}
	if cpuCores, exist := data[NodeCPUCores]; exist {
		value, _ := strconv.ParseInt(cpuCores, 10, 64)
		entity.CPUCores = &value
	}
	if memoryBytes, exist := data[NodeMemoryBytes]; exist {
		value, _ := strconv.ParseInt(memoryBytes, 10, 64)
		entity.MemoryBytes = &value
	}

	return entity
}

func (e NodeEntity) InfluxDBPoint(measurementName string) (*influxdb_client.Point, error) {

	tags := map[string]string{}
	if e.Name != nil {
		tags[NodeName] = *e.Name
	}

	fields := map[string]interface{}{}
	if e.NodeGroup != nil {
		fields[NodeGroup] = *e.NodeGroup
	}
	if e.IsInCluster != nil {
		fields[NodeInCluster] = *e.IsInCluster
	}
	if e.CPUCores != nil {
		fields[NodeCPUCores] = *e.CPUCores
	}
	if e.MemoryBytes != nil {
		fields[NodeMemoryBytes] = *e.MemoryBytes
	}
	if e.CreatedTime != nil {
		fields[NodeCreateTime] = *e.CreatedTime
	}

	return influxdb_client.NewPoint(measurementName, tags, fields, e.Time)
}

func (e NodeEntity) BuildDatahubNode() *datahub_v1alpha1.Node {

	node := &datahub_v1alpha1.Node{
		Capacity: &datahub_v1alpha1.Capacity{},
	}

	if e.Name != nil {
		node.Name = *e.Name
	}
	if e.CPUCores != nil {
		node.Capacity.CpuCores = *e.CPUCores
	}
	if e.MemoryBytes != nil {
		node.Capacity.MemoryBytes = *e.MemoryBytes
	}

	return node
}
