package updateenvvar

import (
	"strings"

	"github.com/containers-ai/federatorai-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func AssignServiceToDeployment(dep *appsv1.Deployment, ns string) {
	if len(dep.Spec.Template.Spec.Containers[0].Env) > 0 {
		for index, value := range dep.Spec.Template.Spec.Containers[0].Env {
			if strings.Contains(value.String(), util.NamespaceService) {
				dep.Spec.Template.Spec.Containers[0].Env[index].Value = strings.Replace(dep.Spec.Template.Spec.Containers[0].Env[index].Value, util.NamespaceService, ns+".svc", -1)
			}
		}
	}
}
func AssignServiceToConfigMap(cm *corev1.ConfigMap, ns string) {
	if strings.Contains(cm.Data[util.OriginComfigMapPrometheusLocation], util.NamespaceService) {
		cm.Data[util.OriginComfigMapPrometheusLocation] = strings.Replace(cm.Data[util.OriginComfigMapPrometheusLocation], util.NamespaceService, ns+".svc", -1)
	}
}
