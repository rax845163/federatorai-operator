package clusterstatus

import (
	"fmt"
	"strconv"
	"strings"

	cluster_status_entity "github.com/containers-ai/alameda/datahub/pkg/entity/influxdb/cluster_status"
	"github.com/containers-ai/alameda/datahub/pkg/entity/influxdb/utils/enumconv"
	"github.com/containers-ai/alameda/datahub/pkg/repository/influxdb"
	"github.com/containers-ai/alameda/pkg/utils/log"
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	"github.com/golang/protobuf/ptypes/timestamp"
	proto_timestmap "github.com/golang/protobuf/ptypes/timestamp"
	influxdb_client "github.com/influxdata/influxdb/client/v2"
	"github.com/pkg/errors"
)

var (
	scope = log.RegisterScope("cluster_status_db_measurement", "cluster_status DB measurement", 0)
)

// ContainerRepository is used to operate node measurement of cluster_status database
type ContainerRepository struct {
	influxDB *influxdb.InfluxDBRepository
}

// IsTag checks the column is tag or not
func (containerRepository *ContainerRepository) IsTag(column string) bool {
	for _, tag := range cluster_status_entity.ContainerTags {
		if column == string(tag) {
			return true
		}
	}
	return false
}

// NewContainerRepository creates the ContainerRepository instance
func NewContainerRepository(influxDBCfg *influxdb.Config) *ContainerRepository {
	return &ContainerRepository{
		influxDB: &influxdb.InfluxDBRepository{
			Address:  influxDBCfg.Address,
			Username: influxDBCfg.Username,
			Password: influxDBCfg.Password,
		},
	}
}

// ListAlamedaContainers list predicted containers have relation with arguments
func (containerRepository *ContainerRepository) ListAlamedaContainers(namespace, name string, kind datahub_v1alpha1.Kind, timeRange *datahub_v1alpha1.TimeRange) ([]*datahub_v1alpha1.Pod, error) {
	pods := []*datahub_v1alpha1.Pod{}
	whereStr := ""

	relationStatement := ""
	podCreatePeriodCondition := containerRepository.getPodCreatePeriodCondition(timeRange)

	switch kind {
	// bypass if Kind is Pod
	case datahub_v1alpha1.Kind_POD:
		// relationStatement = fmt.Sprintf(`("%s" = '%s' AND "%s" = '%s')`,
		// 	cluster_status_entity.ContainerNamespace, namespace,
		// 	cluster_status_entity.ContainerPodName, name)
		if podCreatePeriodCondition != "" {
			relationStatement = strings.Replace(podCreatePeriodCondition, "AND", "", 1)
		}

	case datahub_v1alpha1.Kind_DEPLOYMENT:
		relationStatement = fmt.Sprintf(`("%s" = '%s' AND "%s" = '%s' AND "%s" = '%s' %s)`,
			cluster_status_entity.ContainerNamespace, namespace,
			cluster_status_entity.ContainerTopControllerName, name,
			cluster_status_entity.ContainerTopControllerKind, enumconv.KindDisp[datahub_v1alpha1.Kind_DEPLOYMENT],
			podCreatePeriodCondition)
	case datahub_v1alpha1.Kind_DEPLOYMENTCONFIG:
		relationStatement = fmt.Sprintf(`("%s" = '%s' AND "%s" = '%s' AND "%s" = '%s' %s)`,
			cluster_status_entity.ContainerNamespace, namespace,
			cluster_status_entity.ContainerTopControllerName, name,
			cluster_status_entity.ContainerTopControllerKind, enumconv.KindDisp[datahub_v1alpha1.Kind_DEPLOYMENTCONFIG],
			podCreatePeriodCondition)
	case datahub_v1alpha1.Kind_ALAMEDASCALER:
		relationStatement = fmt.Sprintf(`("%s" = '%s' AND "%s" = '%s' %s)`,
			cluster_status_entity.ContainerAlamedaScalerNamespace, namespace,
			cluster_status_entity.ContainerAlamedaScalerName, name,
			podCreatePeriodCondition)
	default:
		return pods, errors.Errorf("no mapping filter statement with Datahub Kind: %s, skip building relation statement", datahub_v1alpha1.Kind_name[int32(kind)])
	}
	if relationStatement != "" {
		if whereStr != "" {
			whereStr = fmt.Sprintf("%s AND %s", whereStr, relationStatement)
		} else {
			whereStr = fmt.Sprintf("WHERE %s", relationStatement)
		}
	}

	cmd := fmt.Sprintf("SELECT * FROM %s %s GROUP BY \"%s\",\"%s\",\"%s\",\"%s\"",
		string(Container), whereStr,
		string(cluster_status_entity.ContainerNamespace), string(cluster_status_entity.ContainerPodName),
		string(cluster_status_entity.ContainerAlamedaScalerNamespace), string(cluster_status_entity.ContainerAlamedaScalerName))
	scope.Debugf("ListAlamedaContainers command: %s", cmd)
	if results, err := containerRepository.influxDB.QueryDB(cmd, string(influxdb.ClusterStatus)); err == nil {

		containerEntities := make([]*cluster_status_entity.ContainerEntity, 0)

		rows := influxdb.PackMap(results)
		for _, row := range rows {
			for _, data := range row.Data {
				entity := cluster_status_entity.NewContainerEntityFromMap(data)
				containerEntities = append(containerEntities, &entity)
			}
		}

		pods = buildDatahubPodsFromContainerEntities(containerEntities)
		return pods, nil
	} else {
		return pods, err
	}
}

// CreateContainers add containers information container measurement
func (containerRepository *ContainerRepository) CreateContainers(pods []*datahub_v1alpha1.Pod) error {
	points := []*influxdb_client.Point{}
	for _, pod := range pods {
		containerEntities, err := buildContainerEntitiesFromDatahubPod(pod)
		if err != nil {
			return errors.Wrap(err, "create containers failed")
		}
		for _, containerEntity := range containerEntities {
			p, err := containerEntity.InfluxDBPoint(string(Container))
			if err != nil {
				return errors.Wrap(err, "create containers failed")
			}
			points = append(points, p)
		}
	}
	err := containerRepository.influxDB.WritePoints(points, influxdb_client.BatchPointsConfig{
		Database: string(influxdb.ClusterStatus),
	})
	if err != nil {
		return errors.Wrap(err, "create containers to influxdb failed")
	}
	return nil
}

// DeleteContainers set containers' field is_deleted to true into container measurement
func (containerRepository *ContainerRepository) DeleteContainers(pods []*datahub_v1alpha1.Pod) error {
	for _, pod := range pods {
		if pod.GetNamespacedName() == nil || pod.GetAlamedaScaler() == nil {
			continue
		}
		podNS := pod.GetNamespacedName().GetNamespace()
		podName := pod.GetNamespacedName().GetName()
		alaScalerNS := pod.GetAlamedaScaler().GetNamespace()
		alaScalerName := pod.GetAlamedaScaler().GetName()
		cmd := fmt.Sprintf("DROP SERIES FROM %s WHERE \"%s\"='%s' AND \"%s\"='%s' AND \"%s\"='%s' AND \"%s\"='%s'", Container,
			cluster_status_entity.ContainerNamespace, podNS, cluster_status_entity.ContainerPodName, podName,
			cluster_status_entity.ContainerAlamedaScalerNamespace, alaScalerNS, cluster_status_entity.ContainerAlamedaScalerName, alaScalerName)
		scope.Debugf("DeleteContainers command: %s", cmd)
		_, err := containerRepository.influxDB.QueryDB(cmd, string(influxdb.ClusterStatus))
		if err != nil {
			scope.Errorf(err.Error())
		}
	}
	return nil
}

// ListPodsContainers list containers information container measurement
func (containerRepository *ContainerRepository) ListPodsContainers(pods []*datahub_v1alpha1.Pod) ([]*cluster_status_entity.ContainerEntity, error) {

	var (
		cmd                 = ""
		cmdSelectString     = ""
		cmdTagsFilterString = ""
		containerEntities   = make([]*cluster_status_entity.ContainerEntity, 0)
	)

	if len(pods) == 0 {
		return containerEntities, nil
	}

	cmdSelectString = fmt.Sprintf(`select * from "%s" `, Container)
	for _, pod := range pods {

		var (
			namespace = ""
			podName   = ""
		)

		if pod.GetNamespacedName() != nil {
			namespace = pod.GetNamespacedName().GetNamespace()
			podName = pod.GetNamespacedName().GetName()
		}

		cmdTagsFilterString += fmt.Sprintf(`("%s" = '%s' and "%s" = '%s') or `,
			cluster_status_entity.ContainerNamespace, namespace,
			cluster_status_entity.ContainerPodName, podName,
		)
	}
	cmdTagsFilterString = strings.TrimSuffix(cmdTagsFilterString, "or ")

	cmd = fmt.Sprintf("%s where %s", cmdSelectString, cmdTagsFilterString)
	results, err := containerRepository.influxDB.QueryDB(cmd, string(influxdb.ClusterStatus))
	if err != nil {
		return containerEntities, errors.Wrap(err, "list pod containers failed")
	}

	rows := influxdb.PackMap(results)
	for _, row := range rows {
		for _, data := range row.Data {
			entity := cluster_status_entity.NewContainerEntityFromMap(data)
			containerEntities = append(containerEntities, &entity)
		}
	}

	return containerEntities, nil
}

func (containerRepository *ContainerRepository) getPodCreatePeriodCondition(timeRange *datahub_v1alpha1.TimeRange) string {
	if timeRange == nil {
		return ""
	}

	var start int64 = 0
	var end int64 = 0

	if timeRange.StartTime != nil {
		start = timeRange.StartTime.Seconds
	}

	if timeRange.EndTime != nil {
		end = timeRange.EndTime.Seconds
	}

	if start == 0 && end == 0 {
		return ""
	} else if start == 0 && end != 0 {
		period := fmt.Sprintf(`AND "pod_create_time" < %d`, end)
		return period
	} else if start != 0 && end == 0 {
		period := fmt.Sprintf(`AND "pod_create_time" >= %d`, start)
		return period
	} else if start != 0 && end != 0 {
		period := fmt.Sprintf(`AND "pod_create_time" >= %d AND "pod_create_time" < %d`, start, end)
		return period
	}

	return ""
}

func buildContainerEntitiesFromDatahubPod(pod *datahub_v1alpha1.Pod) ([]*cluster_status_entity.ContainerEntity, error) {

	var (
		namespace                                 *string
		podName                                   *string
		podPhase                                  *string
		podMessage                                *string
		podReason                                 *string
		alamedaScalerNamespace                    *string
		alamedaScalerName                         *string
		nodeName                                  *string
		name                                      *string
		statusWaitingReason                       *string
		statusWaitingMessage                      *string
		statusRunningStartedAt                    *int64
		statusTerminatedExitCode                  *int32
		statusTerminatedReason                    *string
		statusTerminatedMessage                   *string
		statusTerminatedStartedAt                 *int64
		statusTerminatedFinishedAt                *int64
		lastTerminationStatusWaitingReason        *string
		lastTerminationStatusWaitingMessage       *string
		lastTerminationStatusRunningStartedAt     *int64
		lastTerminationStatusTerminatedExitCode   *int32
		lastTerminationStatusTerminatedReason     *string
		lastTerminationStatusTerminatedMessage    *string
		lastTerminationStatusTerminatedStartedAt  *int64
		lastTerminationStatusTerminatedFinishedAt *int64
		restartCount                              *int32
		resourceRequestCPU                        *float64
		resourceRequestMemory                     *int64
		resourceLimitCPU                          *float64
		resourceLimitMemory                       *int64
		policy                                    *string
		podCreatedTime                            *int64
		resourceLink                              *string
		topControllerName                         *string
		topControllerKind                         *string
		topControllerReplicas                     *int32
		usedRecommendationID                      *string
	)

	nodeName = &pod.NodeName
	resourceLink = &pod.ResourceLink
	usedRecommendationID = &pod.UsedRecommendationId

	if pod.NamespacedName != nil {
		namespace = &pod.NamespacedName.Namespace
		podName = &pod.NamespacedName.Name
	}
	if pod.AlamedaScaler != nil {
		alamedaScalerNamespace = &pod.AlamedaScaler.Namespace
		alamedaScalerName = &pod.AlamedaScaler.Name
	}
	if pod.StartTime != nil {
		startTime := pod.StartTime.GetSeconds()
		podCreatedTime = &startTime
	}
	if pod.TopController != nil {
		if pod.TopController.NamespacedName != nil {
			topControllerName = &pod.TopController.NamespacedName.Name
		}
		if k, exist := enumconv.KindDisp[pod.TopController.Kind]; exist {
			topControllerKind = &k
		}
		topControllerReplicas = &pod.TopController.Replicas
	}
	if p, exist := enumconv.RecommendationPolicyDisp[pod.GetPolicy()]; exist {
		policy = &p
	}
	if pod.Status != nil {
		if val, ok := datahub_v1alpha1.PodPhase_name[int32(pod.Status.Phase)]; ok {
			podPhase = &val
		} else {
			val = datahub_v1alpha1.PodPhase_name[int32(datahub_v1alpha1.PodPhase_Unknown)]
			podPhase = &val
		}
		podMessage = &pod.Status.Message
		podReason = &pod.Status.Reason
	}

	appName := pod.GetAppName()
	appPartOf := pod.GetAppPartOf()
	enableVPA := pod.GetEnable_VPA()
	enableHPA := pod.GetEnable_HPA()

	containerEntities := make([]*cluster_status_entity.ContainerEntity, 0)
	for _, datahubContainer := range pod.Containers {

		name = &datahubContainer.Name

		for _, metricData := range datahubContainer.GetLimitResource() {
			if data := metricData.GetData(); len(data) == 1 {
				switch metricData.GetMetricType() {
				case datahub_v1alpha1.MetricType_CPU_USAGE_SECONDS_PERCENTAGE:
					val, err := strconv.ParseFloat(data[0].NumValue, 64)
					if err != nil {
						scope.Warnf("convert string: %s to float64 failed, skip assigning value, err: %s", data[0].NumValue, err.Error())
					}
					resourceLimitCPU = &val
				case datahub_v1alpha1.MetricType_MEMORY_USAGE_BYTES:
					val, err := strconv.ParseInt(data[0].NumValue, 10, 64)
					if err != nil {
						scope.Warnf("convert string: %s to int64 failed, skip assigning value, err: %s", data[0].NumValue, err.Error())
					}
					resourceLimitMemory = &val
				default:
					scope.Warnf("no mapping metric type for Datahub.MetricType: %s, skip assigning value", datahub_v1alpha1.MetricType_name[int32(metricData.GetMetricType())])
				}
			}
		}
		for _, metricData := range datahubContainer.GetRequestResource() {
			if data := metricData.GetData(); len(data) == 1 {
				switch metricData.GetMetricType() {
				case datahub_v1alpha1.MetricType_CPU_USAGE_SECONDS_PERCENTAGE:
					val, err := strconv.ParseFloat(data[0].NumValue, 64)
					if err != nil {
						scope.Warnf("convert string: %s to float64 failed, skip assigning value, err: %s", data[0].NumValue, err.Error())
					}
					resourceRequestCPU = &val
				case datahub_v1alpha1.MetricType_MEMORY_USAGE_BYTES:
					val, err := strconv.ParseInt(data[0].NumValue, 10, 64)
					if err != nil {
						scope.Warnf("convert string: %s to int64 failed, skip assigning value, err: %s", data[0].NumValue, err.Error())
					}
					resourceRequestMemory = &val
				default:
					scope.Warnf("no mapping metric type for Datahub.MetricType: %s, skip assigning value", datahub_v1alpha1.MetricType_name[int32(metricData.GetMetricType())])
				}
			}
		}
		if datahubContainer.GetStatus() != nil {
			containerStatus := datahubContainer.GetStatus()
			if containerStatus.GetState() != nil {
				state := containerStatus.GetState()
				if state.GetWaiting() != nil {
					statusWaitingReason = &state.GetWaiting().Reason
					statusWaitingMessage = &state.GetWaiting().Message
				}
				if state.GetRunning() != nil {
					statusRunningStartedAt = &state.GetRunning().GetStartedAt().Seconds
				}
				if state.GetTerminated() != nil {
					statusTerminatedExitCode = &state.GetTerminated().ExitCode
					statusTerminatedReason = &state.GetTerminated().Reason
					statusTerminatedMessage = &state.GetTerminated().Message
					statusTerminatedStartedAt = &state.GetTerminated().GetStartedAt().Seconds
					statusTerminatedFinishedAt = &state.GetTerminated().GetFinishedAt().Seconds
				}
			}
			if containerStatus.GetLastTerminationState() != nil {
				state := containerStatus.GetLastTerminationState()
				if state.GetWaiting() != nil {
					lastTerminationStatusWaitingReason = &state.GetWaiting().Reason
					lastTerminationStatusWaitingMessage = &state.GetWaiting().Message
				}
				if state.GetRunning() != nil {
					lastTerminationStatusRunningStartedAt = &state.GetRunning().GetStartedAt().Seconds
				}
				if state.GetTerminated() != nil {
					lastTerminationStatusTerminatedExitCode = &state.GetTerminated().ExitCode
					lastTerminationStatusTerminatedReason = &state.GetTerminated().Reason
					lastTerminationStatusTerminatedMessage = &state.GetTerminated().Message
					lastTerminationStatusTerminatedStartedAt = &state.GetTerminated().GetStartedAt().Seconds
					lastTerminationStatusTerminatedFinishedAt = &state.GetTerminated().GetFinishedAt().Seconds
				}
			}
			restartCount = &containerStatus.RestartCount
		}

		containerEntity := &cluster_status_entity.ContainerEntity{
			Time:                                      influxdb.ZeroTime,
			Namespace:                                 namespace,
			PodName:                                   podName,
			PodPhase:                                  podPhase,
			PodMessage:                                podMessage,
			PodReason:                                 podReason,
			AlamedaScalerNamespace:                    alamedaScalerNamespace,
			AlamedaScalerName:                         alamedaScalerName,
			NodeName:                                  nodeName,
			Name:                                      name,
			StatusWaitingReason:                       statusWaitingReason,
			StatusWaitingMessage:                      statusWaitingMessage,
			StatusRunningStartedAt:                    statusRunningStartedAt,
			StatusTerminatedExitCode:                  statusTerminatedExitCode,
			StatusTerminatedReason:                    statusTerminatedReason,
			StatusTerminatedMessage:                   statusTerminatedMessage,
			StatusTerminatedStartedAt:                 statusTerminatedStartedAt,
			StatusTerminatedFinishedAt:                statusTerminatedFinishedAt,
			LastTerminationStatusWaitingReason:        lastTerminationStatusWaitingReason,
			LastTerminationStatusWaitingMessage:       lastTerminationStatusWaitingMessage,
			LastTerminationStatusRunningStartedAt:     lastTerminationStatusRunningStartedAt,
			LastTerminationStatusTerminatedExitCode:   lastTerminationStatusTerminatedExitCode,
			LastTerminationStatusTerminatedReason:     lastTerminationStatusTerminatedReason,
			LastTerminationStatusTerminatedMessage:    lastTerminationStatusTerminatedMessage,
			LastTerminationStatusTerminatedStartedAt:  lastTerminationStatusTerminatedStartedAt,
			LastTerminationStatusTerminatedFinishedAt: lastTerminationStatusTerminatedFinishedAt,
			RestartCount:                              restartCount,
			ResourceRequestCPU:                        resourceRequestCPU,
			ResourceRequestMemory:                     resourceRequestMemory,
			ResourceLimitCPU:                          resourceLimitCPU,
			ResourceLimitMemory:                       resourceLimitMemory,
			Policy:                                    policy,
			PodCreatedTime:                            podCreatedTime,
			ResourceLink:                              resourceLink,
			TopControllerName:                         topControllerName,
			TopControllerKind:                         topControllerKind,
			TpoControllerReplicas:                     topControllerReplicas,
			UsedRecommendationID:                      usedRecommendationID,
			AppName:                                   &appName,
			AppPartOf:                                 &appPartOf,
			EnableHPA:                                 &enableHPA,
			EnableVPA:                                 &enableVPA,
		}
		containerEntities = append(containerEntities, containerEntity)
	}
	return containerEntities, nil
}

func buildDatahubPodsFromContainerEntities(containerEntities []*cluster_status_entity.ContainerEntity) []*datahub_v1alpha1.Pod {

	datahubPods := make([]*datahub_v1alpha1.Pod, 0)
	datahubPodMap := make(map[string]*datahub_v1alpha1.Pod)

	for _, containerEntity := range containerEntities {

		podID := getDatahubPodIDString(containerEntity)

		var datahubPod *datahub_v1alpha1.Pod
		datahubPod, exist := datahubPodMap[podID]
		if !exist {

			var (
				podName                string
				podPhase               string
				podMessage             string
				podReason              string
				namespace              string
				resourceLink           string
				alamedaScalerNamespace string
				alamedaScalerName      string
				nodeName               string
				podCreatedTime         int64
				policy                 string
				topControllerNamespace string
				topControllerName      string
				topControllerKind      string
				topControllerReplicas  int32
				usedRecommendationID   string
				appName                string
				appPartOf              string
				enableHPA              bool
				enableVPA              bool
			)

			if containerEntity.PodName != nil {
				podName = *containerEntity.PodName
			}
			if containerEntity.PodPhase != nil {
				podPhase = *containerEntity.PodPhase
			}
			if containerEntity.PodMessage != nil {
				podMessage = *containerEntity.PodMessage
			}
			if containerEntity.PodReason != nil {
				podReason = *containerEntity.PodReason
			}
			if containerEntity.Namespace != nil {
				namespace = *containerEntity.Namespace
			}
			if containerEntity.ResourceLink != nil {
				resourceLink = *containerEntity.ResourceLink
			}
			if containerEntity.AlamedaScalerNamespace != nil {
				alamedaScalerNamespace = *containerEntity.AlamedaScalerNamespace
			}
			if containerEntity.AlamedaScalerName != nil {
				alamedaScalerName = *containerEntity.AlamedaScalerName
			}
			if containerEntity.NodeName != nil {
				nodeName = *containerEntity.NodeName
			}
			if containerEntity.PodCreatedTime != nil {
				podCreatedTime = *containerEntity.PodCreatedTime
			}
			if containerEntity.Policy != nil {
				policy = *containerEntity.Policy
			}
			if containerEntity.Namespace != nil {
				topControllerNamespace = *containerEntity.Namespace
			}
			if containerEntity.TopControllerName != nil {
				topControllerName = *containerEntity.TopControllerName
			}
			if containerEntity.TopControllerKind != nil {
				topControllerKind = *containerEntity.TopControllerKind
			}
			if containerEntity.TpoControllerReplicas != nil {
				topControllerReplicas = *containerEntity.TpoControllerReplicas
			}
			if containerEntity.UsedRecommendationID != nil {
				usedRecommendationID = *containerEntity.UsedRecommendationID
			}

			if containerEntity.AppName != nil {
				appName = *containerEntity.AppName
			}
			if containerEntity.AppPartOf != nil {
				appPartOf = *containerEntity.AppPartOf
			}
			if containerEntity.EnableHPA != nil {
				enableHPA = *containerEntity.EnableHPA
			}
			if containerEntity.EnableVPA != nil {
				enableVPA = *containerEntity.EnableVPA
			}

			datahubPod = &datahub_v1alpha1.Pod{
				NamespacedName: &datahub_v1alpha1.NamespacedName{
					Name:      podName,
					Namespace: namespace,
				},
				ResourceLink: resourceLink,
				Containers:   make([]*datahub_v1alpha1.Container, 0),
				AlamedaScaler: &datahub_v1alpha1.NamespacedName{
					Name:      alamedaScalerName,
					Namespace: alamedaScalerNamespace,
				},
				NodeName: nodeName,
				StartTime: &proto_timestmap.Timestamp{
					Seconds: podCreatedTime,
				},
				Policy: enumconv.RecommendationPolicyEnum[policy],
				TopController: &datahub_v1alpha1.TopController{
					NamespacedName: &datahub_v1alpha1.NamespacedName{
						Namespace: topControllerNamespace,
						Name:      topControllerName,
					},
					Replicas: topControllerReplicas,
					Kind:     enumconv.KindEnum[topControllerKind],
				},
				UsedRecommendationId: usedRecommendationID,
				Status: &datahub_v1alpha1.PodStatus{
					Phase:   datahub_v1alpha1.PodPhase(datahub_v1alpha1.PodPhase_value[podPhase]),
					Message: podMessage,
					Reason:  podReason,
				},
				AppName:    appName,
				AppPartOf:  appPartOf,
				Enable_HPA: enableHPA,
				Enable_VPA: enableVPA,
			}
			datahubPodMap[podID] = datahubPod
		}

		datahubContainer := containerEntityToDatahubContainer(containerEntity)
		datahubPod.Containers = append(datahubPod.Containers, datahubContainer)
	}

	for _, datahubPod := range datahubPodMap {
		copyDatahubPod := datahubPod
		datahubPods = append(datahubPods, copyDatahubPod)
	}

	return datahubPods
}

func containerEntityToDatahubContainer(containerEntity *cluster_status_entity.ContainerEntity) *datahub_v1alpha1.Container {

	var (
		statusWaitingReason                       string
		statusWaitingMessage                      string
		statusRunningStartedAt                    int64
		statusTerminatedExitCode                  int32
		statusTerminatedReason                    string
		statusTerminatedMessage                   string
		statusTerminatedStartedAt                 int64
		statusTerminatedFinishedAt                int64
		lastTerminationStatusWaitingReason        string
		lastTerminationStatusWaitingMessage       string
		lastTerminationStatusRunningStartedAt     int64
		lastTerminationStatusTerminatedExitCode   int32
		lastTerminationStatusTerminatedReason     string
		lastTerminationStatusTerminatedMessage    string
		lastTerminationStatusTerminatedStartedAt  int64
		lastTerminationStatusTerminatedFinishedAt int64
		restartCount                              int32
		resourceLimitCPU                          float64
		resourceLimitMemory                       int64
		resourceRequestCPU                        float64
		resourceRequestMemory                     int64
	)

	if containerEntity.StatusWaitingReason != nil {
		statusWaitingReason = *containerEntity.StatusWaitingReason
	}
	if containerEntity.StatusWaitingMessage != nil {
		statusWaitingMessage = *containerEntity.StatusWaitingMessage
	}
	if containerEntity.StatusRunningStartedAt != nil {
		statusRunningStartedAt = *containerEntity.StatusRunningStartedAt
	}
	if containerEntity.StatusTerminatedExitCode != nil {
		statusTerminatedExitCode = *containerEntity.StatusTerminatedExitCode
	}
	if containerEntity.StatusTerminatedReason != nil {
		statusTerminatedReason = *containerEntity.StatusTerminatedReason
	}
	if containerEntity.StatusTerminatedMessage != nil {
		statusTerminatedMessage = *containerEntity.StatusTerminatedMessage
	}
	if containerEntity.StatusTerminatedStartedAt != nil {
		statusTerminatedStartedAt = *containerEntity.StatusTerminatedStartedAt
	}
	if containerEntity.StatusTerminatedFinishedAt != nil {
		statusTerminatedFinishedAt = *containerEntity.StatusTerminatedFinishedAt
	}
	if containerEntity.LastTerminationStatusWaitingReason != nil {
		lastTerminationStatusWaitingReason = *containerEntity.LastTerminationStatusWaitingReason
	}
	if containerEntity.LastTerminationStatusWaitingMessage != nil {
		lastTerminationStatusWaitingMessage = *containerEntity.LastTerminationStatusWaitingMessage
	}
	if containerEntity.LastTerminationStatusRunningStartedAt != nil {
		lastTerminationStatusRunningStartedAt = *containerEntity.LastTerminationStatusRunningStartedAt
	}
	if containerEntity.LastTerminationStatusTerminatedExitCode != nil {
		lastTerminationStatusTerminatedExitCode = *containerEntity.LastTerminationStatusTerminatedExitCode
	}
	if containerEntity.LastTerminationStatusTerminatedReason != nil {
		lastTerminationStatusTerminatedReason = *containerEntity.LastTerminationStatusTerminatedReason
	}
	if containerEntity.LastTerminationStatusTerminatedMessage != nil {
		lastTerminationStatusTerminatedMessage = *containerEntity.LastTerminationStatusTerminatedMessage
	}
	if containerEntity.LastTerminationStatusTerminatedStartedAt != nil {
		lastTerminationStatusTerminatedStartedAt = *containerEntity.LastTerminationStatusTerminatedStartedAt
	}
	if containerEntity.LastTerminationStatusTerminatedFinishedAt != nil {
		lastTerminationStatusTerminatedFinishedAt = *containerEntity.LastTerminationStatusTerminatedFinishedAt
	}
	if containerEntity.RestartCount != nil {
		restartCount = *containerEntity.RestartCount
	}
	if containerEntity.ResourceLimitCPU != nil {
		resourceLimitCPU = *containerEntity.ResourceLimitCPU
	}
	if containerEntity.ResourceLimitMemory != nil {
		resourceLimitMemory = *containerEntity.ResourceLimitMemory
	}
	if containerEntity.ResourceRequestCPU != nil {
		resourceRequestCPU = *containerEntity.ResourceRequestCPU
	}
	if containerEntity.ResourceRequestMemory != nil {
		resourceRequestMemory = *containerEntity.ResourceRequestMemory
	}

	return &datahub_v1alpha1.Container{
		Name: *containerEntity.Name,
		LimitResource: []*datahub_v1alpha1.MetricData{
			&datahub_v1alpha1.MetricData{
				MetricType: datahub_v1alpha1.MetricType_CPU_USAGE_SECONDS_PERCENTAGE,
				Data: []*datahub_v1alpha1.Sample{
					&datahub_v1alpha1.Sample{
						NumValue: strconv.FormatFloat(resourceLimitCPU, 'f', -1, 64),
					},
				},
			},
			&datahub_v1alpha1.MetricData{
				MetricType: datahub_v1alpha1.MetricType_MEMORY_USAGE_BYTES,
				Data: []*datahub_v1alpha1.Sample{
					&datahub_v1alpha1.Sample{
						NumValue: strconv.FormatInt(resourceLimitMemory, 10),
					},
				},
			},
		},
		RequestResource: []*datahub_v1alpha1.MetricData{
			&datahub_v1alpha1.MetricData{
				MetricType: datahub_v1alpha1.MetricType_CPU_USAGE_SECONDS_PERCENTAGE,
				Data: []*datahub_v1alpha1.Sample{
					&datahub_v1alpha1.Sample{
						NumValue: strconv.FormatFloat(resourceRequestCPU, 'f', -1, 64),
					},
				},
			},
			&datahub_v1alpha1.MetricData{
				MetricType: datahub_v1alpha1.MetricType_MEMORY_USAGE_BYTES,
				Data: []*datahub_v1alpha1.Sample{
					&datahub_v1alpha1.Sample{
						NumValue: strconv.FormatInt(resourceRequestMemory, 10),
					},
				},
			},
		},
		Status: &datahub_v1alpha1.ContainerStatus{
			State: &datahub_v1alpha1.ContainerState{
				Waiting: &datahub_v1alpha1.ContainerStateWaiting{
					Reason:  statusWaitingReason,
					Message: statusWaitingMessage,
				},
				Running: &datahub_v1alpha1.ContainerStateRunning{
					StartedAt: &timestamp.Timestamp{Seconds: statusRunningStartedAt},
				},
				Terminated: &datahub_v1alpha1.ContainerStateTerminated{
					ExitCode:   statusTerminatedExitCode,
					Reason:     statusTerminatedReason,
					Message:    statusTerminatedMessage,
					StartedAt:  &timestamp.Timestamp{Seconds: statusTerminatedStartedAt},
					FinishedAt: &timestamp.Timestamp{Seconds: statusTerminatedFinishedAt},
				},
			},
			LastTerminationState: &datahub_v1alpha1.ContainerState{
				Waiting: &datahub_v1alpha1.ContainerStateWaiting{
					Reason:  lastTerminationStatusWaitingReason,
					Message: lastTerminationStatusWaitingMessage,
				},
				Running: &datahub_v1alpha1.ContainerStateRunning{
					StartedAt: &timestamp.Timestamp{Seconds: lastTerminationStatusRunningStartedAt},
				},
				Terminated: &datahub_v1alpha1.ContainerStateTerminated{
					ExitCode:   lastTerminationStatusTerminatedExitCode,
					Reason:     lastTerminationStatusTerminatedReason,
					Message:    lastTerminationStatusTerminatedMessage,
					StartedAt:  &timestamp.Timestamp{Seconds: lastTerminationStatusTerminatedStartedAt},
					FinishedAt: &timestamp.Timestamp{Seconds: lastTerminationStatusTerminatedFinishedAt},
				},
			},
			RestartCount: restartCount,
		},
	}
}

func getDatahubPodIDString(containerEntity *cluster_status_entity.ContainerEntity) string {

	var (
		namespace         string
		podName           string
		alamedaScalerNS   string
		alamedaScalerName string
		nodeName          string
	)

	if containerEntity.Namespace != nil {
		namespace = *containerEntity.Namespace
	}
	if containerEntity.PodName != nil {
		podName = *containerEntity.PodName
	}
	if containerEntity.AlamedaScalerNamespace != nil {
		alamedaScalerNS = *containerEntity.AlamedaScalerNamespace
	}
	if containerEntity.AlamedaScalerName != nil {
		alamedaScalerName = *containerEntity.AlamedaScalerName
	}
	if containerEntity.NodeName != nil {
		nodeName = *containerEntity.NodeName
	}

	return fmt.Sprintf("%s.%s.%s.%s.%s", namespace, podName, alamedaScalerNS, alamedaScalerName, nodeName)
}
