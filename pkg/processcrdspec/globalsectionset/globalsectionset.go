package globalsectionset

import (
	"strings"

	admission_controller "github.com/containers-ai/alameda/admission-controller"
	"github.com/containers-ai/federatorai-operator/pkg/apis/federatorai/v1alpha1"
	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/alamedaserviceparamter"
	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/updateenvvar"
	"github.com/containers-ai/federatorai-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func GlobalSectionSetParamterToStatefulset(ss *appsv1.StatefulSet, asp *alamedaserviceparamter.AlamedaServiceParamter) {
	switch ss.Name {
	case util.FedemeterInflixDBSSN:
		util.SetStatefulsetImageStruct(ss, asp.Version, util.FedemeterInflixDBSSN)
	}
}

func GlobalSectionSetParamterToDeployment(dep *appsv1.Deployment, asp *alamedaserviceparamter.AlamedaServiceParamter) {
	processDeploymentPrometheusService(dep, asp.PrometheusService) //Global section set DeploymentSpec's PrometheusService

	switch dep.Name {
	case util.AlamedaaiDPN:
		{
			//Global section set DeploymentSpec's image version(only alameda component)
			util.SetImageStruct(dep, asp.Version, util.AlamedaaiCTN)
			//Global section set persistentVolumeClaim to mountPath
			util.SetStorageToVolumeSource(dep, asp.Storages, "alameda-ai-type.pvc", util.AlamedaGroup)
			util.SetStorageToMountPath(dep, asp.Storages, util.AlamedaaiCTN, "alameda-ai-type-storage", util.AlamedaGroup)
		}
	case util.AlamedaoperatorDPN:
		{
			util.SetImageStruct(dep, asp.Version, util.AlamedaoperatorCTN)
			util.SetStorageToVolumeSource(dep, asp.Storages, "alameda-operator-type.pvc", util.AlamedaGroup)
			util.SetStorageToMountPath(dep, asp.Storages, util.AlamedaoperatorCTN, "alameda-operator-type-storage", util.AlamedaGroup)
		}
	case util.AlamedadatahubDPN:
		{
			util.SetImageStruct(dep, asp.Version, util.AlamedadatahubCTN)
			util.SetStorageToVolumeSource(dep, asp.Storages, "alameda-datahub-type.pvc", util.AlamedaGroup)
			util.SetStorageToMountPath(dep, asp.Storages, util.AlamedadatahubCTN, "alameda-datahub-type-storage", util.AlamedaGroup)
		}
	case util.AlamedaevictionerDPN:
		{
			util.SetImageStruct(dep, asp.Version, util.AlamedaevictionerCTN)
			util.SetStorageToVolumeSource(dep, asp.Storages, "alameda-evictioner-type.pvc", util.AlamedaGroup)
			util.SetStorageToMountPath(dep, asp.Storages, util.AlamedaevictionerCTN, "alameda-evictioner-type-storage", util.AlamedaGroup)
		}
	case util.AdmissioncontrollerDPN:
		{
			util.SetImageStruct(dep, asp.Version, util.AdmissioncontrollerCTN)
			util.SetStorageToVolumeSource(dep, asp.Storages, "admission-controller-type.pvc", util.AlamedaGroup)
			util.SetStorageToMountPath(dep, asp.Storages, util.AdmissioncontrollerCTN, "admission-controller-type-storage", util.AlamedaGroup)
		}
	case util.AlamedarecommenderDPN:
		{
			util.SetImageStruct(dep, asp.Version, util.AlamedarecommenderCTN)
			util.SetStorageToVolumeSource(dep, asp.Storages, "alameda-recommender-type.pvc", util.AlamedaGroup)
			util.SetStorageToMountPath(dep, asp.Storages, util.AlamedarecommenderCTN, "alameda-recommender-type-storage", util.AlamedaGroup)
		}
	case util.AlamedaexecutorDPN:
		{
			util.SetImageStruct(dep, asp.Version, util.AlamedaexecutorCTN)
			util.SetStorageToVolumeSource(dep, asp.Storages, "alameda-executor-type.pvc", util.AlamedaGroup)
			util.SetStorageToMountPath(dep, asp.Storages, util.AlamedaexecutorCTN, "alameda-executor-type-storage", util.AlamedaGroup)
		}
	case util.AlamedadispatcherDPN:
		{
			util.SetImageStruct(dep, asp.Version, util.AlamedadispatcherCTN)
			util.SetStorageToVolumeSource(dep, asp.Storages, "alameda-dispatcher-type.pvc", util.AlamedaGroup)
			util.SetStorageToMountPath(dep, asp.Storages, util.AlamedadispatcherCTN, "alameda-dispatcher-type-storage", util.AlamedaGroup)
		}
	case util.AlamedaRabbitMQDPN:
		util.SetImageStruct(dep, asp.Version, util.AlamedaRabbitMQCTN)
	case util.AlamedaanalyzerDPN:
		{
			util.SetImageStruct(dep, asp.Version, util.AlamedaanalyzerCTN)
			util.SetStorageToVolumeSource(dep, asp.Storages, "alameda-analyzer-type.pvc", util.AlamedaGroup)
			util.SetStorageToMountPath(dep, asp.Storages, util.AlamedaanalyzerCTN, "alameda-analyzer-type-storage", util.AlamedaGroup)
		}
	case util.FedemeterDPN:
		{
			util.SetImageStruct(dep, asp.Version, util.FedemeterCTN)
			util.SetStorageToVolumeSource(dep, asp.Storages, "fedemeter-type.pvc", util.FedemeterGroup)
			util.SetStorageToMountPath(dep, asp.Storages, util.FedemeterCTN, "fedemeter-type-storage", util.FedemeterGroup)
		}
	case util.InfluxdbDPN:
		{
			util.SetStorageToVolumeSource(dep, asp.Storages, "my-alameda.influxdb-type.pvc", util.InfluxDBGroup)
			util.SetStorageToMountPath(dep, asp.Storages, util.InfluxdbCTN, "influxdb-type-storage", util.InfluxDBGroup)
		}
	case util.GrafanaDPN:
		{
			util.SetImageStruct(dep, asp.Version, util.GrafanaCTN)
			util.SetStorageToVolumeSource(dep, asp.Storages, "my-alameda.grafana-type.pvc", util.GrafanaGroup)
			util.SetStorageToMountPath(dep, asp.Storages, util.GrafanaCTN, "grafana-type-storage", util.GrafanaGroup)
		}
	case util.AlamedaweavescopeDPN:
		{
			util.SetImageStruct(dep, asp.AlamedaWeavescopeSectionSet, util.AlamedaweavescopeCTN)
			util.SetImagePullPolicy(dep, util.AlamedaweavescopeCTN, asp.AlamedaWeavescopeSectionSet.ImagePullPolicy)
		}
	case util.AlamedaNotifierDPN:
		util.SetImageStruct(dep, asp.Version, util.AlamedaNofitierCTN)
		util.SetStorageToVolumeSource(dep, asp.Storages, "alameda-notifier-type.pvc", util.AlamedaGroup)
		util.SetStorageToMountPath(dep, asp.Storages, util.AlamedaNofitierCTN, "alameda-notifier-type-storage", util.AlamedaGroup)
	case util.FederatoraiAgentDPN:
		util.SetImageStruct(dep, asp.Version, util.FederatoraiAgentCTN)
		util.SetStorageToVolumeSource(dep, asp.Storages, "federatorai-agent-type.pvc", util.AlamedaGroup)
		util.SetStorageToMountPath(dep, asp.Storages, util.FederatoraiAgentCTN, "federatorai-agent-type-storage", util.AlamedaGroup)
	case util.FederatoraiAgentGPUDPN:
		util.SetImageStruct(dep, asp.Version, util.FederatoraiAgentGPUCTN)
		util.SetStorageToVolumeSource(dep, asp.Storages, "federatorai-agent-gpu-type.pvc", util.AlamedaGroup)
		util.SetStorageToMountPath(dep, asp.Storages, util.FederatoraiAgentGPUCTN, "federatorai-agent-gpu-type-storage", util.AlamedaGroup)
	case util.FederatoraiRESTDPN:
		util.SetImageStruct(dep, asp.Version, util.FederatoraiRESTCTN)
		util.SetStorageToVolumeSource(dep, asp.Storages, "federatorai-rest-type.pvc", util.AlamedaGroup)
		util.SetStorageToMountPath(dep, asp.Storages, util.FederatoraiRESTCTN, "federatorai-rest-type-storage", util.AlamedaGroup)
	}

	envVars := getEnvVarsToUpdateByDeployment(dep.Name, asp)
	updateenvvar.UpdateEnvVarsToDeployment(dep, envVars)
}

func GlobalSectionSetParamterToDaemonSet(ds *appsv1.DaemonSet, asp *alamedaserviceparamter.AlamedaServiceParamter) {
	switch ds.Name {
	case util.AlamedaweavescopeAgentDS:
		{
			util.SetDaemonSetImageStruct(ds, asp.AlamedaWeavescopeSectionSet, util.AlamedaweavescopeAgentCTN)
			util.SetDaemonSetImagePullPolicy(ds, util.AlamedaweavescopeAgentCTN, asp.AlamedaWeavescopeSectionSet.ImagePullPolicy)
		}
	}
}

func processConfigMapsPrometheusService(cm *corev1.ConfigMap, prometheusservice string) {
	if strings.Contains(cm.Data[util.OriginComfigMapPrometheusLocation], util.OriginPrometheus_URL) && prometheusservice != "" {
		cm.Data[util.OriginComfigMapPrometheusLocation] = strings.Replace(cm.Data[util.OriginComfigMapPrometheusLocation], util.OriginPrometheus_URL, prometheusservice, -1)
	}
}
func GlobalSectionSetParamterToConfigMap(cm *corev1.ConfigMap, prometheusService string, namespace string) {
	processConfigMapsPrometheusService(cm, prometheusService) //ConfigMapData's PrometheusService
}

func processDeploymentPrometheusService(dep *appsv1.Deployment, prometheusservice string) {
	if flag, envIndex, ctnIndex := isPrometheusServiceDep(dep); flag == true && prometheusservice != "" {
		dep.Spec.Template.Spec.Containers[ctnIndex].Env[envIndex].Value = prometheusservice
	}
}

func isPrometheusServiceDep(dep *appsv1.Deployment) (bool, int, int) {
	for ctnIndex, v := range dep.Spec.Template.Spec.Containers {
		if len(v.Env) > 0 {
			for envIndex, value := range dep.Spec.Template.Spec.Containers[ctnIndex].Env {
				if value.Name == util.OriginDeploymentPrometheusLocation {
					return true, envIndex, ctnIndex
				}
			}
			return false, -1, -1
		}
	}
	return false, -1, -1
}

func GlobalSectionSetParamterToPersistentVolumeClaim(pvc *corev1.PersistentVolumeClaim, asp *alamedaserviceparamter.AlamedaServiceParamter) {
	for _, pvcusage := range v1alpha1.PvcUsage {
		if strings.Contains(pvc.Name, string(pvcusage)) {
			util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.Storages, pvcusage)
		}
	}
}

func getEnvVarsToUpdateByDeployment(deploymentName string, asp *alamedaserviceparamter.AlamedaServiceParamter) []corev1.EnvVar {

	var envVars []corev1.EnvVar

	switch deploymentName {
	case util.AdmissioncontrollerDPN:
		envVars = getAdmissionControllerEnvVarsToUpdate(asp)
	case util.AlamedaevictionerDPN:
		envVars = getAlamedaEvictionerEnvVarsToUpdate(asp)
	case util.AlamedaaiDPN:
		envVars = getAlamedaAIEnvVarsToUpdate(asp)
	default:
	}

	return envVars
}

func getAlamedaAIEnvVarsToUpdate(asp *alamedaserviceparamter.AlamedaServiceParamter) []corev1.EnvVar {

	envVars := make([]corev1.EnvVar, 0)

	switch asp.EnableDispatcher {
	case true:
		envVars = append(envVars, corev1.EnvVar{
			Name:  "PREDICT_QUEUE_ENABLED",
			Value: "true",
		})
	case false:
		envVars = append(envVars, corev1.EnvVar{
			Name:  "PREDICT_QUEUE_ENABLED",
			Value: "false",
		})
	}

	return envVars
}

func getAdmissionControllerEnvVarsToUpdate(asp *alamedaserviceparamter.AlamedaServiceParamter) []corev1.EnvVar {

	envVars := make([]corev1.EnvVar, 0)

	switch asp.Platform {
	case v1alpha1.PlatformOpenshift3_9:
		envVars = append(envVars, corev1.EnvVar{
			Name:  "ALAMEDA_ADMCTL_JSONPATCHVALIDATIONFUNC",
			Value: admission_controller.JsonPatchValidationFuncOpenshift3_9,
		})
	default:
		envVars = append(envVars, corev1.EnvVar{
			Name:  "ALAMEDA_ADMCTL_JSONPATCHVALIDATIONFUNC",
			Value: "",
		})
	}

	return envVars
}

func getAlamedaEvictionerEnvVarsToUpdate(asp *alamedaserviceparamter.AlamedaServiceParamter) []corev1.EnvVar {
	envVars := make([]corev1.EnvVar, 0)

	switch asp.Platform {
	case v1alpha1.PlatformOpenshift3_9:
		envVars = append(envVars, corev1.EnvVar{
			Name:  "ALAMEDA_EVICTIONER_EVICTION_PURGECONTAINERCPUMEMORY",
			Value: "true",
		})
	default:
		envVars = append(envVars, corev1.EnvVar{
			Name:  "ALAMEDA_EVICTIONER_EVICTION_PURGECONTAINERCPUMEMORY",
			Value: "false",
		})
	}

	return envVars
}
