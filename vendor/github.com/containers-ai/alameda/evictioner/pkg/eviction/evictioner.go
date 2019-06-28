package eviction

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	autoscalingv1alpha1 "github.com/containers-ai/alameda/operator/pkg/apis/autoscaling/v1alpha1"
	utilsresource "github.com/containers-ai/alameda/operator/pkg/utils/resources"
	"github.com/containers-ai/alameda/pkg/consts"
	"github.com/containers-ai/alameda/pkg/utils"
	logUtil "github.com/containers-ai/alameda/pkg/utils/log"
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	openshift_apps_v1 "github.com/openshift/api/apps/v1"
	"github.com/pkg/errors"
	"google.golang.org/genproto/googleapis/rpc/code"
	apps_v1 "k8s.io/api/apps/v1"
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
	checkCycle              int64
	datahubClnt             datahub_v1alpha1.DatahubServiceClient
	k8sClienit              client.Client
	evictCfg                Config
	purgeContainerCPUMemory bool
}

// NewEvictioner return Evictioner instance
func NewEvictioner(checkCycle int64,
	datahubClnt datahub_v1alpha1.DatahubServiceClient,
	k8sClienit client.Client,
	evictCfg Config,
	purgeContainerCPUMemory bool) *Evictioner {
	return &Evictioner{
		checkCycle:              checkCycle,
		datahubClnt:             datahubClnt,
		k8sClienit:              k8sClienit,
		evictCfg:                evictCfg,
		purgeContainerCPUMemory: purgeContainerCPUMemory,
	}
}

// Start checking pods need to apply recommendation
func (evictioner *Evictioner) Start() {
	go evictioner.evictProcess()
}

func (evictioner *Evictioner) evictProcess() {
	for {
		if !evictioner.evictCfg.Enable {
			scope.Warn("evictioner is not enabled")
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
		if evictioner.purgeContainerCPUMemory {
			topController := recPod.TopController
			if topController == nil || topController.NamespacedName == nil {
				scope.Errorf("Purge pod (%s,%s) resources failed: get empty topController from PodRecommendation", recPodIns.GetNamespace(), recPodIns.GetName())
				continue

			}
			topControllerNamespace := topController.NamespacedName.Namespace
			topControllerName := topController.NamespacedName.Name
			topControllerKind := topController.Kind
			topControllerInstance, err := evictioner.getTopController(topControllerNamespace, topControllerName, topControllerKind)
			if err != nil {
				scope.Errorf("Purge pod (%s,%s) resources failed: get topController failed: %s", recPodIns.GetNamespace(), recPodIns.GetName(), err.Error())
				continue
			}
			if needToPurge, err := evictioner.needToPurgeTopControllerContainerResources(topControllerInstance, topControllerKind); err != nil {
				scope.Errorf("Purge pod (%s,%s) resources failed: %s", recPodIns.GetNamespace(), recPodIns.GetName(), err.Error())
			} else if needToPurge {
				if err = evictioner.purgeTopControllerContainerResources(topControllerInstance, topControllerKind); err != nil {
					scope.Errorf("Purge pod (%s,%s) resources failed: %s", recPodIns.GetNamespace(), recPodIns.GetName(), err.Error())
					continue
				}
			} else {
				err = evictioner.k8sClienit.Delete(context.TODO(), recPodIns)
				if err != nil {
					scope.Errorf("Evict pod (%s,%s) failed: %s", recPodIns.GetNamespace(), recPodIns.GetName(), err.Error())
				}
			}
		} else {
			err = evictioner.k8sClienit.Delete(context.TODO(), recPodIns)
			if err != nil {
				scope.Errorf("Evict pod (%s,%s) failed: %s", recPodIns.GetNamespace(), recPodIns.GetName(), err.Error())
			}
		}
	}
}

func (evictioner *Evictioner) getTopController(namespace string, name string, kind datahub_v1alpha1.Kind) (interface{}, error) {

	getResource := utilsresource.NewGetResource(evictioner.k8sClienit)

	switch kind {
	case datahub_v1alpha1.Kind_DEPLOYMENT:
		return getResource.GetDeployment(namespace, name)
	case datahub_v1alpha1.Kind_DEPLOYMENTCONFIG:
		return getResource.GetDeploymentConfig(namespace, name)
	default:
		return nil, errors.Errorf("not supported controller type %s", datahub_v1alpha1.Kind_name[int32(kind)])
	}
}

func (evictioner *Evictioner) needToPurgeTopControllerContainerResources(controller interface{}, kind datahub_v1alpha1.Kind) (bool, error) {

	switch kind {
	case datahub_v1alpha1.Kind_DEPLOYMENT:
		deployment := controller.(*apps_v1.Deployment)
		for _, container := range deployment.Spec.Template.Spec.Containers {
			resourceLimits := container.Resources.Limits
			if resourceLimits != nil {
				_, cpuSpecExist := resourceLimits[corev1.ResourceCPU]
				_, memorySpecExist := resourceLimits[corev1.ResourceMemory]
				if cpuSpecExist || memorySpecExist {
					return true, nil
				}
			}
			resourceRequests := container.Resources.Requests
			if resourceRequests != nil {
				_, cpuSpecExist := resourceRequests[corev1.ResourceCPU]
				_, memorySpecExist := resourceRequests[corev1.ResourceMemory]
				if cpuSpecExist || memorySpecExist {
					return true, nil
				}
			}
		}
		return false, nil
	case datahub_v1alpha1.Kind_DEPLOYMENTCONFIG:
		deploymentConfig := controller.(*openshift_apps_v1.DeploymentConfig)
		for _, container := range deploymentConfig.Spec.Template.Spec.Containers {
			resourceLimits := container.Resources.Limits
			if resourceLimits != nil {
				_, cpuSpecExist := resourceLimits[corev1.ResourceCPU]
				_, memorySpecExist := resourceLimits[corev1.ResourceMemory]
				if cpuSpecExist || memorySpecExist {
					return true, nil
				}
			}
			resourceRequests := container.Resources.Requests
			if resourceRequests != nil {
				_, cpuSpecExist := resourceRequests[corev1.ResourceCPU]
				_, memorySpecExist := resourceRequests[corev1.ResourceMemory]
				if cpuSpecExist || memorySpecExist {
					return true, nil
				}
			}
		}
		return false, nil
	default:
		return false, errors.Errorf("not supported controller type %s", datahub_v1alpha1.Kind_name[int32(kind)])
	}
}

func (evictioner *Evictioner) purgeTopControllerContainerResources(controller interface{}, kind datahub_v1alpha1.Kind) error {

	switch kind {
	case datahub_v1alpha1.Kind_DEPLOYMENT:
		deployment := controller.(*apps_v1.Deployment)
		deploymentCopy := deployment.DeepCopy()
		for _, container := range deploymentCopy.Spec.Template.Spec.Containers {
			resourceLimits := container.Resources.Limits
			if resourceLimits != nil {
				delete(resourceLimits, corev1.ResourceCPU)
				delete(resourceLimits, corev1.ResourceMemory)
			}
			resourceRequests := container.Resources.Requests
			if resourceRequests != nil {
				delete(resourceRequests, corev1.ResourceCPU)
				delete(resourceRequests, corev1.ResourceMemory)
			}
		}
		ctx := context.TODO()
		err := evictioner.k8sClienit.Update(ctx, deploymentCopy)
		if err != nil {
			return errors.Wrapf(err, "purge topController failed: %s", err.Error())
		}
		return nil
	case datahub_v1alpha1.Kind_DEPLOYMENTCONFIG:
		deploymentConfig := controller.(*openshift_apps_v1.DeploymentConfig)
		deploymentConfigCopy := deploymentConfig.DeepCopy()
		for _, container := range deploymentConfigCopy.Spec.Template.Spec.Containers {
			resourceLimits := container.Resources.Limits
			if resourceLimits != nil {
				delete(resourceLimits, corev1.ResourceCPU)
				delete(resourceLimits, corev1.ResourceMemory)
			}
			resourceRequests := container.Resources.Requests
			if resourceRequests != nil {
				delete(resourceRequests, corev1.ResourceCPU)
				delete(resourceRequests, corev1.ResourceMemory)
			}
		}
		ctx := context.TODO()
		err := evictioner.k8sClienit.Update(ctx, deploymentConfigCopy)
		if err != nil {
			return errors.Wrapf(err, "purge topController failed: %s", err.Error())
		}
		return nil
	default:
		return errors.Errorf("not supported controller type %s", datahub_v1alpha1.Kind_name[int32(kind)])
	}
}

func (evictioner *Evictioner) listAppliablePodRecommendation() ([]*datahub_v1alpha1.PodRecommendation, error) {

	appliablePodRecList := []*datahub_v1alpha1.PodRecommendation{}
	nowTime := time.Now()
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

	controllerRecommendationInfoMap := NewControllerRecommendationInfoMap(evictioner.k8sClienit, podRecommsPossibleToApply)
	for _, controllerRecommendationInfo := range controllerRecommendationInfoMap {
		podRecommendationInfos := controllerRecommendationInfo.podRecommendationInfos
		sort.Slice(podRecommendationInfos, func(i, j int) bool {
			return podRecommendationInfos[i].pod.ObjectMeta.CreationTimestamp.UnixNano() < podRecommendationInfos[j].pod.ObjectMeta.CreationTimestamp.UnixNano()
		})
	}
	for _, controllerRecommendationInfo := range controllerRecommendationInfoMap {

		// Create eviction restriction
		maxUnavailable := controllerRecommendationInfo.getMaxUnavailable()
		triggerThreshold, err := controllerRecommendationInfo.buildTriggerThreshold()
		if err != nil {
			scope.Errorf("Build triggerThreshold of controller (%s/%s, kind: %s) faild, skip evicting controller's pod: %s",
				controllerRecommendationInfo.namespace, controllerRecommendationInfo.name, controllerRecommendationInfo.kind, err.Error())
			continue
		}
		podRecommendations := make([]*datahub_v1alpha1.PodRecommendation, len(controllerRecommendationInfo.podRecommendationInfos))
		for i := range controllerRecommendationInfo.podRecommendationInfos {
			podRecommendations[i] = controllerRecommendationInfo.podRecommendationInfos[i].recommendation
		}
		evictionRestriction := NewEvictionRestriction(evictioner.k8sClienit, maxUnavailable, triggerThreshold, podRecommendations)

		for _, podRecommendationInfo := range controllerRecommendationInfo.podRecommendationInfos {
			pod := podRecommendationInfo.pod
			podRecommendation := podRecommendationInfo.recommendation
			if !controllerRecommendationInfo.isScalingToolTypeVPA() {
				scope.Infof("Pod (%s/%s) cannot be evicted due to AlamedaScaler's scaling tool is type of %s",
					pod.GetNamespace(), pod.GetName(), controllerRecommendationInfo.alamedaScaler.Spec.ScalingTool.Type)
				continue
			}
			if ok, err := podRecommendationInfo.isApplicableAtTime(nowTime); err != nil {
				scope.Infof("Pod (%s/%s) cannot be evicted due to PodRecommendation validate error, %s",
					pod.GetNamespace(), pod.GetName(), err.Error())
				continue
			} else if !ok {
				scope.Infof("Pod (%s/%s) cannot be evicted due to current time (%d) is not applicable on PodRecommendation's startTime (%d) and endTime(%d) interval",
					pod.GetNamespace(), pod.GetName(), nowTime.Unix(), podRecommendation.GetStartTime().GetSeconds(), podRecommendation.GetEndTime().GetSeconds())
				continue
			}
			if isEvictabel, err := evictionRestriction.IsEvictabel(pod); err != nil {
				scope.Infof("Pod (%s/%s) cannot be evicted due to eviction restriction checking error: %s", pod.GetNamespace(), pod.GetName(), err.Error())
				continue
			} else if !isEvictabel {
				scope.Infof("Pod (%s/%s) cannot be evicted.", pod.GetNamespace(), pod.GetName())
				continue
			} else {
				scope.Infof("Pod (%s/%s) can be evicted.", pod.GetNamespace(), pod.GetName())
				appliablePodRecList = append(appliablePodRecList, podRecommendation)
			}
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
			enableScalerMap[fmt.Sprintf("%s/%s", alamRecomm.GetNamespace(), or.Name)] = scaler.IsEnableExecution()
			return scaler.IsEnableExecution()
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

type podRecommendationInfo struct {
	pod            *corev1.Pod
	recommendation *datahub_v1alpha1.PodRecommendation
}

func (p *podRecommendationInfo) isApplicableAtTime(t time.Time) (bool, error) {

	startTime := p.recommendation.GetStartTime()
	endTime := p.recommendation.GetEndTime()

	if startTime == nil || endTime == nil {
		return false, errors.Errorf("starTime and endTime cannot be nil")
	}

	if startTime.GetSeconds() >= t.Unix() || t.Unix() >= endTime.GetSeconds() {
		return false, nil
	}

	return true, nil
}

type controllerRecommendationInfo struct {
	namespace              string
	name                   string
	kind                   string
	alamedaScaler          *autoscalingv1alpha1.AlamedaScaler
	podRecommendationInfos []*podRecommendationInfo
}

func (c controllerRecommendationInfo) getMaxUnavailable() string {

	var maxUnavailable string

	scalingTool := c.alamedaScaler.Spec.ScalingTool
	if scalingTool.ExecutionStrategy == nil {
		maxUnavailable = autoscalingv1alpha1.DefaultMaxUnavailablePercentage
		return maxUnavailable
	}

	maxUnavailable = scalingTool.ExecutionStrategy.MaxUnavailable
	return maxUnavailable
}

func (c controllerRecommendationInfo) isScalingToolTypeVPA() bool {
	return c.alamedaScaler.IsScalingToolTypeVPA()
}

func (c controllerRecommendationInfo) buildTriggerThreshold() (triggerThreshold, error) {

	var triggerThreshold triggerThreshold

	cpu := c.alamedaScaler.Spec.ScalingTool.ExecutionStrategy.TriggerThreshold.CPU
	cpu = strings.TrimSuffix(cpu, "%")
	cpuValue, err := strconv.ParseFloat(cpu, 64)
	if err != nil {
		return triggerThreshold, errors.Errorf("parse cpu trigger threshold failed: %s", err.Error())
	}
	triggerThreshold.CPU = cpuValue

	memory := c.alamedaScaler.Spec.ScalingTool.ExecutionStrategy.TriggerThreshold.Memory
	memory = strings.TrimSuffix(memory, "%")
	memoryValue, err := strconv.ParseFloat(memory, 64)
	if err != nil {
		return triggerThreshold, errors.Errorf("parse memory trigger threshold failed: %s", err.Error())
	}
	triggerThreshold.Memory = memoryValue

	return triggerThreshold, nil
}

func NewControllerRecommendationInfoMap(client client.Client, podRecommendations []*datahub_v1alpha1.PodRecommendation) map[string]*controllerRecommendationInfo {

	getResource := utilsresource.NewGetResource(client)
	alamedaScalerMap := make(map[string]*autoscalingv1alpha1.AlamedaScaler)
	controllerRecommendationInfoMap := make(map[string]*controllerRecommendationInfo)
	for _, podRecommendation := range podRecommendations {

		// Filter out invalid PodRecommendation
		copyPodRecommendation := proto.Clone(podRecommendation)
		podRecommendation = copyPodRecommendation.(*datahub_v1alpha1.PodRecommendation)
		recommendationNamespacedName := podRecommendation.NamespacedName
		if recommendationNamespacedName == nil {
			scope.Errorf("skip PodRecommendation due to PodRecommendation has empty NamespacedName")
			continue
		}

		// Get AlamedaScaler owns this PodRecommendation and validate the AlamedaScaler is enabled execution.
		alamedaRecommendation, err := getResource.GetAlamedaRecommendation(podRecommendation.NamespacedName.Namespace, podRecommendation.NamespacedName.Name)
		if err != nil {
			scope.Errorf("skip PodRecommendation (%s/%s) due to get AlamedaRecommendation falied: %s", podRecommendation.NamespacedName.Namespace, podRecommendation.NamespacedName.Name, err.Error())
			continue
		}
		alamedaScalerNamespace := ""
		alamedaScalerName := ""
		for _, or := range alamedaRecommendation.OwnerReferences {
			if or.Kind == "AlamedaScaler" {
				alamedaScalerNamespace = alamedaRecommendation.Namespace
				alamedaScalerName = or.Name
				break
			}
		}
		alamedaScaler, exist := alamedaScalerMap[fmt.Sprintf("%s/%s", alamedaScalerNamespace, alamedaScalerName)]
		if !exist {
			alamedaScaler, err = getResource.GetAlamedaScaler(alamedaScalerNamespace, alamedaScalerName)
			if err != nil {
				scope.Errorf("skip PodRecommendation (%s/%s) due to get AlamedaScaler falied: %s", podRecommendation.NamespacedName.Namespace, podRecommendation.NamespacedName.Name, err.Error())
				continue
			}
			alamedaScalerMap[fmt.Sprintf("%s/%s", alamedaScalerNamespace, alamedaScalerName)] = alamedaScaler
		}
		if !alamedaScaler.IsEnableExecution() {
			scope.Errorf("skip PodRecommendation (%s/%s) because it's execution is not enabled.", recommendationNamespacedName.Namespace, recommendationNamespacedName.Name)
			continue
		}

		// Get Pod instance of this PodRecommendation
		podNamespace := recommendationNamespacedName.Namespace
		podName := recommendationNamespacedName.Name
		pod, err := getResource.GetPod(podNamespace, podName)
		if err != nil {
			scope.Errorf("skip PodRecommendation due to get Pod (%s/%s) failed: %s", podNamespace, podName, err.Error())
			continue
		}

		// Get topmost controller namespace, name and kind controlling this pod
		controller := podRecommendation.TopController
		if controller == nil {
			scope.Errorf("skip PodRecommendation (%s/%s) due to PodRecommendation has empty topmost controller", recommendationNamespacedName.Namespace, recommendationNamespacedName.Name)
			continue
		} else if controller.NamespacedName == nil {
			scope.Errorf("skip PodRecommendation (%s/%s) due to topmost controller has empty NamespacedName", recommendationNamespacedName.Namespace, recommendationNamespacedName.Name)
			continue
		}

		// Append podRecommendationInfos into controllerRecommendationInfo
		controllerID := fmt.Sprintf("%s.%s.%s", controller.Kind, controller.NamespacedName.Namespace, controller.NamespacedName.Name)
		_, exist = controllerRecommendationInfoMap[controllerID]
		if !exist {
			controllerRecommendationInfoMap[controllerID] = &controllerRecommendationInfo{
				namespace:              controller.NamespacedName.Namespace,
				name:                   controller.NamespacedName.Name,
				kind:                   datahub_v1alpha1.Kind_name[int32(controller.Kind)],
				alamedaScaler:          alamedaScaler,
				podRecommendationInfos: make([]*podRecommendationInfo, 0),
			}
		}
		podRecommendationInfo := &podRecommendationInfo{
			pod:            pod,
			recommendation: podRecommendation,
		}
		controllerRecommendationInfoMap[controllerID].podRecommendationInfos = append(
			controllerRecommendationInfoMap[controllerID].podRecommendationInfos,
			podRecommendationInfo,
		)
	}

	return controllerRecommendationInfoMap
}
