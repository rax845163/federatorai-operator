package enumconv

import (
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	corev1 "k8s.io/api/core/v1"
)

var PodPhaseEnumDatahubToK8S map[datahub_v1alpha1.PodPhase]corev1.PodPhase = map[datahub_v1alpha1.PodPhase]corev1.PodPhase{
	datahub_v1alpha1.PodPhase_Pending:   corev1.PodPending,
	datahub_v1alpha1.PodPhase_Running:   corev1.PodRunning,
	datahub_v1alpha1.PodPhase_Succeeded: corev1.PodSucceeded,
	datahub_v1alpha1.PodPhase_Failed:    corev1.PodFailed,
	datahub_v1alpha1.PodPhase_Unknown:   corev1.PodUnknown,
}

var PodPhaseEnumK8SToDatahub map[corev1.PodPhase]datahub_v1alpha1.PodPhase = map[corev1.PodPhase]datahub_v1alpha1.PodPhase{
	corev1.PodPending:   datahub_v1alpha1.PodPhase_Pending,
	corev1.PodRunning:   datahub_v1alpha1.PodPhase_Running,
	corev1.PodSucceeded: datahub_v1alpha1.PodPhase_Succeeded,
	corev1.PodFailed:    datahub_v1alpha1.PodPhase_Failed,
	corev1.PodUnknown:   datahub_v1alpha1.PodPhase_Unknown,
}
