package updateparamter

import (
	"strings"

	appsv1 "k8s.io/api/apps/v1"
)

func isProphetstorImage(dep *appsv1.Deployment) bool {
	if strings.Contains(dep.Spec.Template.Spec.Containers[0].Image, "quay.io/prophetstor/") {
		return true
	} else {
		return false
	}
}

func ProcessImageVersion(dep *appsv1.Deployment, version string) *appsv1.Deployment {
	if isProphetstorImage(dep) {
		s := strings.Split(dep.Spec.Template.Spec.Containers[0].Image, ":")
		if len(s) != 0 {
			dep.Spec.Template.Spec.Containers[0].Image = s[0] + ":" + version
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

func MatchImageVersion(dep *appsv1.Deployment, version string) bool {
	if getImageVersion(dep) != version && isProphetstorImage(dep) {
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
	if flag, index := isPrometheusService(dep); flag == true {
		dep.Spec.Template.Spec.Containers[0].Env[index].Value = prometheusservice
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

func MatchPrometheusService(dep *appsv1.Deployment, prometheusservice string) bool {
	if clusterprometheusservice, flag := getPrometheusService(dep); clusterprometheusservice != prometheusservice && flag {
		return true
	} else {
		return false
	}
}

func MatchAlamedaServiceParamter(dep *appsv1.Deployment, version string, prometheusservice string) bool {
	if MatchImageVersion(dep, version) || MatchPrometheusService(dep, prometheusservice) {
		return true
	}
	return false
}
