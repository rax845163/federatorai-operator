package clusterstatus

import (
	"strconv"
	"time"

	"github.com/containers-ai/alameda/datahub/pkg/utils"
	influxdb_client "github.com/influxdata/influxdb/client/v2"
)

type containerTag = string
type containerField = string

const (
	// ContainerTime is the time that container information is saved to the measurement
	ContainerTime containerTag = "time"
	// ContainerNamespace is the container namespace
	ContainerNamespace containerTag = "namespace"
	// ContainerPodName is the name of pod that container is running in
	ContainerPodName containerTag = "pod_name"
	// ContainerAlamedaScalerNamespace is the namespace of AlamedaScaler that container belongs to
	ContainerAlamedaScalerNamespace containerTag = "alameda_scaler_namespace"
	// ContainerAlamedaScalerName is the name of AlamedaScaler that container belongs to
	ContainerAlamedaScalerName containerTag = "alameda_scaler_name"
	// ContainerNodeName is the name of node that container is running in
	ContainerNodeName containerTag = "node_name"
	// ContainerName is the container name
	ContainerName      containerTag = "name"
	ContainerAppName   containerTag = "app_name"
	ContainerAppPartOf containerTag = "app_part_of"

	// ContainerPodPhase is a label for the condition of a pod at the current time
	ContainerPodPhase containerField = "pod_phase"
	// ContainerPodMessage is a human readable message indicating details about why the pod is in this condition
	ContainerPodMessage containerField = "pod_message"
	// ContainerPodReason is a brief CamelCase message indicating details about why the pod is in this state
	ContainerPodReason containerField = "pod_reason"
	// ContainerStatusWaitingReason is a brief reason the container is not yet running
	ContainerStatusWaitingReason containerField = "status_waiting_reason"
	// ContainerStatusWaitingMessage is a message regarding why the container is not yet running
	ContainerStatusWaitingMessage containerField = "status_waiting_message"
	// ContainerStatusRunningStartedAt is a time at which the container was last (re-)started
	ContainerStatusRunningStartedAt containerField = "status_running_start_at"
	// ContainerStatusTerminatedExitCode is a exit status from the last termination of the container
	ContainerStatusTerminatedExitCode containerField = "status_terminated_exit_code"
	// ContainerStatusTerminatedReason is a brief reason from the last termination of the container
	ContainerStatusTerminatedReason containerField = "status_terminated_reason"
	// ContainerStatusTerminatedMessage is a message regarding the last termination of the container
	ContainerStatusTerminatedMessage containerField = "status_terminated_message"
	// ContainerStatusTerminatedStartedAt is a time at which previous execution of the container started
	ContainerStatusTerminatedStartedAt containerField = "status_terminated_started_at"
	// ContainerStatusTerminatedFinishedAt is a time at which the container last terminated
	ContainerStatusTerminatedFinishedAt containerField = "status_terminated_finished_at"
	// ContainerLastTerminationStatusWaitingReason is a last termination brief reason the container is not yet running
	ContainerLastTerminationStatusWaitingReason containerField = "last_termination_status_waiting_reason"
	// ContainerLastTerminationStatusWaitingMessage is a last termination message regarding why the container is not yet running
	ContainerLastTerminationStatusWaitingMessage containerField = "last_termination_status_waiting_message"
	// ContainerLastTerminationStatusRunningStartedAt is a last termination time at which the container was last (re-)started
	ContainerLastTerminationStatusRunningStartedAt containerField = "last_termination_status_running_start_at"
	// ContainerLastTerminationStatusTerminatedExitCode is a last termination exit status from the last termination of the container
	ContainerLastTerminationStatusTerminatedExitCode containerField = "last_termination_status_terminated_exit_code"
	// ContainerLastTerminationStatusTerminatedReason is a last termination brief reason from the last termination of the container
	ContainerLastTerminationStatusTerminatedReason containerField = "last_termination_status_terminated_reason"
	// ContainerLastTerminationStatusTerminatedMessage is a last termination message regarding the last termination of the container
	ContainerLastTerminationStatusTerminatedMessage containerField = "last_termination_status_terminated_message"
	// ContainerLastTerminationStatusTerminatedStartedAt is a last termination time at which previous execution of the container started
	ContainerLastTerminationStatusTerminatedStartedAt containerField = "last_termination_status_terminated_started_at"
	// ContainerLastTerminationStatusTerminatedFinishedAt is a last termination time at which the container last terminated
	ContainerLastTerminationStatusTerminatedFinishedAt containerField = "last_termination_status_terminated_finished_at"
	// ContainerRestartCount is the number of times the container has been restarted
	ContainerRestartCount containerField = "restart_count"
	// ContainerResourceRequestCPU is CPU request of the container
	ContainerResourceRequestCPU containerField = "resource_request_cpu"
	// ContainerResourceRequestMemory is memory request of the container
	ContainerResourceRequestMemory containerField = "resource_request_memroy"
	// ContainerResourceLimitCPU is CPU limit of the container
	ContainerResourceLimitCPU containerField = "resource_limit_cpu"
	// ContainerResourceLimitMemory is memory limit of the container
	ContainerResourceLimitMemory containerField = "resource_limit_memory"
	// ContainerPolicy is the prediction policy of container
	ContainerPolicy containerField = "policy"
	// ContainerPodCreateTime is the creation time of pod
	ContainerPodCreateTime containerField = "pod_create_time"
	// ContainerResourceLink is the resource link of pod
	ContainerResourceLink containerField = "resource_link"
	// ContainerTopControllerName is top controller name of the pod
	ContainerTopControllerName containerField = "top_controller_name"
	// ContainerTopControllerKind is top controller kind of the pod
	ContainerTopControllerKind containerField = "top_controller_kind"
	// ContainerTpoControllerReplicas is the number of replicas of container
	ContainerTpoControllerReplicas containerField = "top_controller_replicas"
	// ContainerUsedRecommendationID is the recommendation id that the pod applied
	ContainerUsedRecommendationID containerField = "used_recommendation_id"
	ContainerEnableVPA            containerField = "enable_VPA"
	ContainerEnableHPA            containerField = "enable_HPA"
)

var (
	// ContainerTags is the list of container measurement tags
	ContainerTags = []containerTag{
		ContainerTime, ContainerNamespace, ContainerPodName,
		ContainerAlamedaScalerNamespace, ContainerAlamedaScalerName,
		ContainerNodeName, ContainerName, ContainerAppName, ContainerAppPartOf,
	}
	// ContainerFields is the list of container measurement fields
	ContainerFields = []containerField{
		ContainerPodPhase, ContainerPodMessage, ContainerPodReason,
		ContainerStatusWaitingReason, ContainerStatusWaitingMessage,
		ContainerStatusRunningStartedAt,
		ContainerStatusTerminatedExitCode, ContainerStatusTerminatedReason, ContainerStatusTerminatedMessage,
		ContainerStatusTerminatedStartedAt, ContainerStatusTerminatedFinishedAt,
		ContainerLastTerminationStatusWaitingReason, ContainerLastTerminationStatusWaitingMessage,
		ContainerLastTerminationStatusRunningStartedAt,
		ContainerLastTerminationStatusTerminatedExitCode, ContainerLastTerminationStatusTerminatedReason, ContainerLastTerminationStatusTerminatedMessage,
		ContainerLastTerminationStatusTerminatedStartedAt, ContainerLastTerminationStatusTerminatedFinishedAt,
		ContainerRestartCount,
		ContainerResourceRequestCPU, ContainerResourceRequestMemory,
		ContainerResourceLimitCPU, ContainerResourceLimitMemory,
		ContainerPolicy,
		ContainerPodCreateTime, ContainerResourceLink,
		ContainerTopControllerName, ContainerTopControllerKind, ContainerTpoControllerReplicas,
		ContainerEnableHPA, ContainerEnableVPA,
	}
)

// ContainerEntity Entity in database
type ContainerEntity struct {
	Time                                      time.Time
	Namespace                                 *string
	PodName                                   *string
	PodPhase                                  *string
	PodMessage                                *string
	PodReason                                 *string
	AlamedaScalerNamespace                    *string
	AlamedaScalerName                         *string
	NodeName                                  *string
	Name                                      *string
	AppName                                   *string
	AppPartOf                                 *string
	StatusWaitingReason                       *string
	StatusWaitingMessage                      *string
	StatusRunningStartedAt                    *int64
	StatusTerminatedExitCode                  *int32
	StatusTerminatedReason                    *string
	StatusTerminatedMessage                   *string
	StatusTerminatedStartedAt                 *int64
	StatusTerminatedFinishedAt                *int64
	LastTerminationStatusWaitingReason        *string
	LastTerminationStatusWaitingMessage       *string
	LastTerminationStatusRunningStartedAt     *int64
	LastTerminationStatusTerminatedExitCode   *int32
	LastTerminationStatusTerminatedReason     *string
	LastTerminationStatusTerminatedMessage    *string
	LastTerminationStatusTerminatedStartedAt  *int64
	LastTerminationStatusTerminatedFinishedAt *int64
	RestartCount                              *int32
	ResourceRequestCPU                        *float64
	ResourceRequestMemory                     *int64
	ResourceLimitCPU                          *float64
	ResourceLimitMemory                       *int64
	Policy                                    *string
	PodCreatedTime                            *int64
	ResourceLink                              *string
	TopControllerName                         *string
	TopControllerKind                         *string
	TpoControllerReplicas                     *int32
	UsedRecommendationID                      *string
	EnableVPA                                 *bool
	EnableHPA                                 *bool
}

// NewContainerEntityFromMap Build entity from map
func NewContainerEntityFromMap(data map[string]string) ContainerEntity {

	// TODO: log error
	tempTimestamp, _ := utils.ParseTime(data[ContainerTime])

	entity := ContainerEntity{
		Time: tempTimestamp,
	}

	if namespace, exist := data[ContainerNamespace]; exist {
		entity.Namespace = &namespace
	}
	if podName, exist := data[ContainerPodName]; exist {
		entity.PodName = &podName
	}
	if podPhase, exist := data[ContainerPodPhase]; exist {
		entity.PodPhase = &podPhase
	}
	if podMessage, exist := data[ContainerPodMessage]; exist {
		entity.PodMessage = &podMessage
	}
	if podReason, exist := data[ContainerPodReason]; exist {
		entity.PodReason = &podReason
	}
	if alamedaScalerNamespace, exist := data[ContainerAlamedaScalerNamespace]; exist {
		entity.AlamedaScalerNamespace = &alamedaScalerNamespace
	}
	if alamedaScalerName, exist := data[ContainerAlamedaScalerName]; exist {
		entity.AlamedaScalerName = &alamedaScalerName
	}
	if nodeName, exist := data[ContainerNodeName]; exist {
		entity.NodeName = &nodeName
	}
	if name, exist := data[ContainerName]; exist {
		entity.Name = &name
	}
	if statusWaitingReason, exist := data[ContainerStatusWaitingReason]; exist {
		entity.StatusWaitingReason = &statusWaitingReason
	}
	if statusWaitingMessage, exist := data[ContainerStatusWaitingMessage]; exist {
		entity.StatusWaitingMessage = &statusWaitingMessage
	}
	if statusRunningStartedAt, exist := data[ContainerStatusRunningStartedAt]; exist {
		value, _ := strconv.ParseInt(statusRunningStartedAt, 10, 64)
		entity.StatusRunningStartedAt = &value
	}
	if statusTerminatedExitCode, exist := data[ContainerStatusTerminatedExitCode]; exist {
		temp, _ := strconv.ParseInt(statusTerminatedExitCode, 10, 64)
		value := int32(temp)
		entity.StatusTerminatedExitCode = &value
	}
	if statusTerminatedReason, exist := data[ContainerStatusTerminatedReason]; exist {
		entity.StatusTerminatedReason = &statusTerminatedReason
	}
	if statusTerminatedMessage, exist := data[ContainerStatusTerminatedMessage]; exist {
		entity.StatusTerminatedMessage = &statusTerminatedMessage
	}
	if statusTerminatedStartedAt, exist := data[ContainerStatusTerminatedStartedAt]; exist {
		value, _ := strconv.ParseInt(statusTerminatedStartedAt, 10, 64)
		entity.StatusTerminatedStartedAt = &value
	}
	if statusTerminatedFinishedAt, exist := data[ContainerStatusTerminatedFinishedAt]; exist {
		value, _ := strconv.ParseInt(statusTerminatedFinishedAt, 10, 64)
		entity.StatusTerminatedFinishedAt = &value
	}
	if lastTerminationStatusWaitingReason, exist := data[ContainerLastTerminationStatusWaitingReason]; exist {
		entity.LastTerminationStatusWaitingReason = &lastTerminationStatusWaitingReason
	}
	if lastTerminationStatusWaitingMessage, exist := data[ContainerLastTerminationStatusWaitingMessage]; exist {
		entity.LastTerminationStatusWaitingMessage = &lastTerminationStatusWaitingMessage
	}
	if lastTerminationStatusRunningStartedAt, exist := data[ContainerLastTerminationStatusRunningStartedAt]; exist {
		value, _ := strconv.ParseInt(lastTerminationStatusRunningStartedAt, 10, 64)
		entity.LastTerminationStatusRunningStartedAt = &value
	}
	if lastTerminationStatusTerminatedExitCode, exist := data[ContainerLastTerminationStatusTerminatedExitCode]; exist {
		temp, _ := strconv.ParseInt(lastTerminationStatusTerminatedExitCode, 10, 64)
		value := int32(temp)
		entity.LastTerminationStatusTerminatedExitCode = &value
	}
	if lastTerminationStatusTerminatedReason, exist := data[ContainerLastTerminationStatusTerminatedReason]; exist {
		entity.LastTerminationStatusTerminatedReason = &lastTerminationStatusTerminatedReason
	}
	if lastTerminationStatusTerminatedMessage, exist := data[ContainerLastTerminationStatusTerminatedMessage]; exist {
		entity.LastTerminationStatusTerminatedMessage = &lastTerminationStatusTerminatedMessage
	}
	if lastTerminationStatusTerminatedStartedAt, exist := data[ContainerLastTerminationStatusTerminatedStartedAt]; exist {
		value, _ := strconv.ParseInt(lastTerminationStatusTerminatedStartedAt, 10, 64)
		entity.LastTerminationStatusTerminatedStartedAt = &value
	}
	if lastTerminationStatusTerminatedFinishedAt, exist := data[ContainerLastTerminationStatusTerminatedFinishedAt]; exist {
		value, _ := strconv.ParseInt(lastTerminationStatusTerminatedFinishedAt, 10, 64)
		entity.LastTerminationStatusTerminatedFinishedAt = &value
	}
	if restartCount, exist := data[ContainerRestartCount]; exist {
		temp, _ := strconv.ParseInt(restartCount, 10, 64)
		value := int32(temp)
		entity.RestartCount = &value
	}
	if resourceRequestCPU, exist := data[ContainerResourceRequestCPU]; exist {
		value, _ := strconv.ParseFloat(resourceRequestCPU, 64)
		entity.ResourceRequestCPU = &value
	}
	if resourceRequestMemory, exist := data[ContainerResourceRequestMemory]; exist {
		value, _ := strconv.ParseInt(resourceRequestMemory, 10, 64)
		entity.ResourceRequestMemory = &value
	}
	if resourceLimitCPU, exist := data[ContainerResourceLimitCPU]; exist {
		value, _ := strconv.ParseFloat(resourceLimitCPU, 64)
		entity.ResourceLimitCPU = &value
	}
	if resourceLimitMemory, exist := data[ContainerResourceLimitMemory]; exist {
		value, _ := strconv.ParseInt(resourceLimitMemory, 10, 64)
		entity.ResourceLimitMemory = &value
	}
	if policy, exist := data[ContainerPolicy]; exist {
		entity.Policy = &policy
	}
	if podCreatedTime, exist := data[ContainerPodCreateTime]; exist {
		value, _ := strconv.ParseInt(podCreatedTime, 10, 64)
		entity.PodCreatedTime = &value
	}
	if resourceLink, exist := data[ContainerResourceLink]; exist {
		entity.ResourceLink = &resourceLink
	}
	if topControllerName, exist := data[ContainerTopControllerName]; exist {
		entity.TopControllerName = &topControllerName
	}
	if topControllerKind, exist := data[ContainerTopControllerKind]; exist {
		entity.TopControllerKind = &topControllerKind
	}
	if tpoControllerReplicas, exist := data[ContainerTpoControllerReplicas]; exist {
		temp, _ := strconv.ParseInt(tpoControllerReplicas, 10, 64)
		value := int32(temp)
		entity.TpoControllerReplicas = &value
	}
	if usedRecommendationID, exist := data[ContainerUsedRecommendationID]; exist {
		entity.UsedRecommendationID = &usedRecommendationID
	}
	if appName, exist := data[ContainerAppName]; exist {
		entity.AppName = &appName
	}
	if appPartOf, exist := data[ContainerAppPartOf]; exist {
		entity.AppPartOf = &appPartOf
	}
	if enableVPA, exist := data[ContainerEnableVPA]; exist {
		b, _ := strconv.ParseBool(enableVPA)
		entity.EnableVPA = &b
	}
	if enableHPA, exist := data[ContainerEnableHPA]; exist {
		b, _ := strconv.ParseBool(enableHPA)
		entity.EnableHPA = &b
	}

	return entity
}

func (e ContainerEntity) InfluxDBPoint(measurementName string) (*influxdb_client.Point, error) {

	tags := map[string]string{}
	if e.Namespace != nil {
		tags[ContainerNamespace] = *e.Namespace
	}
	if e.PodName != nil {
		tags[ContainerPodName] = *e.PodName
	}
	if e.NodeName != nil {
		tags[ContainerNodeName] = *e.NodeName
	}
	if e.Name != nil {
		tags[ContainerName] = *e.Name
	}
	if e.AlamedaScalerNamespace != nil {
		tags[ContainerAlamedaScalerNamespace] = *e.AlamedaScalerNamespace
	}
	if e.AlamedaScalerName != nil {
		tags[ContainerAlamedaScalerName] = *e.AlamedaScalerName
	}
	if e.AppName != nil {
		tags[ContainerAppName] = *e.AppName
	}
	if e.AppPartOf != nil {
		tags[ContainerAppPartOf] = *e.AppPartOf
	}

	fields := map[string]interface{}{}
	if e.PodPhase != nil {
		fields[ContainerPodPhase] = *e.PodPhase
	}
	if e.PodMessage != nil {
		fields[ContainerPodMessage] = *e.PodMessage
	}
	if e.PodReason != nil {
		fields[ContainerPodReason] = *e.PodReason
	}
	if e.StatusWaitingReason != nil {
		fields[ContainerStatusWaitingReason] = *e.StatusWaitingReason
	}
	if e.StatusWaitingMessage != nil {
		fields[ContainerStatusWaitingMessage] = *e.StatusWaitingMessage
	}
	if e.StatusRunningStartedAt != nil {
		fields[ContainerStatusRunningStartedAt] = *e.StatusRunningStartedAt
	}
	if e.StatusTerminatedExitCode != nil {
		fields[ContainerStatusTerminatedExitCode] = *e.StatusTerminatedExitCode
	}
	if e.StatusTerminatedReason != nil {
		fields[ContainerStatusTerminatedReason] = *e.StatusTerminatedReason
	}
	if e.StatusTerminatedMessage != nil {
		fields[ContainerStatusTerminatedMessage] = *e.StatusTerminatedMessage
	}
	if e.StatusTerminatedStartedAt != nil {
		fields[ContainerStatusTerminatedStartedAt] = *e.StatusTerminatedStartedAt
	}
	if e.StatusTerminatedFinishedAt != nil {
		fields[ContainerStatusTerminatedFinishedAt] = *e.StatusTerminatedFinishedAt
	}
	if e.LastTerminationStatusWaitingReason != nil {
		fields[ContainerLastTerminationStatusWaitingReason] = *e.LastTerminationStatusWaitingReason
	}
	if e.LastTerminationStatusWaitingMessage != nil {
		fields[ContainerLastTerminationStatusWaitingMessage] = *e.LastTerminationStatusWaitingMessage
	}
	if e.LastTerminationStatusRunningStartedAt != nil {
		fields[ContainerLastTerminationStatusRunningStartedAt] = *e.LastTerminationStatusRunningStartedAt
	}
	if e.LastTerminationStatusTerminatedExitCode != nil {
		fields[ContainerLastTerminationStatusTerminatedExitCode] = *e.LastTerminationStatusTerminatedExitCode
	}
	if e.LastTerminationStatusTerminatedReason != nil {
		fields[ContainerLastTerminationStatusTerminatedReason] = *e.LastTerminationStatusTerminatedReason
	}
	if e.LastTerminationStatusTerminatedMessage != nil {
		fields[ContainerLastTerminationStatusTerminatedMessage] = *e.LastTerminationStatusTerminatedMessage
	}
	if e.LastTerminationStatusTerminatedStartedAt != nil {
		fields[ContainerLastTerminationStatusTerminatedStartedAt] = *e.LastTerminationStatusTerminatedStartedAt
	}
	if e.LastTerminationStatusTerminatedFinishedAt != nil {
		fields[ContainerLastTerminationStatusTerminatedFinishedAt] = *e.LastTerminationStatusTerminatedFinishedAt
	}
	if e.RestartCount != nil {
		fields[ContainerRestartCount] = *e.RestartCount
	}
	if e.Policy != nil {
		fields[ContainerPolicy] = *e.Policy
	}
	if e.ResourceRequestCPU != nil {
		fields[ContainerResourceRequestCPU] = *e.ResourceRequestCPU
	}
	if e.ResourceRequestMemory != nil {
		fields[ContainerResourceRequestMemory] = *e.ResourceRequestMemory
	}
	if e.ResourceLimitCPU != nil {
		fields[ContainerResourceLimitCPU] = *e.ResourceLimitCPU
	}
	if e.ResourceLimitMemory != nil {
		fields[ContainerResourceLimitMemory] = *e.ResourceLimitMemory
	}
	if e.PodCreatedTime != nil {
		fields[ContainerPodCreateTime] = *e.PodCreatedTime
	}
	if e.ResourceLink != nil {
		fields[ContainerResourceLink] = *e.ResourceLink
	}
	if e.TopControllerName != nil {
		fields[ContainerTopControllerName] = *e.TopControllerName
	}
	if e.TopControllerKind != nil {
		fields[ContainerTopControllerKind] = *e.TopControllerKind
	}
	if e.TpoControllerReplicas != nil {
		fields[ContainerTpoControllerReplicas] = *e.TpoControllerReplicas
	}
	if e.UsedRecommendationID != nil {
		fields[ContainerUsedRecommendationID] = *e.UsedRecommendationID
	}
	if e.EnableVPA != nil {
		fields[ContainerEnableVPA] = *e.EnableVPA
	}
	if e.EnableHPA != nil {
		fields[ContainerEnableHPA] = *e.EnableHPA
	}

	return influxdb_client.NewPoint(measurementName, tags, fields, e.Time)
}
