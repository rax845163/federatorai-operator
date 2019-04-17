package processcrdspec

import (
	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/alamedaserviceparamter"
	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/updateenvvar"
	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/updateparamter"
	"github.com/containers-ai/federatorai-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func findindex(dep *appsv1.Deployment) int { //find volumeMount path's locat
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
	if index := findindex(dep); index != -1 { //find deployment and get index
		switch dep.Name {
		case "alameda-ai":
			{
				if !util.IsEmpty(asp.AlamedaAILog) { // if user set PVCspec
					pvcs := &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "alameda-ai.pvc"}
					vs := corev1.VolumeSource{PersistentVolumeClaim: pvcs}
					dep.Spec.Template.Spec.Volumes[index].VolumeSource = vs
				}
			}
		case "alameda-operator":
			{
				if !util.IsEmpty(asp.AlamedaOperatorLog) {
					pvcs := &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "alameda-operator.pvc"}
					vs := corev1.VolumeSource{PersistentVolumeClaim: pvcs}
					dep.Spec.Template.Spec.Volumes[index].VolumeSource = vs
				}
			}
		case "alameda-datahub":
			{
				if !util.IsEmpty(asp.AlamedaDatahubLog) {
					pvcs := &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "alameda-datahub.pvc"}
					vs := corev1.VolumeSource{PersistentVolumeClaim: pvcs}
					dep.Spec.Template.Spec.Volumes[index].VolumeSource = vs
				}
			}
		case "alameda-evictioner":
			{
				if !util.IsEmpty(asp.AlamedaEvictionerLog) {
					pvcs := &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "alameda-evictioner.pvc"}
					vs := corev1.VolumeSource{PersistentVolumeClaim: pvcs}
					dep.Spec.Template.Spec.Volumes[index].VolumeSource = vs
				}
			}
		case "admission-controller":
			{
				if !util.IsEmpty(asp.AdmissionControllerLog) {
					pvcs := &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "admission-controller.pvc"}
					vs := corev1.VolumeSource{PersistentVolumeClaim: pvcs}
					dep.Spec.Template.Spec.Volumes[index].VolumeSource = vs
				}
			}
		case "alameda-influxdb":
			{
				if !util.IsEmpty(asp.InfluxdbPVCSet) {
					pvcs := &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "my-alameda.influxdb.pvc"}
					vs := corev1.VolumeSource{PersistentVolumeClaim: pvcs}
					dep.Spec.Template.Spec.Volumes[index].VolumeSource = vs
				}
			}
		case "alameda-grafana":
			{
				if !util.IsEmpty(asp.GrafanaPVCSet) {
					pvcs := &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "my-alameda.grafana.pvc"}
					vs := corev1.VolumeSource{PersistentVolumeClaim: pvcs}
					dep.Spec.Template.Spec.Volumes[index].VolumeSource = vs
				}
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
			pvc.Spec = asp.AlamedaAILog
			return pvc
		}
	case "alameda-operator.pvc":
		{
			pvc.Spec = asp.AlamedaOperatorLog
		}
	case "alameda-datahub.pvc":
		{
			pvc.Spec = asp.AlamedaDatahubLog
		}
	case "alameda-evictioner.pvc":
		{
			pvc.Spec = asp.AlamedaEvictionerLog
		}
	case "admission-controller.pvc":
		{
			pvc.Spec = asp.AdmissionControllerLog
		}
	case "my-alameda.influxdb.pvc":
		{
			pvc.Spec = asp.InfluxdbPVCSet
		}
	case "my-alameda.grafana.pvc":
		{
			pvc.Spec = asp.GrafanaPVCSet
		}
	}
	return pvc
}
