package pod

import (
	utilsresource "github.com/containers-ai/alameda/operator/pkg/utils/resources"
	"github.com/containers-ai/alameda/pkg/consts"
	"github.com/containers-ai/alameda/pkg/utils/datahub/enumconv"
	logUtil "github.com/containers-ai/alameda/pkg/utils/log"
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	scope = logUtil.RegisterScope("datahubpodutils", "datahub pod utils", 0)
)

// NewStatus return pod status struct of datahub
func NewStatus(pod *corev1.Pod) *datahub_v1alpha1.PodStatus {
	return &datahub_v1alpha1.PodStatus{
		Message: pod.Status.Message,
		Reason:  pod.Status.Reason,
		Phase:   enumconv.PodPhaseEnumK8SToDatahub[pod.Status.Phase],
	}
}

// GetReplicasFromPod return number of replicas of pod
func GetReplicasFromPod(pod *corev1.Pod, client client.Client) int32 {
	getResource := utilsresource.NewGetResource(client)

	for _, or := range pod.OwnerReferences {
		if or.Kind == consts.K8S_KIND_REPLICASET {
			rs, err := getResource.GetReplicaSet(pod.GetNamespace(), or.Name)
			if err == nil {
				return rs.Status.Replicas
			} else {
				scope.Errorf("Get replicaset for number of replicas failed due to %s", err.Error())
			}
		} else if or.Kind == consts.K8S_KIND_REPLICATIONCONTROLLER {
			rc, err := getResource.GetReplicationController(pod.GetNamespace(), or.Name)
			if err == nil {
				return rc.Status.Replicas
			} else {
				scope.Errorf("Get replicationcontroller for number of replicas failed due to %s", err.Error())
			}
		}
	}
	return int32(1)
}
