package updateparamter

import (
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var (
	log = logf.Log.WithName("controller_alamedaservice")
)

func isProphetstorImage(dep *appsv1.Deployment) bool {
	if strings.Contains(dep.Spec.Template.Spec.Containers[0].Image, "quay.io/prophetstor/") {
		return true
	} else {
		return false
	}
}

func ProcessImageVersion(dep *appsv1.Deployment, version string) *appsv1.Deployment {
	log.V(1).Info("ProcessImageVersion", "AlamedaServiceImageVersion", version)
	if isProphetstorImage(dep) {
		s := strings.Split(dep.Spec.Template.Spec.Containers[0].Image, ":")
		if len(s) != 0 && version != "" {
			dep.Spec.Template.Spec.Containers[0].Image = s[0] + ":" + version
			log.V(1).Info("ProcessImageVersion", "AlamedaComponentDep.Containers.Image", dep.Spec.Template.Spec.Containers[0].Image)
		}
	}
	return dep
}

func getImageVersion(dep *appsv1.Deployment) string {
	var version = ""
	if isProphetstorImage(dep) {
		s := strings.Split(dep.Spec.Template.Spec.Containers[0].Image, ":")
		if len(s) != 0 {
			version = s[1]
		}
	}
	return version
}

func MisMatchImageVersion(dep *appsv1.Deployment, version string) bool {
	if (getImageVersion(dep) != version && version != "") && isProphetstorImage(dep) {
		return true
	} else {
		return false
	}
}

func isPrometheusService(dep *appsv1.Deployment) (bool, int) {
	if len(dep.Spec.Template.Spec.Containers[0].Env) > 0 {
		for index, value := range dep.Spec.Template.Spec.Containers[0].Env {
			if value.Name == "ALAMEDA_DATAHUB_PROMETHEUS_URL" {
				return true, index
			}
		}
		return false, 0
	}
	return false, 0
}

func ProcessPrometheusService(dep *appsv1.Deployment, prometheusservice string) *appsv1.Deployment {
	log.V(1).Info("ProcessPrometheusService", "AlamedaServicePrometheusService", prometheusservice)
	if flag, index := isPrometheusService(dep); flag == true && prometheusservice != "" {
		dep.Spec.Template.Spec.Containers[0].Env[index].Value = prometheusservice
		log.V(1).Info("ProcessPrometheusService", "DatahubDep.Env.ALAMEDA_DATAHUB_PROMETHEUS_URL", dep.Spec.Template.Spec.Containers[0].Env[index].Value)
	}
	return dep
}

func getPrometheusService(dep *appsv1.Deployment) (string, bool) {
	var flag = false
	var prometheusservice = ""
	if Matchflag, index := isPrometheusService(dep); Matchflag == true {
		prometheusservice = dep.Spec.Template.Spec.Containers[0].Env[index].Value
		return prometheusservice, Matchflag
	}
	return prometheusservice, flag
}

func MisMatchPrometheusService(dep *appsv1.Deployment, prometheusservice string) bool {
	if clusterPrometheusService, flag := getPrometheusService(dep); (clusterPrometheusService != prometheusservice && prometheusservice != "") && flag {
		return true
	}
	return false
}

func MisMatchAlamedaServiceParamter(dep *appsv1.Deployment, version string, prometheusservice string) bool {
	misMatchIV := MisMatchImageVersion(dep, version)
	micMatchPS := MisMatchPrometheusService(dep, prometheusservice)

	log.V(1).Info("MisMatchAlamedaServiceParamter", "ComponentDeployment", dep.Name, "MatchImageVersion", misMatchIV, "MatchPrometheusService", micMatchPS)

	if misMatchIV || micMatchPS {
		return true
	}
	return false
}

func MisMatchPVC(dep *appsv1.Deployment, claimName string) bool {
	if flag := isPVC(dep); flag == true && dep.Spec.Template.Spec.Volumes[0].VolumeSource.PersistentVolumeClaim != nil {
		if dep.Spec.Template.Spec.Volumes[0].VolumeSource.PersistentVolumeClaim.ClaimName != claimName {
			return true
		}
	}
	return false
}

func isPVC(dep *appsv1.Deployment) bool {
	if len(dep.Spec.Template.Spec.Volumes) > 0 {
		if dep.Spec.Template.Spec.Volumes[0].Name == "grafana-storage" {
			return true
		}
	}
	return false
}

func ProcessPVC(dep *appsv1.Deployment, claimname string) *appsv1.Deployment {
	if flag := isPVC(dep); flag == true && claimname != "" {
		pvcs := &corev1.PersistentVolumeClaimVolumeSource{ClaimName: claimname}
		vs := corev1.VolumeSource{PersistentVolumeClaim: pvcs}
		dep.Spec.Template.Spec.Volumes[0].VolumeSource = vs
	}
	return dep
}

/*
func isPrometheusService(dep *appsv1.Deployment) (bool, int) {
	if len(dep.Spec.Template.Spec.Containers[0].Env) > 0 {
		for index, value := range dep.Spec.Template.Spec.Containers[0].Env {
			if value.Name == "ALAMEDA_DATAHUB_PROMETHEUS_URL" {
				return true, index
			}
		}
		return false, 0
	}
	return false, 0
}
*/
