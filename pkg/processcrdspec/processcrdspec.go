package processcrdspec

import (
	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/alamedaserviceparamter"
	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/componentsectionset"
	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/globalsectionset"
	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/updateenvvar"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func ParamterToDeployment(dep *appsv1.Deployment, asp *alamedaserviceparamter.AlamedaServiceParamter) *appsv1.Deployment {
	updateenvvar.AssignServiceToDeployment(dep, dep.Namespace)      //DeploymentSpec's service
	globalsectionset.GlobalSectionSetParamterToDeployment(dep, asp) //DeploymentSpec's Global Section Set
	componentsectionset.SectionSetParamterToDeployment(dep, asp)    //DeploymentSpec's Component Section Set
	return dep
}
func ParamterToConfigMap(cm *corev1.ConfigMap, asp *alamedaserviceparamter.AlamedaServiceParamter) *corev1.ConfigMap {
	updateenvvar.AssignServiceToConfigMap(cm, cm.Namespace)                         //ConfigMapSpec's service
	globalsectionset.GlobalSectionSetParamterToConfigMap(cm, asp.PrometheusService) //ConfigMapSpec's PrometheusService
	return cm
}
func ParamterToPersistentVolumeClaim(pvc *corev1.PersistentVolumeClaim, asp *alamedaserviceparamter.AlamedaServiceParamter) *corev1.PersistentVolumeClaim {
	globalsectionset.GlobalSectionSetParamterToPersistentVolumeClaim(pvc, asp)
	componentsectionset.SectionSetParamterToPersistentVolumeClaim(pvc, asp) //PersistentVolumeClaim's Component Section Set
	return pvc
}
