package container

import (
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	"github.com/golang/protobuf/ptypes/timestamp"
	corev1 "k8s.io/api/core/v1"
)

func NewStatus(containerStatus *corev1.ContainerStatus) *datahub_v1alpha1.ContainerStatus {
	state := &datahub_v1alpha1.ContainerState{}
	if containerStatus.State.Running != nil {
		state.Running = &datahub_v1alpha1.ContainerStateRunning{
			StartedAt: &timestamp.Timestamp{
				Seconds: containerStatus.State.Running.StartedAt.Unix(),
			},
		}
	} else if containerStatus.State.Terminated != nil {
		state.Terminated = &datahub_v1alpha1.ContainerStateTerminated{
			ExitCode: containerStatus.State.Terminated.ExitCode,
			Reason:   containerStatus.State.Terminated.Reason,
			Message:  containerStatus.State.Terminated.Message,
			StartedAt: &timestamp.Timestamp{
				Seconds: containerStatus.State.Terminated.StartedAt.Unix(),
			},
			FinishedAt: &timestamp.Timestamp{
				Seconds: containerStatus.State.Terminated.FinishedAt.Unix(),
			},
		}
	} else if containerStatus.State.Waiting != nil {
		state.Waiting = &datahub_v1alpha1.ContainerStateWaiting{
			Reason:  containerStatus.State.Waiting.Reason,
			Message: containerStatus.State.Waiting.Message,
		}
	}
	lastTerminationState := &datahub_v1alpha1.ContainerState{}
	if containerStatus.LastTerminationState.Running != nil {
		lastTerminationState.Running = &datahub_v1alpha1.ContainerStateRunning{
			StartedAt: &timestamp.Timestamp{
				Seconds: containerStatus.LastTerminationState.Running.StartedAt.Unix(),
			},
		}
	} else if containerStatus.LastTerminationState.Terminated != nil {
		lastTerminationState.Terminated = &datahub_v1alpha1.ContainerStateTerminated{
			ExitCode: containerStatus.LastTerminationState.Terminated.ExitCode,
			Reason:   containerStatus.LastTerminationState.Terminated.Reason,
			Message:  containerStatus.LastTerminationState.Terminated.Message,
			StartedAt: &timestamp.Timestamp{
				Seconds: containerStatus.LastTerminationState.Terminated.StartedAt.Unix(),
			},
			FinishedAt: &timestamp.Timestamp{
				Seconds: containerStatus.LastTerminationState.Terminated.FinishedAt.Unix(),
			},
		}
	} else if containerStatus.LastTerminationState.Waiting != nil {
		lastTerminationState.Waiting = &datahub_v1alpha1.ContainerStateWaiting{
			Reason:  containerStatus.LastTerminationState.Waiting.Reason,
			Message: containerStatus.LastTerminationState.Waiting.Message,
		}
	}
	return &datahub_v1alpha1.ContainerStatus{
		RestartCount:         containerStatus.RestartCount,
		State:                state,
		LastTerminationState: lastTerminationState,
	}
}
