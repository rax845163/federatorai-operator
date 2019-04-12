package processcrdspec

import (
	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/alamedaserviceparamter"
	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/updateenvvar"
	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/updateparamter"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func findindex(dep *appsv1.Deployment) int {
	if len(dep.Spec.Template.Spec.Volumes) > 0 {
		for index, value := range dep.Spec.Template.Spec.Volumes {
			if value.Name == "alameda-ai-log-storage" {
				return index
			}
			if value.Name == "alameda-operator-log-storage" {
				return index
			}
			if value.Name == "alameda-datahub-log-storage" {
				return index
			}
			if value.Name == "alameda-evictioner-log-storage" {
				return index
			}
			if value.Name == "admission-controller-log-storage" {
				return index
			}
			if value.Name == "influxdb-storage" {
				return index
			}
			if value.Name == "grafana-storage" {
				return index
			}
		}
		return -1
	}
	return -1
}
func ParamterToDeployment(dep *appsv1.Deployment, asp *alamedaserviceparamter.AlamedaServiceParamter) *appsv1.Deployment {
	dep = updateenvvar.AssignServiceToDeployment(dep, dep.Namespace) //DeploymentSpec's service
	dep = updateparamter.ProcessImageVersion(dep, asp.Version)
	dep = updateparamter.ProcessDeploymentsPrometheusService(dep, asp.PrometheusService)
	//dep = updateparamter.ProcessDepPVC(dep, asp.InfluxdbPVCSet.Flag, asp.GrafanaPVCSet.Flag) //if user set pvc
	if index := findindex(dep); index != -1 {
		switch dep.Name {
		case "alameda-ai":
			{
				dep = updateparamter.ProcessDepPVC(dep, index, asp.AlamedaAILog)
			}
		case "alameda-operator":
			{
				dep = updateparamter.ProcessDepPVC(dep, index, asp.AlamedaOperatorLog)
			}
		case "alameda-datahub":
			{
				dep = updateparamter.ProcessDepPVC(dep, index, asp.AlamedaDatahubLog)
			}
		case "alameda-evictioner":
			{
				dep = updateparamter.ProcessDepPVC(dep, index, asp.AlamedaEvictionerLog)
			}
		case "admission-controller":
			{
				dep = updateparamter.ProcessDepPVC(dep, index, asp.AdmissionControllerLog)
			}
		case "alameda-influxdb":
			{
				dep = updateparamter.ProcessDepPVC(dep, index, asp.InfluxdbPVCSet)
			}
		case "alameda-grafana":
			{
				dep = updateparamter.ProcessDepPVC(dep, index, asp.GrafanaPVCSet)
			}
		}
	}
	return dep
}

func ParamterToConfigMap(cm *corev1.ConfigMap, asp *alamedaserviceparamter.AlamedaServiceParamter) *corev1.ConfigMap {
	cm = updateenvvar.AssignServiceToConfigMap(cm, cm.Namespace) //ConfigMapSpec's service
	cm = updateparamter.ProcessConfigMapsPrometheusService(cm, asp.PrometheusService)
	return cm
}
func ParamterToPersistentVolumeClaim(pvc *corev1.PersistentVolumeClaim, asp *alamedaserviceparamter.AlamedaServiceParamter) *corev1.PersistentVolumeClaim {
	switch pvc.Name {
	case "alameda-ai.pvc":
		{
			pvc = updateparamter.ProcessComponentLogPVC(pvc, asp.AlamedaAILog)
		}
	case "alameda-operator.pvc":
		{
			pvc = updateparamter.ProcessComponentLogPVC(pvc, asp.AlamedaOperatorLog)
		}
	case "alameda-datahub.pvc":
		{
			pvc = updateparamter.ProcessComponentLogPVC(pvc, asp.AlamedaDatahubLog)
		}
	case "alameda-evictioner.pvc":
		{
			pvc = updateparamter.ProcessComponentLogPVC(pvc, asp.AlamedaEvictionerLog)
		}
	case "admission-controller.pvc":
		{
			pvc = updateparamter.ProcessComponentLogPVC(pvc, asp.AdmissionControllerLog)
		}
	case "my-alameda.influxdb.pvc":
		{
			pvc = updateparamter.ProcessComponentLogPVC(pvc, asp.InfluxdbPVCSet)
		}
	case "my-alameda.grafana.pvc":
		{
			pvc = updateparamter.ProcessComponentLogPVC(pvc, asp.GrafanaPVCSet)
		}
	}
	return pvc
}
