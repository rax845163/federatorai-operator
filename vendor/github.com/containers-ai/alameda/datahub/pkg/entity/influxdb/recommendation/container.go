package recommendation

type containerTag = string
type containerField = string

const (
	// ContainerTime is the time to apply recommendation
	ContainerTime containerTag = "time"
	// ContainerNamespace is recommended container namespace
	ContainerNamespace containerTag = "namespace"
	// ContainerName is recommended container name
	ContainerName containerTag = "name"
	// ContainerPodName is pod name of recommended container
	ContainerPodName     containerTag = "pod_name"
	ContainerGranularity containerTag = "granularity"

	// ContainerPolicy is recommended policy
	ContainerPolicy     containerField = "policy"
	ContainerPolicyTime containerField = "policy_time"
	// ContainerResourceRequestCPU is recommended CPU request
	ContainerResourceRequestCPU containerField = "resource_request_cpu"
	// ContainerResourceRequestMemory is recommended memory request
	ContainerResourceRequestMemory containerField = "resource_request_memory"
	// ContainerResourceLimitCPU is recommended CPU limit
	ContainerResourceLimitCPU containerField = "resource_limit_cpu"
	// ContainerResourceLimitMemory is recommended memory limit
	ContainerResourceLimitMemory containerField = "resource_limit_memory"
	// ContainerInitialResourceRequestCPU is recommended initial CPU request
	ContainerInitialResourceRequestCPU containerField = "initial_resource_request_cpu"
	// ContainerInitialResourceRequestMemory is recommended initial memory request
	ContainerInitialResourceRequestMemory containerField = "initial_resource_request_memory"
	// ContainerInitialResourceLimitCPU is recommended initial CPU limit
	ContainerInitialResourceLimitCPU containerField = "initial_resource_limit_cpu"
	// ContainerInitialResourceLimitMemory is recommended initial memory limit
	ContainerInitialResourceLimitMemory containerField = "initial_resource_limit_memory"
	// ContainerStartTime is recommended start time
	ContainerStartTime containerField = "start_time"
	// ContainerEndTime is recommended end time
	ContainerEndTime containerField = "end_time"
	// ContainerTopControllerName is top controller name of the pod
	ContainerTopControllerName containerField = "top_controller_name"
	// ContainerTopControllerKind is top controller kind of the pod
	ContainerTopControllerKind containerField = "top_controller_kind"
)

const (
	ContainerMetricKindLimit       = "limit"
	ContainerMetricKindRequest     = "request"
	ContainerMetricKindInitLimit   = "initLimit"
	ContainerMetricKindInitRequest = "initRequest"
)

var (
	ContainerMetricKinds = []string{
		ContainerMetricKindLimit,
		ContainerMetricKindRequest,
		ContainerMetricKindInitLimit,
		ContainerMetricKindInitRequest,
	}
)

var (
	// ContainerTags is list of tags of alameda_container_recommendation measurement
	ContainerTags = []containerTag{
		ContainerTime,
		ContainerNamespace,
		ContainerName,
		ContainerPodName,
		ContainerGranularity,
	}
	// ContainerFields is list of fields of alameda_container_recommendation measurement
	ContainerFields = []containerField{
		ContainerPolicy,
		ContainerResourceRequestCPU,
		ContainerResourceRequestMemory,
		ContainerResourceLimitCPU,
		ContainerResourceLimitMemory,
		ContainerInitialResourceRequestCPU,
		ContainerInitialResourceRequestMemory,
		ContainerInitialResourceLimitCPU,
		ContainerInitialResourceLimitMemory,
		ContainerStartTime, ContainerEndTime,
		ContainerTopControllerName,
		ContainerTopControllerKind,
	}
)
