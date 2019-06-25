package recommendation

type controllerTag = string
type controllerField = string

const (
	ControllerTime      controllerTag = "time"
	ControllerNamespace controllerTag = "namespace"
	ControllerName      controllerTag = "name"
	ControllerType      controllerTag = "type"

	ControllerKind              controllerField = "kind"
	ControllerCurrentReplicas   controllerField = "current_replicas"
	ControllerDesiredReplicas   controllerField = "desired_replicas"
	ControllerCreateTime        controllerField = "create_time"
	ControllerCurrentCPURequest controllerField = "current_cpu_requests"
	ControllerCurrentMEMRequest controllerField = "current_mem_requests"
	ControllerCurrentCPULimit   controllerField = "current_cpu_limits"
	ControllerCurrentMEMLimit   controllerField = "current_mem_limits"
	ControllerDesiredCPULimit   controllerField = "desired_cpu_limits"
	ControllerDesiredMEMLimit   controllerField = "desired_mem_limits"
)

var (
	// ControllerTags is list of tags of alameda_controller_recommendation measurement
	ControllerTags = []controllerTag{
		ControllerTime,
		ControllerNamespace,
		ControllerName,
		ControllerType,
	}
	// ControllerFields is list of fields of alameda_controller_recommendation measurement
	ControllerField = []controllerField{
		ControllerCurrentReplicas,
		ControllerDesiredReplicas,
		ControllerCreateTime,
		ControllerKind,

		ControllerCurrentCPURequest,
		ControllerCurrentMEMRequest,
		ControllerCurrentCPULimit,
		ControllerCurrentMEMLimit,
		ControllerDesiredCPULimit,
		ControllerDesiredMEMLimit,
	}
)
