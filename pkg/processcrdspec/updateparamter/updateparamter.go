package updateparamter

import (
	"strings"

	"github.com/containers-ai/federatorai-operator/pkg/apis/federatorai/v1alpha1"
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

func ProcessConfigMapsPrometheusService(cm *corev1.ConfigMap, prometheusservice string) *corev1.ConfigMap {
	log.V(1).Info("ProcessPrometheusService", "AlamedaServicePrometheusService", prometheusservice)
	if strings.Contains(cm.Data["prometheus.yaml"], "https://prometheus-k8s.openshift-monitoring.svc:9091") && prometheusservice != "" {
		cm.Data["prometheus.yaml"] = strings.Replace(cm.Data["prometheus.yaml"], "https://prometheus-k8s.openshift-monitoring.svc:9091", prometheusservice, -1)
	}
	return cm
}
func ProcessDeploymentsPrometheusService(dep *appsv1.Deployment, prometheusservice string) *appsv1.Deployment {
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

/*
func isGrafanaPVC(dep *appsv1.Deployment) (bool, string) {
	if len(dep.Spec.Template.Spec.Volumes) > 0 {
		if dep.Spec.Template.Spec.Volumes[0].Name == "grafana-storage" {
			return true, "my-alameda.grafana.pvc"
		}
	}
	return false, ""
}
func isInfluxDBPVC(dep *appsv1.Deployment) (bool, string) {
	if len(dep.Spec.Template.Spec.Volumes) > 0 {
		if dep.Spec.Template.Spec.Volumes[0].Name == "influxdb-storage" {
			return true, "my-alameda.influxdb.pvc"
		}
	}
	return false, ""
}

func ProcessDepPVC(dep *appsv1.Deployment, influxdbFlag, grafanaFlag bool) *appsv1.Deployment {
	if flag, claimname := isGrafanaPVC(dep); flag == true && claimname != "" && grafanaFlag {
		pvcs := &corev1.PersistentVolumeClaimVolumeSource{ClaimName: claimname}
		vs := corev1.VolumeSource{PersistentVolumeClaim: pvcs}
		dep.Spec.Template.Spec.Volumes[0].VolumeSource = vs
	}
	if flag, claimname := isInfluxDBPVC(dep); flag == true && claimname != "" && influxdbFlag {
		pvcs := &corev1.PersistentVolumeClaimVolumeSource{ClaimName: claimname}
		vs := corev1.VolumeSource{PersistentVolumeClaim: pvcs}
		dep.Spec.Template.Spec.Volumes[0].VolumeSource = vs
	}
	return dep
}*/
func ProcessDepPVC(dep *appsv1.Deployment, index int, value interface{}) *appsv1.Deployment {
	switch v := value.(type) {
	case v1alpha1.AlamedaAILog:
		{
			if v.Flag {
				pvcs := &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "alameda-ai.pvc"}
				vs := corev1.VolumeSource{PersistentVolumeClaim: pvcs}
				dep.Spec.Template.Spec.Volumes[index].VolumeSource = vs
			}
		}
	case v1alpha1.AlamedaDatahubLog:
		{
			if v.Flag {
				pvcs := &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "alameda-datahub.pvc"}
				vs := corev1.VolumeSource{PersistentVolumeClaim: pvcs}
				dep.Spec.Template.Spec.Volumes[index].VolumeSource = vs
			}
		}
	case v1alpha1.AlamedaOperatorLog:
		{
			if v.Flag {
				pvcs := &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "alameda-operator.pvc"}
				vs := corev1.VolumeSource{PersistentVolumeClaim: pvcs}
				dep.Spec.Template.Spec.Volumes[index].VolumeSource = vs
			}

		}
	case v1alpha1.AlamedaEvictionerLog:
		{
			if v.Flag {
				pvcs := &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "alameda-evictioner.pvc"}
				vs := corev1.VolumeSource{PersistentVolumeClaim: pvcs}
				dep.Spec.Template.Spec.Volumes[index].VolumeSource = vs
			}

		}
	case v1alpha1.AdmissionControllerLog:
		{
			if v.Flag {
				pvcs := &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "admission-controller.pvc"}
				vs := corev1.VolumeSource{PersistentVolumeClaim: pvcs}
				dep.Spec.Template.Spec.Volumes[index].VolumeSource = vs
			}

		}
	case v1alpha1.AlamedaServiceSpecInfluxdbPVCSet:
		{
			if v.Flag {
				pvcs := &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "my-alameda.influxdb.pvc"}
				vs := corev1.VolumeSource{PersistentVolumeClaim: pvcs}
				dep.Spec.Template.Spec.Volumes[index].VolumeSource = vs
			}

		}
	case v1alpha1.AlamedaServiceSpecGrafanaPVCSet:
		{
			if v.Flag {
				pvcs := &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "my-alameda.grafana.pvc"}
				vs := corev1.VolumeSource{PersistentVolumeClaim: pvcs}
				dep.Spec.Template.Spec.Volumes[index].VolumeSource = vs
			}

		}
	}
	return dep
}

func ProcessComponentLogPVC(pvc *corev1.PersistentVolumeClaim, value interface{}) *corev1.PersistentVolumeClaim {
	switch v := value.(type) {
	case v1alpha1.AlamedaAILog:
		{
			pvc.Spec = v.Spec
			return pvc
		}
	case v1alpha1.AlamedaDatahubLog:
		{
			pvc.Spec = v.Spec
			return pvc
		}
	case v1alpha1.AlamedaOperatorLog:
		{
			pvc.Spec = v.Spec
			return pvc
		}
	case v1alpha1.AlamedaEvictionerLog:
		{
			pvc.Spec = v.Spec
			return pvc
		}
	case v1alpha1.AdmissionControllerLog:
		{
			pvc.Spec = v.Spec
			return pvc
		}
	case v1alpha1.AlamedaServiceSpecInfluxdbPVCSet:
		{
			pvc.Spec = v.Spec
			return pvc
		}
	case v1alpha1.AlamedaServiceSpecGrafanaPVCSet:
		{
			pvc.Spec = v.Spec
			return pvc
		}
	}
	return pvc
}
