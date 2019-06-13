package datahub

import (
	"context"
	"math"
	"strconv"

	"github.com/containers-ai/alameda/admission-controller/pkg/recommendator/resource"
	"github.com/containers-ai/alameda/pkg/framework/datahub"
	"github.com/containers-ai/alameda/pkg/utils/log"
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pkg/errors"
	core_v1 "k8s.io/api/core/v1"
	k8s_resource "k8s.io/apimachinery/pkg/api/resource"
)

var (
	scope               = log.RegisterScope("resource-recommendator", "Datahub resource recommendator", 0)
	k8sKind_DatahubKind = map[string]datahub_v1alpha1.Kind{
		"Pod":              datahub_v1alpha1.Kind_POD,
		"Deployment":       datahub_v1alpha1.Kind_DEPLOYMENT,
		"DeploymentConfig": datahub_v1alpha1.Kind_DEPLOYMENTCONFIG,
	}
	datahubKind_K8SKind = map[datahub_v1alpha1.Kind]string{
		datahub_v1alpha1.Kind_POD:              "Pod",
		datahub_v1alpha1.Kind_DEPLOYMENT:       "Deployment",
		datahub_v1alpha1.Kind_DEPLOYMENTCONFIG: "DeploymentConfig",
	}
	datahubMetricType_K8SResourceName = map[datahub_v1alpha1.MetricType]core_v1.ResourceName{
		datahub_v1alpha1.MetricType_CPU_USAGE_SECONDS_PERCENTAGE: core_v1.ResourceCPU,
		datahub_v1alpha1.MetricType_MEMORY_USAGE_BYTES:           core_v1.ResourceMemory,
	}
)

var _ resource.ResourceRecommendator = &datahubResourceRecommendator{}

type datahubResourceRecommendator struct {
	datahubServiceClient datahub_v1alpha1.DatahubServiceClient
}

func NewDatahubResourceRecommendator(client datahub_v1alpha1.DatahubServiceClient) (resource.ResourceRecommendator, error) {

	return &datahubResourceRecommendator{
		datahubServiceClient: client,
	}, nil
}

func (dr *datahubResourceRecommendator) ListControllerPodResourceRecommendations(req resource.ListControllerPodResourceRecommendationsRequest) ([]*resource.PodResourceRecommendation, error) {

	recommendations := make([]*resource.PodResourceRecommendation, 0)

	datahubRequest, err := buildListAvailablePodRecommendationsRequest(req)
	if err != nil {
		return recommendations, errors.Wrap(err, "list controller pod resource recommendations failed")
	}
	scope.Debugf("query ListAvailablePodRecommendations to datahub, send request: %+v", datahubRequest)
	resp, err := dr.datahubServiceClient.ListAvailablePodRecommendations(context.Background(), datahubRequest)
	scope.Debugf("query ListAvailablePodRecommendations to datahub, received response: %+v", resp)
	if err != nil {
		return recommendations, errors.Wrap(err, "list controller pod resource recommendations failed")
	} else if _, err := datahub.IsResponseStatusOK(resp.Status); err != nil {
		return recommendations, errors.Wrap(err, "list controller pod resource recommendations failed")
	}

	for _, datahubPodRecommendation := range resp.GetPodRecommendations() {
		podRecommendation := buildPodResourceRecommendationFromDatahubPodRecommendation(datahubPodRecommendation)
		recommendations = append(recommendations, podRecommendation)
	}

	return recommendations, nil
}

func buildListAvailablePodRecommendationsRequest(request resource.ListControllerPodResourceRecommendationsRequest) (*datahub_v1alpha1.ListPodRecommendationsRequest, error) {

	var datahubRequest *datahub_v1alpha1.ListPodRecommendationsRequest

	datahubKind, exist := k8sKind_DatahubKind[request.Kind]
	if !exist {
		return datahubRequest, errors.Errorf("build Datahub ListPodRecommendationsRequest failed: no mapping Datahub kind for k8s kind: %s", request.Kind)
	}

	var queryTime *timestamp.Timestamp
	var err error
	if request.Time != nil {
		queryTime, err = ptypes.TimestampProto(*request.Time)
		if err != nil {
			return datahubRequest, errors.Errorf("build Datahub ListPodRecommendationsRequest failed: convert time.Time to google.Timestamp failed: %s", err.Error())
		}
	}

	datahubRequest = &datahub_v1alpha1.ListPodRecommendationsRequest{
		NamespacedName: &datahub_v1alpha1.NamespacedName{
			Namespace: request.Namespace,
			Name:      request.Name,
		},
		Kind: datahubKind,
		QueryCondition: &datahub_v1alpha1.QueryCondition{
			TimeRange: &datahub_v1alpha1.TimeRange{
				ApplyTime: queryTime,
			},
			Order: datahub_v1alpha1.QueryCondition_DESC,
			Limit: 1,
		},
	}
	return datahubRequest, nil
}

// TODO assign value of datahub.PodRecommendation.AssignedPodName to resource.Recommendation.AssignedPodName
func buildPodResourceRecommendationFromDatahubPodRecommendation(datahubPodRecommendation *datahub_v1alpha1.PodRecommendation) *resource.PodResourceRecommendation {

	namespace := ""
	name := ""
	if namespacedName := datahubPodRecommendation.GetNamespacedName(); namespacedName != nil {
		namespace = namespacedName.Namespace
		name = namespacedName.Name
	}

	startTime, _ := ptypes.Timestamp(datahubPodRecommendation.GetStartTime())
	endTime, _ := ptypes.Timestamp(datahubPodRecommendation.GetEndTime())

	topControllerKind := ""
	topControllerName := ""
	if datahubPodRecommendation.TopController != nil {
		topControllerKind = datahubKind_K8SKind[datahubPodRecommendation.TopController.Kind]
		if datahubPodRecommendation.TopController.NamespacedName != nil {
			topControllerName = datahubPodRecommendation.TopController.NamespacedName.Name
		}
	}

	podRecommendation := &resource.PodResourceRecommendation{
		Namespace:                        namespace,
		Name:                             name,
		TopControllerKind:                topControllerKind,
		TopControllerName:                topControllerName,
		ContainerResourceRecommendations: make([]*resource.ContainerResourceRecommendation, 0),
		ValidStartTime:                   startTime,
		ValidEndTime:                     endTime,
	}
	for _, datahubContainerRecommendation := range datahubPodRecommendation.GetContainerRecommendations() {
		containerResourceRecommendation := buildContainerResourceRecommendationFromDatahubContainerRecommendation(datahubContainerRecommendation)
		podRecommendation.ContainerResourceRecommendations = append(podRecommendation.ContainerResourceRecommendations, containerResourceRecommendation)
	}

	return podRecommendation
}

func buildContainerResourceRecommendationFromDatahubContainerRecommendation(datahubContainerRecommendation *datahub_v1alpha1.ContainerRecommendation) *resource.ContainerResourceRecommendation {

	containerResourceRecommendation := &resource.ContainerResourceRecommendation{
		Name: datahubContainerRecommendation.Name,
	}

	resourceLimitMap := datahubMetricDataSliceToMetricTypeValueMap(datahubContainerRecommendation.GetLimitRecommendations())
	containerResourceRecommendation.Limits = buildK8SReosurceListFromMetricTypeValueMap(resourceLimitMap)

	resourceRequestMap := datahubMetricDataSliceToMetricTypeValueMap(datahubContainerRecommendation.GetRequestRecommendations())
	containerResourceRecommendation.Requests = buildK8SReosurceListFromMetricTypeValueMap(resourceRequestMap)

	return containerResourceRecommendation
}

func datahubMetricDataSliceToMetricTypeValueMap(metricDataSlice []*datahub_v1alpha1.MetricData) map[datahub_v1alpha1.MetricType]string {

	resourceMap := make(map[datahub_v1alpha1.MetricType]string)

	for _, metricData := range metricDataSlice {
		sample := choseOneSample(metricData.GetData())
		if sample != nil {
			resourceMap[metricData.MetricType] = sample.NumValue
		}
	}

	return resourceMap
}

func choseOneSample(samples []*datahub_v1alpha1.Sample) *datahub_v1alpha1.Sample {

	if len(samples) > 0 {
		return samples[0]
	} else {
		return nil
	}
}

func buildK8SReosurceListFromMetricTypeValueMap(metricTypeValueMap map[datahub_v1alpha1.MetricType]string) core_v1.ResourceList {

	resourceList := make(core_v1.ResourceList)

	for metricType, value := range metricTypeValueMap {

		resourceUnit := ""
		if metricType == datahub_v1alpha1.MetricType_CPU_USAGE_SECONDS_PERCENTAGE {
			cpuMilliCores, err := strconv.ParseFloat(value, 64)
			if err != nil {

			}
			cpuMilliCores = math.Ceil(cpuMilliCores)
			value = strconv.FormatFloat(cpuMilliCores, 'f', 0, 64)
			resourceUnit = "m"
		}
		value = value + resourceUnit

		quantity, err := k8s_resource.ParseQuantity(value)
		if err != nil {
			scope.Warnf("parse value to k8s resource.Quantity failed, skip this recommendation: metricType:%s, value: %s, errMsg: %s", datahub_v1alpha1.MetricType_name[int32(metricType)], value, err.Error())
			continue
		}

		if k8sResourceName, exist := datahubMetricType_K8SResourceName[metricType]; !exist {
			scope.Warnf("no mapping k8s core_v1.ResourceName found for Datahub MetricType, skip this recommendation: metricType: %d", metricType)
			continue
		} else {
			resourceList[k8sResourceName] = quantity
		}
	}

	return resourceList
}
