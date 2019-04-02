package updateresource

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var (
	resourceEmpty = func(value interface{}) bool {
		switch v := value.(type) {
		case corev1.Protocol:
			{
				if v == "" {
					return true
				}
				return false
			}
		case string:
			if v == "" {
				return true
			}
			return false
		case int32:
			if v == 0 {
				return true
			}
			return false
		case *int32:
			if v == nil {
				return true
			}
			return false
		}
		return false
	}
	okdServiceDefaultProtocol                 corev1.Protocol = corev1.ProtocolTCP
	okdServiceDefaultTargetPort                               = func(port int32) int32 { return port }
	okdDeploymentDefaultDefaultMode           *int32          = func(i int32) *int32 { return &i }(420)
	okdDeploymentDefaultEnvFieldRefAPIVersion string          = "v1"
	log                                                       = logf.Log.WithName("controller_alamedaservice")
)

func MatchResourceService(foundSv, resourceSv *corev1.Service) bool {
	modify := false
	if !equality.Semantic.DeepEqual(foundSv.Labels, resourceSv.Labels) {
		modify = true
		log.V(-1).Info("change Labels")
		foundSv.Labels = resourceSv.Labels
	}
	for index, value := range resourceSv.Spec.Ports {
		if resourceEmpty(value.Protocol) {
			resourceSv.Spec.Ports[index].Protocol = okdServiceDefaultProtocol
		}
		if resourceEmpty(value.TargetPort.IntVal) {
			resourceSv.Spec.Ports[index].TargetPort.IntVal = okdServiceDefaultTargetPort(resourceSv.Spec.Ports[index].Port)
		}
	}
	if !equality.Semantic.DeepEqual(foundSv.Spec.Ports, resourceSv.Spec.Ports) {
		modify = true

		log.V(-1).Info("change Ports")
		foundSv.Spec.Ports = resourceSv.Spec.Ports
	}
	if !equality.Semantic.DeepEqual(foundSv.Spec.Selector, resourceSv.Spec.Selector) {
		modify = true
		log.V(-1).Info("change Selector")
		foundSv.Spec.Selector = resourceSv.Spec.Selector
	}
	return modify
}

func MatchResourceDeployment(foundDep, resourceDep *appsv1.Deployment) bool {
	modify := false
	if !equality.Semantic.DeepEqual(foundDep.Labels, resourceDep.Labels) {
		modify = true
		log.V(-1).Info("change Labels")
		foundDep.Labels = resourceDep.Labels
	}
	matchSelectorAndReplicas(&modify, &foundDep.Spec, &resourceDep.Spec)
	matchTemplate(&modify, &foundDep.Spec.Template, &resourceDep.Spec.Template)
	return modify
}
func matchSelectorAndReplicas(modify *bool, foundDep, resourceDep *appsv1.DeploymentSpec) {
	if !equality.Semantic.DeepEqual(foundDep.Selector, resourceDep.Selector) {
		*modify = true
		log.V(-1).Info("change Selector")
		foundDep.Selector = resourceDep.Selector
	}
	if !equality.Semantic.DeepEqual(foundDep.Replicas, resourceDep.Replicas) {
		*modify = true
		log.V(-1).Info("change Replicas")
		foundDep.Replicas = resourceDep.Replicas
	}
}
func matchTemplate(modify *bool, foundDep, resourceDep *corev1.PodTemplateSpec) {
	matchTemplateObjectMeta(modify, &foundDep.ObjectMeta, &resourceDep.ObjectMeta)
	matchTemplatePodSpec(modify, &foundDep.Spec, &resourceDep.Spec)
}
func matchTemplateObjectMeta(modify *bool, foundDep, resourceDep *metav1.ObjectMeta) {
	if foundDep.Name != resourceDep.Name {
		*modify = true
		log.V(-1).Info("change TemplateObjectMetaName")
		foundDep.Name = resourceDep.Name
	}
	if !equality.Semantic.DeepEqual(foundDep.Labels, resourceDep.Labels) {
		*modify = true
		log.V(-1).Info("change TemplateObjectMetaLabels")
		foundDep.Labels = resourceDep.Labels
	}
}
func matchTemplatePodSpec(modify *bool, foundDep, resourceDep *corev1.PodSpec) {
	if foundDep.ServiceAccountName != resourceDep.ServiceAccountName {
		*modify = true
		log.V(-1).Info("change ServiceAccountName")
		foundDep.ServiceAccountName = resourceDep.ServiceAccountName
	}
	for resourceIndex, resourceContainerValue := range resourceDep.Containers {
		for foundIndex, foundContainerValue := range foundDep.Containers {
			if foundContainerValue.Name == resourceContainerValue.Name {
				if foundDep.Containers[foundIndex].Image != resourceDep.Containers[resourceIndex].Image {
					*modify = true
					log.V(-1).Info("change Image")
					foundDep.Containers[foundIndex].Image = resourceDep.Containers[resourceIndex].Image
				}
				if foundDep.Containers[foundIndex].ImagePullPolicy != resourceDep.Containers[resourceIndex].ImagePullPolicy {
					*modify = true
					log.V(-1).Info("change ImagePullPolicy")
					foundDep.Containers[foundIndex].ImagePullPolicy = resourceDep.Containers[resourceIndex].ImagePullPolicy
				}
				if !equality.Semantic.DeepEqual(foundDep.Containers[foundIndex].Ports, resourceDep.Containers[resourceIndex].Ports) {
					*modify = true
					log.V(-1).Info("change Ports")
					foundDep.Containers[foundIndex].Ports = resourceDep.Containers[resourceIndex].Ports
				}
				if !equality.Semantic.DeepEqual(foundDep.Containers[foundIndex].Resources, resourceDep.Containers[resourceIndex].Resources) {
					*modify = true
					log.V(-1).Info("change Resources")
					foundDep.Containers[foundIndex].Resources = resourceDep.Containers[resourceIndex].Resources
				}
				if !equality.Semantic.DeepEqual(foundDep.Containers[foundIndex].VolumeMounts, resourceDep.Containers[resourceIndex].VolumeMounts) {
					*modify = true
					log.V(-1).Info("change VolumeMounts")
					foundDep.Containers[foundIndex].VolumeMounts = resourceDep.Containers[resourceIndex].VolumeMounts
				}
				for index, value := range resourceDep.Containers[resourceIndex].Env {
					if value.ValueFrom != nil {
						if value.ValueFrom.FieldRef != nil {
							if resourceEmpty(value.ValueFrom.FieldRef.APIVersion) {
								resourceDep.Containers[resourceIndex].Env[index].ValueFrom.FieldRef.APIVersion = okdDeploymentDefaultEnvFieldRefAPIVersion
							}
						}
					}
				}
				if !equality.Semantic.DeepEqual(foundDep.Containers[foundIndex].Env, resourceDep.Containers[resourceIndex].Env) {
					*modify = true
					log.V(-1).Info("change Env")
					foundDep.Containers[foundIndex].Env = resourceDep.Containers[resourceIndex].Env
				}
			}
		}
	}
	for resourceIndex, resourceVolumeValue := range resourceDep.Volumes {
		for foundIndex, foundVolumeValue := range foundDep.Volumes {
			if foundVolumeValue.Name == resourceVolumeValue.Name {
				if resourceDep.Volumes[resourceIndex].VolumeSource.Secret != nil {
					if resourceEmpty(resourceDep.Volumes[resourceIndex].VolumeSource.Secret.DefaultMode) {
						resourceDep.Volumes[resourceIndex].VolumeSource.Secret.DefaultMode = okdDeploymentDefaultDefaultMode
					}
				}
				if resourceDep.Volumes[resourceIndex].VolumeSource.ConfigMap != nil {
					if resourceEmpty(resourceDep.Volumes[resourceIndex].VolumeSource.ConfigMap.DefaultMode) {
						resourceDep.Volumes[resourceIndex].VolumeSource.ConfigMap.DefaultMode = okdDeploymentDefaultDefaultMode
					}
				}
				if !equality.Semantic.DeepEqual(foundDep.Volumes[foundIndex].VolumeSource, resourceDep.Volumes[resourceIndex].VolumeSource) {
					*modify = true
					log.V(-1).Info("change VolumeSource")
					foundDep.Volumes[foundIndex].VolumeSource = resourceDep.Volumes[resourceIndex].VolumeSource
				}
			}
		}
	}
}
