package clusterstatus

type controllerTag = string
type controllerField = string

const (
	ControllerTime           controllerTag = "time"
	ControllerNamespace      controllerTag = "namespace"
	ControllerName           controllerTag = "name"
	ControllerOwnerNamespace controllerTag = "owner_namespace"
	ControllerOwnerName      controllerTag = "owner_name"

	ControllerKind            controllerField = "kind"
	ControllerOwnerKind       controllerField = "owner_kind"
	ControllerReplicas        controllerField = "replicas"
	ControllerEnableExecution controllerField = "enable_execution"
	ControllerPolicy          controllerField = "policy"
)

var (
	// ControllerTags is list of tags of alameda_controller_recommendation measurement
	ControllerTags = []controllerTag{
		ControllerTime,
		ControllerNamespace,
		ControllerName,
		ControllerOwnerNamespace,
		ControllerOwnerName,
	}
	// ControllerFields is list of fields of alameda_controller_recommendation measurement
	ControllerFields = []controllerField{
		ControllerKind,
		ControllerOwnerKind,
		ControllerReplicas,
		ControllerEnableExecution,
		ControllerPolicy,
	}
)
