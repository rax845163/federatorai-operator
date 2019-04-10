package updateresource

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var (
	//the source file has no value
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
	//okd default value
	okdServiceDefaultProtocol                 corev1.Protocol = corev1.ProtocolTCP
	okdServiceDefaultTargetPort                               = func(port int32) int32 { return port }
	okdDeploymentDefaultDefaultMode           *int32          = func(i int32) *int32 { return &i }(420)
	okdDeploymentDefaultEnvFieldRefAPIVersion string          = "v1"
	log                                                       = logf.Log.WithName("controller_alamedaservice")
)

func MisMatchResourceConfigMap(clusterCM, sourceCM *corev1.ConfigMap) bool {
	modify := false
	if !equality.Semantic.DeepEqual(clusterCM.Data, sourceCM.Data) {
		modify = true
		log.V(-1).Info("change Data")
		clusterCM.Data = sourceCM.Data
	}
	return modify
}
func MisMatchResourceService(clusterSv, sourceSv *corev1.Service) bool {
	modify := false
	if !equality.Semantic.DeepEqual(clusterSv.Labels, sourceSv.Labels) {
		modify = true
		log.V(-1).Info("change Labels")
		clusterSv.Labels = sourceSv.Labels
	}
	for index, value := range sourceSv.Spec.Ports {
		if resourceEmpty(value.Protocol) {
			sourceSv.Spec.Ports[index].Protocol = okdServiceDefaultProtocol
		}
		if resourceEmpty(value.TargetPort.IntVal) {
			sourceSv.Spec.Ports[index].TargetPort.IntVal = okdServiceDefaultTargetPort(sourceSv.Spec.Ports[index].Port)
		}
	}

	if !equality.Semantic.DeepEqual(clusterSv.Spec.Ports, sourceSv.Spec.Ports) {
		modify = true

		log.V(-1).Info("change Ports")
		clusterSv.Spec.Ports = sourceSv.Spec.Ports
	}
	if !equality.Semantic.DeepEqual(clusterSv.Spec.Selector, sourceSv.Spec.Selector) {
		modify = true
		log.V(-1).Info("change Selector")
		clusterSv.Spec.Selector = sourceSv.Spec.Selector
	}
	return modify
}

func MisMatchResourceDeployment(clusterDep, sourceDep *appsv1.Deployment) bool {
	modify := false
	if !equality.Semantic.DeepEqual(clusterDep.Labels, sourceDep.Labels) {
		modify = true
		log.V(-1).Info("change Labels")
		clusterDep.Labels = sourceDep.Labels
	}
	misMatchSelectorAndReplicas(&modify, &clusterDep.Spec, &sourceDep.Spec)
	misMatchTemplate(&modify, &clusterDep.Spec.Template, &sourceDep.Spec.Template)
	return modify
}
func misMatchSelectorAndReplicas(modify *bool, clusterDep, sourceDep *appsv1.DeploymentSpec) {
	if !equality.Semantic.DeepEqual(clusterDep.Selector, sourceDep.Selector) {
		*modify = true
		log.V(-1).Info("change Selector")
		clusterDep.Selector = sourceDep.Selector
	}
	if !equality.Semantic.DeepEqual(clusterDep.Replicas, sourceDep.Replicas) {
		*modify = true
		log.V(-1).Info("change Replicas")
		clusterDep.Replicas = sourceDep.Replicas
	}
}
func misMatchTemplate(modify *bool, clusterDep, sourceDep *corev1.PodTemplateSpec) {
	misMatchTemplateObjectMeta(modify, &clusterDep.ObjectMeta, &sourceDep.ObjectMeta)
	misMatchTemplatePodSpec(modify, &clusterDep.Spec, &sourceDep.Spec)
}
func misMatchTemplateObjectMeta(modify *bool, clusterDep, sourceDep *metav1.ObjectMeta) {
	if clusterDep.Name != sourceDep.Name {
		*modify = true
		log.V(-1).Info("change TemplateObjectMetaName")
		clusterDep.Name = sourceDep.Name
	}
	if !equality.Semantic.DeepEqual(clusterDep.Labels, sourceDep.Labels) {
		*modify = true
		log.V(-1).Info("change TemplateObjectMetaLabels")
		clusterDep.Labels = sourceDep.Labels
	}
}
func misMatchTemplatePodSpec(modify *bool, clusterDep, sourceDep *corev1.PodSpec) {
	if clusterDep.ServiceAccountName != sourceDep.ServiceAccountName {
		*modify = true
		log.V(-1).Info("change ServiceAccountName")
		clusterDep.ServiceAccountName = sourceDep.ServiceAccountName
	}
	for sourceIndex, sourceContainerValue := range sourceDep.Containers {
		for clusterIndex, clusterContainerValue := range clusterDep.Containers {
			if clusterContainerValue.Name == sourceContainerValue.Name {
				if clusterDep.Containers[clusterIndex].Image != sourceDep.Containers[sourceIndex].Image {
					*modify = true
					log.V(-1).Info("change Image")
					clusterDep.Containers[clusterIndex].Image = sourceDep.Containers[sourceIndex].Image
				}
				if clusterDep.Containers[clusterIndex].ImagePullPolicy != sourceDep.Containers[sourceIndex].ImagePullPolicy {
					*modify = true
					log.V(-1).Info("change ImagePullPolicy")
					clusterDep.Containers[clusterIndex].ImagePullPolicy = sourceDep.Containers[sourceIndex].ImagePullPolicy
				}
				if !equality.Semantic.DeepEqual(clusterDep.Containers[clusterIndex].Ports, sourceDep.Containers[sourceIndex].Ports) {
					*modify = true
					log.V(-1).Info("change Ports")
					clusterDep.Containers[clusterIndex].Ports = sourceDep.Containers[sourceIndex].Ports
				}
				if !equality.Semantic.DeepEqual(clusterDep.Containers[clusterIndex].Resources, sourceDep.Containers[sourceIndex].Resources) {
					*modify = true
					log.V(-1).Info("change Resources")
					clusterDep.Containers[clusterIndex].Resources = sourceDep.Containers[sourceIndex].Resources
				}
				if !equality.Semantic.DeepEqual(clusterDep.Containers[clusterIndex].VolumeMounts, sourceDep.Containers[sourceIndex].VolumeMounts) {
					*modify = true
					log.V(-1).Info("change VolumeMounts")
					clusterDep.Containers[clusterIndex].VolumeMounts = sourceDep.Containers[sourceIndex].VolumeMounts
				}
				for index, value := range sourceDep.Containers[sourceIndex].Env {
					if value.ValueFrom != nil {
						if value.ValueFrom.FieldRef != nil {
							if resourceEmpty(value.ValueFrom.FieldRef.APIVersion) {
								sourceDep.Containers[sourceIndex].Env[index].ValueFrom.FieldRef.APIVersion = okdDeploymentDefaultEnvFieldRefAPIVersion
							}
						}
					}
				}
				if !equality.Semantic.DeepEqual(clusterDep.Containers[clusterIndex].Env, sourceDep.Containers[sourceIndex].Env) {
					*modify = true
					log.V(-1).Info("change Env")
					clusterDep.Containers[clusterIndex].Env = sourceDep.Containers[sourceIndex].Env
				}
			}
		}
	}
	for sourceIndex, sourceVolumeValue := range sourceDep.Volumes {
		for clusterIndex, clusterVolumeValue := range clusterDep.Volumes {
			if clusterVolumeValue.Name == sourceVolumeValue.Name {
				if sourceDep.Volumes[sourceIndex].VolumeSource.Secret != nil {
					if resourceEmpty(sourceDep.Volumes[sourceIndex].VolumeSource.Secret.DefaultMode) {
						sourceDep.Volumes[sourceIndex].VolumeSource.Secret.DefaultMode = okdDeploymentDefaultDefaultMode
					}
				}
				if sourceDep.Volumes[sourceIndex].VolumeSource.ConfigMap != nil {
					if resourceEmpty(sourceDep.Volumes[sourceIndex].VolumeSource.ConfigMap.DefaultMode) {
						sourceDep.Volumes[sourceIndex].VolumeSource.ConfigMap.DefaultMode = okdDeploymentDefaultDefaultMode
					}
				}
				if sourceDep.Volumes[sourceIndex].VolumeSource.DownwardAPI != nil {
					if resourceEmpty(sourceDep.Volumes[sourceIndex].VolumeSource.DownwardAPI.DefaultMode) {
						sourceDep.Volumes[sourceIndex].VolumeSource.DownwardAPI.DefaultMode = okdDeploymentDefaultDefaultMode
					}
					for itemsIndex, _ := range sourceDep.Volumes[sourceIndex].VolumeSource.DownwardAPI.Items {
						if sourceDep.Volumes[sourceIndex].VolumeSource.DownwardAPI.Items[itemsIndex].FieldRef != nil {
							if resourceEmpty(sourceDep.Volumes[sourceIndex].VolumeSource.DownwardAPI.Items[itemsIndex].FieldRef.APIVersion) {
								sourceDep.Volumes[sourceIndex].VolumeSource.DownwardAPI.Items[itemsIndex].FieldRef.APIVersion = okdDeploymentDefaultEnvFieldRefAPIVersion
							}
						}
					}
				}
				if !equality.Semantic.DeepEqual(clusterDep.Volumes[clusterIndex].VolumeSource, sourceDep.Volumes[sourceIndex].VolumeSource) {
					*modify = true
					log.V(-1).Info("change VolumeSource")
					clusterDep.Volumes[clusterIndex].VolumeSource = sourceDep.Volumes[sourceIndex].VolumeSource
				}
			}
		}
	}
}
