package metric

import (
	"fmt"
	"github.com/containers-ai/alameda/datahub/pkg/dao"
	"github.com/containers-ai/alameda/datahub/pkg/kubernetes/metadata"
	"github.com/containers-ai/alameda/datahub/pkg/metric"
	"sort"
)

// MetricsDAO DAO interface of metric data.
type MetricsDAO interface {
	ListPodMetrics(ListPodMetricsRequest) (PodsMetricMap, error)
	ListNodesMetric(ListNodeMetricsRequest) (NodesMetricMap, error)
}

// ListPodMetricsRequest Argument of method ListPodMetrics
type ListPodMetricsRequest struct {
	Namespace metadata.NamespaceName
	PodName   metadata.PodName
	dao.QueryCondition
}

// ListNodeMetricsRequest Argument of method ListNodeMetrics
type ListNodeMetricsRequest struct {
	NodeNames []metadata.NodeName
	dao.QueryCondition
}

// GetNodeNames Return nodes name in request
func (r ListNodeMetricsRequest) GetNodeNames() []metadata.NodeName {
	return r.NodeNames
}

// GetEmptyNodeNames Return slice with one empty string element
func (r ListNodeMetricsRequest) GetEmptyNodeNames() []metadata.NodeName {
	return []metadata.NodeName{""}
}

// ContainerMetric Metric model to represent one container metric
type ContainerMetric struct {
	Namespace     metadata.NamespaceName
	PodName       metadata.PodName
	ContainerName metadata.ContainerName
	Metrics       map[metric.ContainerMetricType][]metric.Sample
}

// BuildPodMetric Build PodMetric consist of the receiver in ContainersMetricMap.
func (c *ContainerMetric) BuildPodMetric() *PodMetric {

	containersMetricMap := ContainersMetricMap{}
	containersMetricMap[c.NamespacePodContainerName()] = c

	return &PodMetric{
		Namespace:           c.Namespace,
		PodName:             c.PodName,
		ContainersMetricMap: &containersMetricMap,
	}
}

// NamespacePodContainerName Return identity of the container metric.
func (c ContainerMetric) NamespacePodContainerName() metadata.NamespacePodContainerName {
	return metadata.NamespacePodContainerName(fmt.Sprintf("%s/%s/%s", c.Namespace, c.PodName, c.ContainerName))
}

// SortByTimestamp Sort each metric samples by timestamp in input order
func (c *ContainerMetric) SortByTimestamp(order dao.Order) {

	for _, samples := range c.Metrics {
		if order == dao.Asc {
			sort.Sort(metric.SamplesByAscTimestamp(samples))
		} else {
			sort.Sort(metric.SamplesByDescTimestamp(samples))
		}
	}
}

// Limit Slicing each metric samples element
func (c *ContainerMetric) Limit(limit int) {

	if limit == 0 {
		return
	}

	for metricType, samples := range c.Metrics {
		c.Metrics[metricType] = samples[:limit]
	}
}

// ContainersMetricMap Containers metric map
type ContainersMetricMap map[metadata.NamespacePodContainerName]*ContainerMetric

// BuildPodsMetricMap Build PodsMetricMap base on current ContainersMetricMap
func (c ContainersMetricMap) BuildPodsMetricMap() *PodsMetricMap {

	var (
		podsMetricMap = &PodsMetricMap{}
	)

	for _, containerMetric := range c {
		podsMetricMap.AddContainerMetric(containerMetric)
	}

	return podsMetricMap
}

// Merge Merge current ContainersMetricMap with input ContainersMetricMap
func (c *ContainersMetricMap) Merge(in *ContainersMetricMap) {

	for namespacePodContainerName, containerMetric := range *in {
		if existedContainerMetric, exist := (*c)[namespacePodContainerName]; exist {
			for metricType, metrics := range containerMetric.Metrics {
				existedContainerMetric.Metrics[metricType] = append(existedContainerMetric.Metrics[metricType], metrics...)
			}
			(*c)[namespacePodContainerName] = existedContainerMetric
		} else {
			(*c)[namespacePodContainerName] = containerMetric
		}
	}
}

// PodMetric Metric model to represent one pod's metric
type PodMetric struct {
	Namespace           metadata.NamespaceName
	PodName             metadata.PodName
	ContainersMetricMap *ContainersMetricMap
}

// NamespacePodName Return identity of the pod metric
func (p PodMetric) NamespacePodName() metadata.NamespacePodName {
	return metadata.NamespacePodName(fmt.Sprintf("%s/%s", p.Namespace, p.PodName))
}

// Merge Merge current PodMetric with input PodMetric
func (p *PodMetric) Merge(in *PodMetric) {
	p.ContainersMetricMap.Merge(in.ContainersMetricMap)
}

// SortByTimestamp Sort each container metric's content
func (p *PodMetric) SortByTimestamp(order dao.Order) {

	for _, containerMetric := range *p.ContainersMetricMap {
		containerMetric.SortByTimestamp(order)
	}
}

// Limit Slicing each container metric content
func (p *PodMetric) Limit(limit int) {

	for _, containerMetric := range *p.ContainersMetricMap {
		containerMetric.Limit(limit)
	}
}

// PodsMetricMap Pods' metric map
type PodsMetricMap map[metadata.NamespacePodName]*PodMetric

// AddContainerMetric Add container metric into PodsMetricMap
func (p *PodsMetricMap) AddContainerMetric(c *ContainerMetric) {

	podMetric := c.BuildPodMetric()
	namespacePodName := podMetric.NamespacePodName()
	if existedPodMetric, exist := (*p)[namespacePodName]; exist {
		existedPodMetric.Merge(podMetric)
	} else {
		(*p)[namespacePodName] = podMetric
	}
}

// SortByTimestamp Sort each pod metric's content
func (p *PodsMetricMap) SortByTimestamp(order dao.Order) {

	for _, podMetric := range *p {
		podMetric.SortByTimestamp(order)
	}
}

// Limit Slicing each pod metric content
func (p *PodsMetricMap) Limit(limit int) {

	for _, podMetric := range *p {
		podMetric.Limit(limit)
	}
}

// NodeMetric Metric model to represent one node metric
type NodeMetric struct {
	NodeName metadata.NodeName
	Metrics  map[metric.NodeMetricType][]metric.Sample
}

// Merge Merge current NodeMetric with input NodeMetric
func (n *NodeMetric) Merge(in *NodeMetric) {

	for metricType, metrics := range in.Metrics {
		n.Metrics[metricType] = append(n.Metrics[metricType], metrics...)
	}
}

// SortByTimestamp Sort each metric samples by timestamp in input order
func (n *NodeMetric) SortByTimestamp(order dao.Order) {

	for _, samples := range n.Metrics {
		if order == dao.Asc {
			sort.Sort(metric.SamplesByAscTimestamp(samples))
		} else {
			sort.Sort(metric.SamplesByDescTimestamp(samples))
		}
	}
}

// Limit Slicing each metric samples element
func (n *NodeMetric) Limit(limit int) {

	if limit == 0 {
		return
	}

	for metricType, samples := range n.Metrics {
		n.Metrics[metricType] = samples[:limit]
	}
}

// NodesMetricMap Nodes' metric map
type NodesMetricMap map[metadata.NodeName]*NodeMetric

// AddNodeMetric Add node metric into NodesMetricMap
func (n *NodesMetricMap) AddNodeMetric(nodeMetric *NodeMetric) {

	nodeName := nodeMetric.NodeName
	if existNodeMetric, exist := (*n)[nodeName]; exist {
		existNodeMetric.Merge(nodeMetric)
	} else {
		(*n)[nodeName] = nodeMetric
	}
}

// SortByTimestamp Sort each node metric's content
func (n *NodesMetricMap) SortByTimestamp(order dao.Order) {

	for _, nodeMetric := range *n {
		nodeMetric.SortByTimestamp(order)
	}
}

// Limit Limit each node metric's content
func (n *NodesMetricMap) Limit(limit int) {

	for _, nodeMetric := range *n {
		nodeMetric.Limit(limit)
	}
}
