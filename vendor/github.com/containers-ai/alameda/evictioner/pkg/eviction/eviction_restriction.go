package eviction

import (
	"fmt"
	"math"

	datahubutils "github.com/containers-ai/alameda/datahub/pkg/utils"
	autoscalingv1alpha1 "github.com/containers-ai/alameda/operator/pkg/apis/autoscaling/v1alpha1"
	utilsresource "github.com/containers-ai/alameda/operator/pkg/utils/resources"
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	core_v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type EvictionRestriction interface {
	IsEvictabel(pod *core_v1.Pod) (isEvictabel bool, err error)
}

type podReplicaStatus struct {
	preservedPodCount float64

	evictedPodCount float64
	runningPodCount float64
}

type evictionRestriction struct {
	preservationPercentage float64
	triggerThreshold       triggerThreshold

	alamedaScalerMap map[string]*autoscalingv1alpha1.AlamedaScaler

	podIDToPodRecommendationMap            map[string]*datahub_v1alpha1.PodRecommendation
	podIDToAlamedaResourceIDMap            map[string]string
	alamedaResourceIDToPodReplicaStatusMap map[string]*podReplicaStatus
}

func NewEvictionRestriction(client client.Client, preservationPercentage float64, triggerThreshold triggerThreshold, podRecommendations []*datahub_v1alpha1.PodRecommendation) EvictionRestriction {

	podIDToPodRecommendationMap := make(map[string]*datahub_v1alpha1.PodRecommendation)
	podIDToAlamedaResourceIDMap := make(map[string]string)
	alamedaResourceIDToPodReplicaStatusMap := make(map[string]*podReplicaStatus)
	for _, podRecommendation := range podRecommendations {

		copyPodRecommendation := proto.Clone(podRecommendation)
		podRecommendation = copyPodRecommendation.(*datahub_v1alpha1.PodRecommendation)

		if podRecommendation.NamespacedName == nil {
			scope.Warnf("skip PodRecommendation due to get nil NamespacedName")
			continue
		}

		podRecommendationNamespace := podRecommendation.NamespacedName.Namespace
		podRecommendationName := podRecommendation.NamespacedName.Name
		podNamespace := podRecommendationNamespace
		podName := podRecommendationName
		podID := fmt.Sprintf("%s/%s", podNamespace, podName)
		podIDToPodRecommendationMap[podID] = podRecommendation

		topController := podRecommendation.TopController
		if topController == nil || topController.NamespacedName == nil {
			scope.Warnf("skip PodRecommendation (%s/%s) due to get empty topController", podRecommendationNamespace, podRecommendationName)
			continue
		}
		alamedaResourceNamespace := topController.NamespacedName.Namespace
		alamedaResourceName := topController.NamespacedName.Name
		alamedaResourceKind := topController.Kind
		alamedaResourceID := fmt.Sprintf("%s.%s.%s", alamedaResourceNamespace, alamedaResourceName, alamedaResourceKind)
		podIDToAlamedaResourceIDMap[podID] = alamedaResourceID

		if _, exist := alamedaResourceIDToPodReplicaStatusMap[alamedaResourceID]; !exist {
			podReplicaStatus, err := buildPodReplicaStatus(client, alamedaResourceNamespace, alamedaResourceName, alamedaResourceKind, preservationPercentage)
			if err != nil {
				scope.Warnf("skip PodRecommendation (%s/%s) due to build PodReplicaStatus failed: %s", podRecommendationNamespace, podRecommendationName, err.Error())
				continue
			}
			alamedaResourceIDToPodReplicaStatusMap[alamedaResourceID] = &podReplicaStatus
		}
	}

	e := &evictionRestriction{
		preservationPercentage: preservationPercentage,
		triggerThreshold:       triggerThreshold,

		podIDToPodRecommendationMap:            podIDToPodRecommendationMap,
		podIDToAlamedaResourceIDMap:            podIDToAlamedaResourceIDMap,
		alamedaResourceIDToPodReplicaStatusMap: alamedaResourceIDToPodReplicaStatusMap,
	}
	return e
}

func (e *evictionRestriction) IsEvictabel(pod *core_v1.Pod) (bool, error) {

	podNamespace := pod.Namespace
	podName := pod.Name
	podID := fmt.Sprintf("%s/%s", podNamespace, podName)

	podRecommendation := e.podIDToPodRecommendationMap[podID]
	if !e.isPodEvictable(pod, podRecommendation) {
		return false, nil
	}

	ok, err := e.canRollingUpdatePod(podID)
	if err != nil {
		scope.Errorf("check if rolling update can perform on pod (%s) failed: %s", podID, err.Error())
		return false, err
	} else if !ok {
		return false, nil
	}

	return true, nil
}

func (e *evictionRestriction) canRollingUpdatePod(podID string) (bool, error) {

	alamedaResourceID, exist := e.podIDToAlamedaResourceIDMap[podID]
	if !exist {
		return false, errors.Errorf("topController owning pod does not exist")
	}
	podReplicaStatus, exist := e.alamedaResourceIDToPodReplicaStatusMap[alamedaResourceID]
	if !exist {
		return false, errors.Errorf("PodReplicaStatus of pod does not exit")
	}
	if podReplicaStatus.runningPodCount-podReplicaStatus.evictedPodCount > podReplicaStatus.preservedPodCount {
		podReplicaStatus.evictedPodCount++
		return true, nil
	} else {
		podRecommendationID := podID
		scope.Infof("Pod (%s) is not evictable, current running replicas count %.f is not greater then preseved replicas count %.f , ignore PodRecommendation (%s)",
			podID,
			podReplicaStatus.runningPodCount-podReplicaStatus.evictedPodCount,
			podReplicaStatus.preservedPodCount,
			podRecommendationID)
		return false, nil
	}

}

func (e *evictionRestriction) isPodEvictable(pod *core_v1.Pod, podRecomm *datahub_v1alpha1.PodRecommendation) bool {
	ctRecomms := podRecomm.GetContainerRecommendations()
	containers := pod.Spec.Containers
	for _, container := range containers {
		for _, recContainer := range ctRecomms {
			if container.Name != recContainer.GetName() {
				continue
			}
			if e.isContainerEvictable(pod, &container, recContainer) {
				return true
			}
		}
	}
	return false
}

func (e *evictionRestriction) isContainerEvictable(pod *core_v1.Pod, container *core_v1.Container, recContainer *datahub_v1alpha1.ContainerRecommendation) bool {
	cpuTriggerThreshold := e.triggerThreshold.CPU
	memoryTriggerThreshold := e.triggerThreshold.Memory

	if &container.Resources == nil || container.Resources.Limits == nil || container.Resources.Requests == nil {
		scope.Infof("Pod %s/%s selected to evict due to some resource of container %s not defined.",
			pod.GetNamespace(), pod.GetName(), recContainer.GetName())
		return true
	}

	for _, resourceType := range []core_v1.ResourceName{
		core_v1.ResourceMemory,
		core_v1.ResourceCPU,
	} {
		// resource limit check
		if _, ok := container.Resources.Limits[resourceType]; !ok {
			scope.Infof("Pod %s/%s selected to evict due to resource limit %s of container %s not defined.",
				pod.GetNamespace(), pod.GetName(), resourceType, recContainer.GetName())
			return true
		}

		for _, limitRec := range recContainer.GetLimitRecommendations() {
			if resourceType == core_v1.ResourceMemory && limitRec.GetMetricType() == datahub_v1alpha1.MetricType_MEMORY_USAGE_BYTES && len(limitRec.GetData()) > 0 {
				if limitRecVal, err := datahubutils.StringToFloat64(limitRec.GetData()[0].GetNumValue()); err == nil {
					limitRecVal = math.Ceil(limitRecVal)
					limitQuan := container.Resources.Limits[resourceType]
					delta := (math.Abs(float64(100*(limitRecVal-float64(limitQuan.Value())))) / float64(limitQuan.Value()))
					scope.Infof("Resource limit of %s pod %s/%s container %s checking eviction threshold (%v perentage). Current setting: %v, Recommended setting: %v",
						resourceType, pod.GetNamespace(), pod.GetName(), recContainer.GetName(), memoryTriggerThreshold, limitQuan.Value(), limitRecVal)
					if delta >= memoryTriggerThreshold {
						scope.Infof("Decide to evict pod %s/%s due to delta is %v >= %v (threshold)", pod.GetNamespace(), pod.GetName(), delta, memoryTriggerThreshold)
						return true
					}
				}
			}
			if resourceType == core_v1.ResourceCPU && limitRec.GetMetricType() == datahub_v1alpha1.MetricType_CPU_USAGE_SECONDS_PERCENTAGE && len(limitRec.GetData()) > 0 {
				if limitRecVal, err := datahubutils.StringToFloat64(limitRec.GetData()[0].GetNumValue()); err == nil {
					limitRecVal = math.Ceil(limitRecVal)
					limitQuan := container.Resources.Limits[resourceType]
					delta := (math.Abs(float64(100*(limitRecVal-float64(limitQuan.MilliValue())))) / float64(limitQuan.MilliValue()))
					scope.Infof("Resource limit of %s pod %s/%s container %s checking eviction threshold (%v perentage). Current setting: %v, Recommended setting: %v",
						resourceType, pod.GetNamespace(), pod.GetName(), recContainer.GetName(), cpuTriggerThreshold, limitQuan.MilliValue(), limitRecVal)
					if delta >= cpuTriggerThreshold {
						scope.Infof("Decide to evict pod %s/%s due to delta is %v >= %v (threshold)", pod.GetNamespace(), pod.GetName(), delta, cpuTriggerThreshold)
						return true
					}
				}
			}
		}

		// resource request check
		if _, ok := container.Resources.Requests[resourceType]; !ok {
			scope.Infof("Pod %s/%s selected to evict due to resource request %s of container %s not defined.",
				pod.GetNamespace(), pod.GetName(), resourceType, recContainer.GetName())
			return true
		}
		for _, reqRec := range recContainer.GetRequestRecommendations() {
			if resourceType == core_v1.ResourceMemory && reqRec.GetMetricType() == datahub_v1alpha1.MetricType_MEMORY_USAGE_BYTES && len(reqRec.GetData()) > 0 {
				if requestRecVal, err := datahubutils.StringToFloat64(reqRec.GetData()[0].GetNumValue()); err == nil {
					requestRecVal = math.Ceil(requestRecVal)
					requestQuan := container.Resources.Requests[resourceType]
					delta := (math.Abs(float64(100*(requestRecVal-float64(requestQuan.Value())))) / float64(requestQuan.Value()))
					scope.Infof("Resource request of %s pod %s/%s container %s checking eviction threshold (%v perentage). Current setting: %v, Recommended setting: %v",
						resourceType, pod.GetNamespace(), pod.GetName(), recContainer.GetName(), memoryTriggerThreshold, requestQuan.Value(), requestRecVal)
					if delta >= memoryTriggerThreshold {
						scope.Infof("Decide to evict pod %s/%s due to delta is %v >= %v (threshold)", pod.GetNamespace(), pod.GetName(), delta, memoryTriggerThreshold)
						return true
					}
				}
			}
			if resourceType == core_v1.ResourceCPU && reqRec.GetMetricType() == datahub_v1alpha1.MetricType_CPU_USAGE_SECONDS_PERCENTAGE && len(reqRec.GetData()) > 0 {
				if requestRecVal, err := datahubutils.StringToFloat64(reqRec.GetData()[0].GetNumValue()); err == nil {
					requestRecVal = math.Ceil(requestRecVal)
					requestQuan := container.Resources.Requests[resourceType]
					delta := (math.Abs(float64(100*(requestRecVal-float64(requestQuan.MilliValue())))) / float64(requestQuan.MilliValue()))
					scope.Infof("Resource request of %s pod %s/%s container %s checking eviction threshold (%v perentage). Current setting: %v, Recommended setting: %v",
						resourceType, pod.GetNamespace(), pod.GetName(), recContainer.GetName(), cpuTriggerThreshold, requestQuan.MilliValue(), requestRecVal)
					if delta >= cpuTriggerThreshold {
						scope.Infof("Decide to evict pod %s/%s due to delta is %v >= %v (threshold)", pod.GetNamespace(), pod.GetName(), delta, cpuTriggerThreshold)
						return true
					}
				}
			}
		}
	}
	return false
}

func buildPodReplicaStatus(k8sClient client.Client, namespace string, name string, kind datahub_v1alpha1.Kind, preservationPercentage float64) (podReplicaStatus, error) {

	podReplicaStatus := podReplicaStatus{}

	var pods []core_v1.Pod
	listResource := utilsresource.NewListResources(k8sClient)
	switch kind {
	case datahub_v1alpha1.Kind_DEPLOYMENT:
		currentPods, err := listResource.ListPodsByDeployment(namespace, name)
		if err != nil {
			return podReplicaStatus, errors.Errorf("%s", err.Error())
		}
		pods = currentPods
	case datahub_v1alpha1.Kind_DEPLOYMENTCONFIG:
		currentPods, err := listResource.ListPodsByDeploymentConfig(namespace, name)
		if err != nil {
			return podReplicaStatus, errors.Errorf("%s", err.Error())
		}
		pods = currentPods
	default:
		return podReplicaStatus, errors.Errorf("not supported controller type %s",
			datahub_v1alpha1.Kind_name[int32(kind)],
		)
	}

	for _, pod := range pods {
		if pod.Status.Phase == core_v1.PodRunning {
			podReplicaStatus.runningPodCount++
		}
	}
	podReplicaStatus.preservedPodCount = math.Ceil(float64(len(pods)) * (preservationPercentage / 100))

	return podReplicaStatus, nil
}
