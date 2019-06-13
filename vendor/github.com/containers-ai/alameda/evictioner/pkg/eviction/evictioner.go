package eviction

import (
	"context"
	"fmt"
	"math"
	"time"

	datahubutils "github.com/containers-ai/alameda/datahub/pkg/utils"
	autoscalingv1alpha1 "github.com/containers-ai/alameda/operator/pkg/apis/autoscaling/v1alpha1"
	utilsresource "github.com/containers-ai/alameda/operator/pkg/utils/resources"
	"github.com/containers-ai/alameda/pkg/consts"
	"github.com/containers-ai/alameda/pkg/utils"
	logUtil "github.com/containers-ai/alameda/pkg/utils/log"
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/genproto/googleapis/rpc/code"
	corev1 "k8s.io/api/core/v1"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	scope = logUtil.RegisterScope("evictioner", "alamedascaler evictioner", 0)
)

// Evictioner deletes pods which need to apply recommendation
type Evictioner struct {
	checkCycle  int64
	datahubClnt datahub_v1alpha1.DatahubServiceClient
	k8sClienit  client.Client
	evictCfg    Config
}

// NewEvictioner return Evictioner instance
func NewEvictioner(checkCycle int64,
	datahubClnt datahub_v1alpha1.DatahubServiceClient,
	k8sClienit client.Client,
	evictCfg Config) *Evictioner {
	return &Evictioner{
		checkCycle:  checkCycle,
		datahubClnt: datahubClnt,
		k8sClienit:  k8sClienit,
		evictCfg:    evictCfg,
	}
}

// Start checking pods need to apply recommendation
func (evictioner *Evictioner) Start() {
	go evictioner.evictProcess()
}

func (evictioner *Evictioner) evictProcess() {
	for {
		if !evictioner.evictCfg.Enable {
			return
		}
		appliablePodRecList, err := evictioner.listAppliablePodRecommendation()
		if err != nil {
			scope.Error(err.Error())
		}
		scope.Debugf("Applicable pod recommendation lists: %s", utils.InterfaceToString(appliablePodRecList))
		evictioner.evictPods(appliablePodRecList)
		time.Sleep(time.Duration(evictioner.checkCycle) * time.Second)
	}
}

func (evictioner *Evictioner) evictPods(recPodList []*datahub_v1alpha1.PodRecommendation) {
	for _, recPod := range recPodList {
		recPodIns := &corev1.Pod{}
		err := evictioner.k8sClienit.Get(context.TODO(), types.NamespacedName{
			Namespace: recPod.GetNamespacedName().GetNamespace(),
			Name:      recPod.GetNamespacedName().GetName(),
		}, recPodIns)
		if err != nil {
			if !k8s_errors.IsNotFound(err) {
				scope.Error(err.Error())
			}
			continue
		}
		err = evictioner.k8sClienit.Delete(context.TODO(), recPodIns)
		if err != nil {
			scope.Errorf("Evict pod (%s,%s) failed: %s", recPodIns.GetNamespace(), recPodIns.GetName(), err.Error())
		}
	}
}

func (evictioner *Evictioner) listAppliablePodRecommendation() ([]*datahub_v1alpha1.PodRecommendation, error) {

	appliablePodRecList := []*datahub_v1alpha1.PodRecommendation{}
	nowTimestamp := time.Now().Unix()

	resp, err := evictioner.listPodRecommsPossibleToApply(nowTimestamp)
	if err != nil {
		return appliablePodRecList, err
	} else if resp.Status == nil {
		return appliablePodRecList, fmt.Errorf("Receive nil status from datahub")
	} else if resp.Status.Code != int32(code.Code_OK) {
		return appliablePodRecList, fmt.Errorf("Status code not 0: receive status code: %d,message: %s", resp.GetStatus().GetCode(), resp.GetStatus().GetMessage())
	}

	podRecommsPossibleToApply := resp.GetPodRecommendations()
	scope.Debugf("Possible applicable pod recommendation lists: %s", utils.InterfaceToString(podRecommsPossibleToApply))

	enableScalerMap := map[string]bool{}
	for _, rec := range podRecommsPossibleToApply {
		startTime := rec.GetStartTime().GetSeconds()
		endTime := rec.GetEndTime().GetSeconds()
		if startTime >= nowTimestamp || nowTimestamp >= endTime {
			continue
		}

		if rec.GetNamespacedName() == nil {
			scope.Warn("receive pod recommendation with nil NamespacedName, skip this recommendation")
			continue
		}

		recNS := rec.GetNamespacedName().GetNamespace()
		recName := rec.GetNamespacedName().GetName()
		pod, err := evictioner.getPodInfo(recNS, recName)
		if err != nil {
			scope.Errorf("Get Pod (%s/%s) failed due to %s.", recNS, recName, err.Error())
			continue
		}

		alamRecomm, err := evictioner.getAlamRecommInfo(recNS, recName)
		if err != nil {
			scope.Errorf("Get AlamedaRecommendation (%s/%s) failed due to %s.", recNS, recName, err.Error())
			continue
		}

		if !evictioner.isPodEnableExecution(alamRecomm, enableScalerMap) {
			scope.Debugf("Pod (%s/%s) cannot be evicted because its execution is not enabled.", pod.GetNamespace(), pod.GetName())
			continue
		}

		if evictioner.isPodEvictable(pod, rec) {
			scope.Debugf("Pod (%s/%s) can be evicted.", pod.GetNamespace(), pod.GetName())
			appliablePodRecList = append(appliablePodRecList, rec)
		}
	}
	return appliablePodRecList, nil
}

func (evictioner *Evictioner) listPodRecommsPossibleToApply(nowTimestamp int64) (*datahub_v1alpha1.ListPodRecommendationsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	in := &datahub_v1alpha1.ListPodRecommendationsRequest{
		QueryCondition: &datahub_v1alpha1.QueryCondition{
			TimeRange: &datahub_v1alpha1.TimeRange{
				ApplyTime: &timestamp.Timestamp{
					Seconds: nowTimestamp,
				},
			},
			Order: datahub_v1alpha1.QueryCondition_DESC,
			Limit: 1,
		},
	}
	scope.Debugf("Request of ListAvailablePodRecommendations is %s.", utils.InterfaceToString(in))

	return evictioner.datahubClnt.ListAvailablePodRecommendations(ctx, in)
}

func (evictioner *Evictioner) getPodInfo(namespace, name string) (*corev1.Pod, error) {
	getResource := utilsresource.NewGetResource(evictioner.k8sClienit)
	pod, err := getResource.GetPod(namespace, name)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			scope.Debugf(err.Error())
		} else {
			scope.Errorf(err.Error())
		}
	}
	return pod, err
}

func (evictioner *Evictioner) getAlamRecommInfo(namespace, name string) (*autoscalingv1alpha1.AlamedaRecommendation, error) {
	getResource := utilsresource.NewGetResource(evictioner.k8sClienit)
	alamRecomm, err := getResource.GetAlamedaRecommendation(namespace, name)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			scope.Debugf(err.Error())
		} else {
			scope.Errorf(err.Error())
		}
	}
	return alamRecomm, err
}

func (evictioner *Evictioner) isContainerEvictable(pod *corev1.Pod, container *corev1.Container, recContainer *datahub_v1alpha1.ContainerRecommendation) bool {
	cpuTriggerThreshold := evictioner.evictCfg.TriggerThreshold.CPU
	memoryTriggerThreshold := evictioner.evictCfg.TriggerThreshold.Memory

	if &container.Resources == nil || container.Resources.Limits == nil || container.Resources.Requests == nil {
		scope.Infof("Pod %s/%s selected to evict due to some resource of container %s not defined.",
			pod.GetNamespace(), pod.GetName(), recContainer.GetName())
		return true
	}

	for _, resourceType := range []corev1.ResourceName{
		corev1.ResourceMemory,
		corev1.ResourceCPU,
	} {
		// resource limit check
		if _, ok := container.Resources.Limits[resourceType]; !ok {
			scope.Infof("Pod %s/%s selected to evict due to resource limit %s of container %s not defined.",
				pod.GetNamespace(), pod.GetName(), resourceType, recContainer.GetName())
			return true
		}

		for _, limitRec := range recContainer.GetLimitRecommendations() {
			if resourceType == corev1.ResourceMemory && limitRec.GetMetricType() == datahub_v1alpha1.MetricType_MEMORY_USAGE_BYTES && len(limitRec.GetData()) > 0 {
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
			if resourceType == corev1.ResourceCPU && limitRec.GetMetricType() == datahub_v1alpha1.MetricType_CPU_USAGE_SECONDS_PERCENTAGE && len(limitRec.GetData()) > 0 {
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
			if resourceType == corev1.ResourceMemory && reqRec.GetMetricType() == datahub_v1alpha1.MetricType_MEMORY_USAGE_BYTES && len(reqRec.GetData()) > 0 {
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
			if resourceType == corev1.ResourceCPU && reqRec.GetMetricType() == datahub_v1alpha1.MetricType_CPU_USAGE_SECONDS_PERCENTAGE && len(reqRec.GetData()) > 0 {
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

func (evictioner *Evictioner) isPodEvictable(pod *corev1.Pod, podRecomm *datahub_v1alpha1.PodRecommendation) bool {
	ctRecomms := podRecomm.GetContainerRecommendations()
	containers := pod.Spec.Containers
	for _, container := range containers {
		for _, recContainer := range ctRecomms {
			if container.Name != recContainer.GetName() {
				continue
			}
			if evictioner.isContainerEvictable(pod, &container, recContainer) {
				return true
			}
		}
	}
	return false
}

func (evictioner *Evictioner) isPodEnableExecution(alamRecomm *autoscalingv1alpha1.AlamedaRecommendation, enableScalerMap map[string]bool) bool {

	for _, or := range alamRecomm.OwnerReferences {
		if or.Kind != consts.K8S_KIND_ALAMEDASCALER {
			continue
		}

		if enabled, ok := enableScalerMap[fmt.Sprintf("%s/%s", alamRecomm.GetNamespace(), or.Name)]; enabled && ok {
			return true
		} else if !enabled && ok {
			return false
		}

		scaler, err := evictioner.getAlamedaScalerInfo(alamRecomm.GetNamespace(), or.Name)
		if err == nil {
			enableScalerMap[fmt.Sprintf("%s/%s", alamRecomm.GetNamespace(), or.Name)] = scaler.Spec.EnableExecution
			return scaler.Spec.EnableExecution
		}
		return false
	}
	return false
}

func (evictioner *Evictioner) getAlamedaScalerInfo(namespace, name string) (*autoscalingv1alpha1.AlamedaScaler, error) {
	getResource := utilsresource.NewGetResource(evictioner.k8sClienit)
	scaler, err := getResource.GetAlamedaScaler(namespace, name)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			scope.Debugf(err.Error())
		} else {
			scope.Errorf(err.Error())
		}
	}
	return scaler, err
}
