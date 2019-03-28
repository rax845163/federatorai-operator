package updateenvvar

import (
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func AssignNameSpace(ns string) {

}
func AssignDeployment(dep *appsv1.Deployment, ns string) *appsv1.Deployment {
	if len(dep.Spec.Template.Spec.Containers[0].Env) > 0 {
		for index, value := range dep.Spec.Template.Spec.Containers[0].Env {
			if strings.Contains(value.String(), "federatorai.svc") {
				dep.Spec.Template.Spec.Containers[0].Env[index].Value = strings.Replace(dep.Spec.Template.Spec.Containers[0].Env[index].Value, "federatorai.svc", ns+".svc", -1)
			}
		}
	}
	return dep
}
func AssignConfigMap(cm *corev1.ConfigMap, ns string) *corev1.ConfigMap {
	if strings.Contains(cm.Data["prometheus.yaml"], "federatorai.svc") {
		cm.Data["prometheus.yaml"] = strings.Replace(cm.Data["prometheus.yaml"], "federatorai.svc", ns+".svc", -1)
	}
	return cm
}
