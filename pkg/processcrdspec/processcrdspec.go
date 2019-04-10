package processcrdspec

import (
	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/alamedaserviceparamter"
	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/updateenvvar"
	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/updateparamter"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func ParamterToDeployment(dep *appsv1.Deployment, asp *alamedaserviceparamter.AlamedaServiceParamter) *appsv1.Deployment {
	dep = updateenvvar.AssignServiceToDeployment(dep, dep.Namespace) //DeploymentSpec's service
	dep = updateparamter.ProcessImageVersion(dep, asp.Version)
	dep = updateparamter.ProcessDeploymentsPrometheusService(dep, asp.PrometheusService)
	dep = updateparamter.ProcessDepPVC(dep, asp.InfluxdbPVCSet.Flag, asp.GrafanaPVCSet.Flag) //if user set pvc
	return dep
}

func ParamterToConfigMap(cm *corev1.ConfigMap, asp *alamedaserviceparamter.AlamedaServiceParamter) *corev1.ConfigMap {
	cm = updateenvvar.AssignServiceToConfigMap(cm, cm.Namespace) //ConfigMapSpec's service
	cm = updateparamter.ProcessConfigMapsPrometheusService(cm, asp.PrometheusService)
	return cm
}
func ParamterToPersistentVolumeClaim(pvc *corev1.PersistentVolumeClaim, asp *alamedaserviceparamter.AlamedaServiceParamter) *corev1.PersistentVolumeClaim {
	if pvc.Name == "my-alameda.grafana.pvc" {
		pvc = updateparamter.ProcessGrafanaPVC(pvc, asp.GrafanaPVCSet)
	}
	if pvc.Name == "my-alameda.influxdb.pvc" {
		pvc = updateparamter.ProcessInfluxDBPVC(pvc, asp.InfluxdbPVCSet)
	}
	return pvc
}
