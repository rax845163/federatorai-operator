package util

import (
	corev1 "k8s.io/api/core/v1"
)

func IsEmpty(pvcspec corev1.PersistentVolumeClaimSpec) bool {
	if len(pvcspec.AccessModes) > 0 && pvcspec.Resources.Requests != nil {
		return false
	}
	return true
}
