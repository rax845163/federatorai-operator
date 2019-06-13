package alamedarecommendation

import (
	"github.com/containers-ai/alameda/datahub/pkg/utils"
	autoscaling_v1alpha1 "github.com/containers-ai/alameda/operator/pkg/apis/autoscaling/v1alpha1"
	logUtil "github.com/containers-ai/alameda/pkg/utils/log"
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	alamedarecommendationReconcilerScope = logUtil.RegisterScope("alamedarecommendation_reconciler", "alamedarecommendation_reconciler", 0)
)

// Reconciler reconciles AlamedaRecommendation object
type Reconciler struct {
	client                client.Client
	alamedaRecommendation *autoscaling_v1alpha1.AlamedaRecommendation
}

// NewReconciler creates Reconciler object
func NewReconciler(client client.Client, alamedaRecommendation *autoscaling_v1alpha1.AlamedaRecommendation) *Reconciler {
	return &Reconciler{
		client:                client,
		alamedaRecommendation: alamedaRecommendation,
	}
}

// UpdateResourceRecommendation updates resource of AlamedaRecommendation
func (reconciler *Reconciler) UpdateResourceRecommendation(podRecommendation *datahub_v1alpha1.PodRecommendation) (*autoscaling_v1alpha1.AlamedaRecommendation, error) {
	for alamedaContainerIdx, alamedaContainer := range reconciler.alamedaRecommendation.Spec.Containers {
		for _, containerRecommendation := range podRecommendation.ContainerRecommendations {
			if alamedaContainer.Name == containerRecommendation.Name {
				if alamedaContainer.Resources.Limits == nil {
					alamedaContainer.Resources.Limits = corev1.ResourceList{}
				}
				if alamedaContainer.Resources.Requests == nil {
					alamedaContainer.Resources.Requests = corev1.ResourceList{}
				}
				for _, limitRecommendation := range containerRecommendation.LimitRecommendations {
					if limitRecommendation.MetricType == datahub_v1alpha1.MetricType_CPU_USAGE_SECONDS_PERCENTAGE {
						cpuLimitTime := int64(0)
						for _, data := range limitRecommendation.Data {
							curNanoSec := utils.TimeStampToNanoSecond(data.Time)
							if numVal, err := utils.StringToInt64(data.NumValue); err == nil && curNanoSec > cpuLimitTime {
								alamedaContainer.Resources.Limits[corev1.ResourceCPU] = *resource.NewMilliQuantity(numVal, resource.DecimalSI)
								cpuLimitTime = curNanoSec
							} else if err != nil {
								alamedarecommendationReconcilerScope.Error(err.Error())
							}
						}

					} else if limitRecommendation.MetricType == datahub_v1alpha1.MetricType_MEMORY_USAGE_BYTES {
						memoryLimitTime := int64(0)
						for _, data := range limitRecommendation.Data {
							curNanoSec := utils.TimeStampToNanoSecond(data.Time)
							if numVal, err := utils.StringToInt64(data.NumValue); err == nil && curNanoSec > memoryLimitTime {
								alamedaContainer.Resources.Limits[corev1.ResourceMemory] = *resource.NewQuantity(numVal, resource.BinarySI)
								memoryLimitTime = curNanoSec
							} else if err != nil {
								alamedarecommendationReconcilerScope.Error(err.Error())
							}
						}
					}
				}
				for _, requestRecommendation := range containerRecommendation.RequestRecommendations {
					if requestRecommendation.MetricType == datahub_v1alpha1.MetricType_CPU_USAGE_SECONDS_PERCENTAGE {
						cpuRequestTime := int64(0)
						for _, data := range requestRecommendation.Data {
							curNanoSec := utils.TimeStampToNanoSecond(data.Time)
							if numVal, err := utils.StringToInt64(data.NumValue); err == nil && curNanoSec > cpuRequestTime {
								alamedaContainer.Resources.Requests[corev1.ResourceCPU] = *resource.NewMilliQuantity(numVal, resource.DecimalSI)
								cpuRequestTime = curNanoSec
							} else if err != nil {
								alamedarecommendationReconcilerScope.Error(err.Error())
							}
						}

					} else if requestRecommendation.MetricType == datahub_v1alpha1.MetricType_MEMORY_USAGE_BYTES {
						memoryRequestTime := int64(0)
						for _, data := range requestRecommendation.Data {
							curNanoSec := utils.TimeStampToNanoSecond(data.Time)
							if numVal, err := utils.StringToInt64(data.NumValue); err == nil && curNanoSec > memoryRequestTime {
								alamedaContainer.Resources.Requests[corev1.ResourceMemory] = *resource.NewQuantity(numVal, resource.BinarySI)
								memoryRequestTime = curNanoSec
							} else if err != nil {
								alamedarecommendationReconcilerScope.Error(err.Error())
							}
						}
					}
				}
			}
		}
		reconciler.alamedaRecommendation.Spec.Containers[alamedaContainerIdx] = alamedaContainer
	}
	return reconciler.alamedaRecommendation, nil
}
