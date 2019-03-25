package imageversion

import (
	"strings"

	appsv1 "k8s.io/api/apps/v1"
)

func ProcessImageVersion(dep *appsv1.Deployment, version string) *appsv1.Deployment {
	if isProphetstorImage(dep) {
		s := strings.Split(dep.Spec.Template.Spec.Containers[0].Image, ":")
		if len(s) != 0 {
			dep.Spec.Template.Spec.Containers[0].Image = s[0] + ":" + version
		}
	}
	return dep
}
func isProphetstorImage(dep *appsv1.Deployment) bool {
	if strings.Contains(dep.Spec.Template.Spec.Containers[0].Image, "quay.io/prophetstor/") {
		return true
	} else {
		return false
	}
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
