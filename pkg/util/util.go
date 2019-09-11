package util

import (
	"context"
	"fmt"
	"strings"

	"github.com/containers-ai/federatorai-operator/pkg/apis/federatorai/v1alpha1"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

type GroupEnums string

const (
	FedemeterGroup GroupEnums = "alameda/fedemeter"
	AlamedaGroup   GroupEnums = "alameda"
	GrafanaGroup   GroupEnums = "grafana"
	InfluxDBGroup  GroupEnums = "influxdb"
	//deployment name
	AlamedaaiDPN              = "alameda-ai"
	AlamedaoperatorDPN        = "alameda-operator"
	AlamedadatahubDPN         = "alameda-datahub"
	AlamedaevictionerDPN      = "alameda-evictioner"
	AdmissioncontrollerDPN    = "admission-controller"
	AlamedarecommenderDPN     = "alameda-recommender"
	AlamedaexecutorDPN        = "alameda-executor"
	AlamedadispatcherDPN      = "alameda-ai-dispatcher"
	AlamedaRabbitMQDPN        = "alameda-rabbitmq"
	FedemeterDPN              = "fedemeter-api"
	GrafanaDPN                = "alameda-grafana"
	InfluxdbDPN               = "alameda-influxdb"
	AlamedaweavescopeDPN      = "alameda-weave-scope-app"
	AlamedaweavescopeProbeDPN = "alameda-weave-scope-cluster-agent"
	AlamedaanalyzerDPN        = "alameda-analyzer"
	AlamedaNotifierDPN        = "alameda-notifier"
	FederatoraiAgentDPN       = "federatorai-agent"
	//DaemonSet name
	AlamedaweavescopeAgentDS = "alameda-weave-scope-agent"
	//container name
	AlamedaaiCTN              = "alameda-ai-engine"
	AlamedaoperatorCTN        = "alameda-operator"
	AlamedadatahubCTN         = "alameda-datahub"
	AlamedaevictionerCTN      = "alameda-evictioner"
	AdmissioncontrollerCTN    = "admission-controller"
	AlamedarecommenderCTN     = "alameda-recommender"
	AlamedaexecutorCTN        = "alameda-executor"
	AlamedadispatcherCTN      = "ai-dispatcher"
	AlamedaRabbitMQCTN        = "rabbitmq"
	FedemeterCTN              = "fedemeter-api"
	FedemeterInfluxDBCTN      = "fedemeter-influxdb"
	GetTokenCTN               = "gettoken"
	GrafanaCTN                = "grafana"
	InfluxdbCTN               = "influxdb"
	AlamedaweavescopeCTN      = "alameda-weave-scope-app"
	AlamedaweavescopeProbeCTN = "alameda-weave-scope-cluster-agent"
	AlamedaweavescopeAgentCTN = "alameda-weave-scope-agent"
	AlamedaanalyzerCTN        = "alameda-analyzer"
	AlamedaNofitierCTN        = "alameda-notifier"
	FederatoraiAgentCTN       = "federatorai-agent"
	//Statefulset name
	FedemeterInflixDBSSN = "fedemeter-influxdb"
	//CRD NAME
	AlamedaScalerName         = "alamedascalers.autoscaling.containers.ai"
	AlamedaRecommendationName = "alamedarecommendations.autoscaling.containers.ai"
	// Service name
	AdmissionControllerServiceName  = "admission-controller"
	AlamedaAIServiceName            = "alameda-ai"
	AlamedaDatahubServiceName       = "alameda-datahub"
	AlamedaGrafanaServiceName       = "alameda-grafana"
	AlamedaInfluxdbServiceName      = "alameda-influxdb"
	AlamedaRabbitmqServiceName      = "alameda-rabbitmq"
	AlamedaWeaveScopeAppServiceName = "alameda-weave-scope-app"
	FedemeterInfluxdbServiceName    = "fedemeter-influxdb"
	FedemeterAPIServiceName         = "fedemeter-api"
	//CRD Version
	OriAlamedaOperatorVersion = "v0.3.8"
	//AlamedaService modify Prometheus's var
	OriginPrometheus_URL               = "https://prometheus-k8s.openshift-monitoring.svc:9091"
	OriginDeploymentPrometheusLocation = "ALAMEDA_DATAHUB_PROMETHEUS_URL"
	OriginComfigMapPrometheusLocation  = "prometheus.yaml"
	NamespaceService                   = "federatorai.svc"
	DefaultNamespace                   = "federatorai"
	NamespaceServiceAccount            = "serviceaccount:federatorai"
	//MountPath
	DataMountPath = "/var/lib"
	LogMountPath  = "/var/log"
	//Recommandation config
	OriginComfigMapRecommandation = "config.toml"
	//Execution config
	OriginComfigMapExecution = "config.yml"
	//fedemeter config
	OriginComfigMapFedemeterInfluxDBHost = "FEDEMETER_INFLUXDB_HOST"
	OriginComfigMapFedemeterDataHubHost  = "DATAHUB_HOST"
	//Delete Deployment When Modify ConfigMap or Service(Temporary strategy)
	GrafanaYaml            = "Deployment/alameda-grafanaDM.yaml"
	GrafanaDatasourcesName = "grafana-datasources"
)

var (
	//if disable resource protection
	Disable_operand_resource_protection = "false"
	log                                 = logf.Log.WithName("controller_alamedaservice")
	//AlamedaScaler version
	AlamedaScalerVersion        = []string{"v1", "v2"}
	V1scalerOperatorVersionList = []string{"v0.3.6", "v0.3.7", "v0.3.8", "v0.3.9", "v0.3.10", "v0.3.11", "v0.3.12"}
)

// GetServiceAddress returns address combining dns name with port number base on port name
func GetServiceAddress(svc *corev1.Service, portName string) (string, error) {

	portNum := int32(0)
	exist := false
	for _, port := range svc.Spec.Ports {
		if port.Name == portName {
			portNum = port.Port
			exist = true
			break
		}
	}
	if !exist {
		return "", errors.New("port name does not exist")
	}

	namespace := svc.Namespace
	name := svc.Name
	return fmt.Sprintf("%s.%s.svc:%d", name, namespace, portNum), nil
}

// GetServiceDNS returns service dns
func GetServiceDNS(svc *corev1.Service) string {
	namespace := svc.Namespace
	name := svc.Name
	return fmt.Sprintf("%s.%s.svc", name, namespace)
}

func SetBootStrapImageStruct(dep *appsv1.Deployment, componentspec v1alpha1.AlamedaComponentSpec, ctn string) {
	for index, value := range dep.Spec.Template.Spec.InitContainers {
		if value.Name == ctn {
			if componentspec.BootStrapContainer.Image != "" || componentspec.BootStrapContainer.Version != "" {
				image := fmt.Sprintf("%s:%s", componentspec.BootStrapContainer.Image, componentspec.BootStrapContainer.Version)
				dep.Spec.Template.Spec.InitContainers[index].Image = image
			}
			dep.Spec.Template.Spec.InitContainers[index].ImagePullPolicy = componentspec.BootStrapContainer.ImagePullPolicy
		}
	}
}

//if user section schema set image then AlamedaService set Containers image
func setImage(dep *appsv1.Deployment, ctn string, image string) {
	for index, value := range dep.Spec.Template.Spec.Containers {
		if value.Name == ctn {
			newImage := ""
			oriImage := dep.Spec.Template.Spec.Containers[index].Image
			imageStrutct := strings.Split(oriImage, ":")
			if len(imageStrutct) != 0 {
				newImage = fmt.Sprintf("%s:%s", image, imageStrutct[len(imageStrutct)-1])
				dep.Spec.Template.Spec.Containers[index].Image = newImage
			}
		}
	}
}

//if user section schema set image version then AlamedaService set Containers image version
func setImageVersion(dep *appsv1.Deployment, ctn string, version string) {
	for index, value := range dep.Spec.Template.Spec.Containers {
		if value.Name == ctn {
			newImage := ""
			oriImage := dep.Spec.Template.Spec.Containers[index].Image
			imageStrutct := strings.Split(oriImage, ":")
			if len(imageStrutct) != 0 {
				newImage = fmt.Sprintf("%s:%s", strings.Join(imageStrutct[:len(imageStrutct)-1], ":"), version)
				dep.Spec.Template.Spec.Containers[index].Image = newImage
			}
			log.V(1).Info("SetImageVersion", "Deployment.Name", dep.Name, "Container.Name", dep.Spec.Template.Spec.Containers[index].Name, "Image", newImage)
		}
	}
}

func SetStatefulsetImageStruct(ss *appsv1.StatefulSet, value interface{}, ctn string) {
	switch v := value.(type) {
	case string:
		{
			//set global schema image version
			if v != "" {
				setStatefulsetImageVersion(ss, ctn, v)
			}
		}
	}
}

func setStatefulsetImageVersion(ss *appsv1.StatefulSet, ctn string, version string) {
	for index, value := range ss.Spec.Template.Spec.Containers {
		if value.Name == ctn {
			newImage := ""
			oriImage := ss.Spec.Template.Spec.Containers[index].Image
			imageStrutct := strings.Split(oriImage, ":")
			if len(imageStrutct) != 0 {
				newImage = fmt.Sprintf("%s:%s", strings.Join(imageStrutct[:len(imageStrutct)-1], ":"), version)
				ss.Spec.Template.Spec.Containers[index].Image = newImage
			}
			log.V(1).Info("SetImageVersion", ss.Spec.Template.Spec.Containers[index].Name, newImage)
		}
	}
}

//if user set related image struct then AlamedaService set Containers image structure
func SetImageStruct(dep *appsv1.Deployment, value interface{}, ctn string) {
	switch v := value.(type) {
	case string:
		{
			//set global schema image version
			if v != "" {
				setImageVersion(dep, ctn, v)
			}
		}
	case v1alpha1.AlamedaComponentSpec:
		{
			//set section schema image
			if v.Image != "" {
				setImage(dep, ctn, v.Image)
			}
			//set section schema image version
			if v.Version != "" {
				setImageVersion(dep, ctn, v.Version)
			}
		}
	}
}

//if user section schema set pullpolicy then AlamedaService set Containers image's pullpolicy
func SetImagePullPolicy(dep *appsv1.Deployment, ctn string, imagePullPolicy corev1.PullPolicy) {
	for index, value := range dep.Spec.Template.Spec.Containers {
		if value.Name == ctn {
			dep.Spec.Template.Spec.Containers[index].ImagePullPolicy = imagePullPolicy
			log.V(1).Info("SetImagePullPolicy", "Deployment.Name", dep.Name, "Container.Name", dep.Spec.Template.Spec.Containers[index].Name, "ImagePullPolicy", imagePullPolicy)
		}
	}
}

func setDaemonSetImage(ds *appsv1.DaemonSet, ctn string, image string) {
	for index, value := range ds.Spec.Template.Spec.Containers {
		if value.Name == ctn {
			newImage := ""
			oriImage := ds.Spec.Template.Spec.Containers[index].Image
			imageStrutct := strings.Split(oriImage, ":")
			if len(imageStrutct) != 0 {
				newImage = fmt.Sprintf("%s:%s", image, imageStrutct[len(imageStrutct)-1])
				ds.Spec.Template.Spec.Containers[index].Image = newImage
			}
		}
	}
}

func setDaemonSetImageVersion(ds *appsv1.DaemonSet, ctn string, version string) {
	for index, value := range ds.Spec.Template.Spec.Containers {
		if value.Name == ctn {
			newImage := ""
			oriImage := ds.Spec.Template.Spec.Containers[index].Image
			imageStrutct := strings.Split(oriImage, ":")
			if len(imageStrutct) != 0 {
				newImage = fmt.Sprintf("%s:%s", strings.Join(imageStrutct[:len(imageStrutct)-1], ":"), version)
				ds.Spec.Template.Spec.Containers[index].Image = newImage
			}
			log.V(1).Info("SetDaemonSetImageVersion", ds.Spec.Template.Spec.Containers[index].Name, newImage)
		}
	}
}

func SetDaemonSetImageStruct(ds *appsv1.DaemonSet, value interface{}, ctn string) {
	switch v := value.(type) {
	case string:
		{
			//set global schema image version
			if v != "" {
				setDaemonSetImageVersion(ds, ctn, v)
			}
		}
	case v1alpha1.AlamedaComponentSpec:
		{
			//set section schema image
			if v.Image != "" {
				setDaemonSetImage(ds, ctn, v.Image)
			}
			//set section schema image version
			if v.Version != "" {
				setDaemonSetImageVersion(ds, ctn, v.Version)
			}
		}
	}
}

func SetDaemonSetImagePullPolicy(ds *appsv1.DaemonSet, ctn string, imagePullPolicy corev1.PullPolicy) {
	for index, value := range ds.Spec.Template.Spec.Containers {
		if value.Name == ctn {
			ds.Spec.Template.Spec.Containers[index].ImagePullPolicy = imagePullPolicy
			log.V(1).Info("SetDaemonSetImagePullPolicy", ds.Spec.Template.Spec.Containers[index].Name, imagePullPolicy)
		}
	}
}

func SetStatefulSetImagePullPolicy(sts *appsv1.StatefulSet, ctn string, imagePullPolicy corev1.PullPolicy) {
	for index, container := range sts.Spec.Template.Spec.Containers {
		if container.Name == ctn {
			sts.Spec.Template.Spec.Containers[index].ImagePullPolicy = imagePullPolicy
			log.V(1).Info("SetStatefulSetImagePullPolicy", sts.Spec.Template.Spec.Containers[index].Name, imagePullPolicy)
		}
	}
}

//if user set storage log then find VolumeSource path's location
func getVolumeLogIndex(dep *appsv1.Deployment) int {
	if len(dep.Spec.Template.Spec.Volumes) > 0 {
		for index, value := range dep.Spec.Template.Spec.Volumes {
			if value.Name == "alameda-ai-log-storage" {
				return index
			}
			if value.Name == "alameda-operator-log-storage" {
				return index
			}
			if value.Name == "alameda-datahub-log-storage" {
				return index
			}
			if value.Name == "alameda-evictioner-log-storage" {
				return index
			}
			if value.Name == "admission-controller-log-storage" {
				return index
			}
			if value.Name == "alameda-recommender-log-storage" {
				return index
			}
			if value.Name == "alameda-executor-log-storage" {
				return index
			}
			if value.Name == "alameda-dispatcher-log-storage" {
				return index
			}
			if value.Name == "alameda-analyzer-log-storage" {
				return index
			}
			if value.Name == "fedemeter-log-storage" {
				return index
			}
		}
		return -1
	}
	return -1
}

//if user set storage data then find VolumeSource path's location
func getVolumeDataIndex(dep *appsv1.Deployment) int {
	if len(dep.Spec.Template.Spec.Volumes) > 0 {
		for index, value := range dep.Spec.Template.Spec.Volumes {
			if value.Name == "alameda-ai-data-storage" {
				return index
			}
			if value.Name == "alameda-operator-data-storage" {
				return index
			}
			if value.Name == "alameda-datahub-data-storage" {
				return index
			}
			if value.Name == "alameda-evictioner-data-storage" {
				return index
			}
			if value.Name == "admission-controller-data-storage" {
				return index
			}
			if value.Name == "alameda-recommender-data-storage" {
				return index
			}
			if value.Name == "alameda-executor-data-storage" {
				return index
			}
			if value.Name == "alameda-dispatcher-data-storage" {
				return index
			}
			if value.Name == "alameda-analyzer-data-storage" {
				return index
			}
			if value.Name == "influxdb-data-storage" {
				return index
			}
			if value.Name == "grafana-data-storage" {
				return index
			}
			if value.Name == "fedemeter-data-storage" {
				return index
			}
		}
		return -1
	}
	return -1
}

//if user set ephemeral then AlamedaService set Deployment VolumeSource is EmptyDir
func setEmptyDir(dep *appsv1.Deployment, index int, size string) {
	if size != "" {
		quantity := resource.MustParse(size)
		emptydir := &corev1.EmptyDirVolumeSource{SizeLimit: &quantity}
		vs := corev1.VolumeSource{EmptyDir: emptydir}
		dep.Spec.Template.Spec.Volumes[index].VolumeSource = vs
	} else {
		vs := corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}
		dep.Spec.Template.Spec.Volumes[index].VolumeSource = vs
	}
	log.V(1).Info("SetVolumeSource", "Deployment.Name", dep.Name)
}

//if user set pvc then AlamedaService set Deployment VolumeSource is PersistentVolumeClaim
func setVolumeSource(dep *appsv1.Deployment, index int, claimName string) {
	pvcs := &corev1.PersistentVolumeClaimVolumeSource{ClaimName: claimName}
	vs := corev1.VolumeSource{PersistentVolumeClaim: pvcs}
	dep.Spec.Template.Spec.Volumes[index].VolumeSource = vs
	log.V(1).Info("SetVolumeSource", "Deployment.Name", dep.Name, "PVCs", pvcs)
}

//if user set pvc then AlamedaService set pvc to Deployment's VolumeSource
func SetStorageToVolumeSource(dep *appsv1.Deployment, storagestructs []v1alpha1.StorageSpec, volumeName string, group GroupEnums) {
	for _, v := range storagestructs {
		if !v.StorageIsEmpty() {
			if index := getVolumeLogIndex(dep); index != -1 && v.Usage == v1alpha1.Log {
				setVolumeSource(dep, index, strings.Replace(volumeName, "type", string(v1alpha1.Log), -1))
			}
			if index := getVolumeDataIndex(dep); index != -1 && v.Usage == v1alpha1.Data {
				setVolumeSource(dep, index, strings.Replace(volumeName, "type", string(v1alpha1.Data), -1))
			}
			if v.Usage == v1alpha1.Empty && group == AlamedaGroup {
				if index := getVolumeLogIndex(dep); index != -1 {
					setVolumeSource(dep, index, strings.Replace(volumeName, "type", string(v1alpha1.Log), -1))
				}
				if index := getVolumeDataIndex(dep); index != -1 {
					setVolumeSource(dep, index, strings.Replace(volumeName, "type", string(v1alpha1.Data), -1))
				}
			} else if v.Usage == v1alpha1.Empty && group != AlamedaGroup {
				if index := getVolumeDataIndex(dep); index != -1 {
					setVolumeSource(dep, index, strings.Replace(volumeName, "type", string(v1alpha1.Data), -1))
				}
			}
		}
		if v.Type == v1alpha1.Ephemeral {
			if index := getVolumeLogIndex(dep); index != -1 && v.Usage == v1alpha1.Log {
				setEmptyDir(dep, index, v.Size)
			}
			if index := getVolumeDataIndex(dep); index != -1 && v.Usage == v1alpha1.Data {
				setEmptyDir(dep, index, v.Size)
			}
		}
	}
}

func setMountPath(dep *appsv1.Deployment, volumeName string, mountPath string, ctn string, group GroupEnums) {
	for index, value := range dep.Spec.Template.Spec.Containers {
		if value.Name == ctn {
			for _, v := range dep.Spec.Template.Spec.Containers[index].VolumeMounts {
				if v.Name == volumeName { //if global schema has been set up
					return
				}
			}
			if group == AlamedaGroup {
				vm := corev1.VolumeMount{Name: volumeName, MountPath: mountPath}
				dep.Spec.Template.Spec.Containers[index].VolumeMounts = append([]corev1.VolumeMount{vm}, dep.Spec.Template.Spec.Containers[index].VolumeMounts...)
			} else {
				vm := corev1.VolumeMount{Name: volumeName, MountPath: mountPath, SubPath: string(group)}
				dep.Spec.Template.Spec.Containers[index].VolumeMounts = append([]corev1.VolumeMount{vm}, dep.Spec.Template.Spec.Containers[index].VolumeMounts...)
			}

		}
	}
}

//if user set pvc then AlamedaService set pvc to Deployment's MountPath
func SetStorageToMountPath(dep *appsv1.Deployment, storagestructs []v1alpha1.StorageSpec, ctn string, volumeName string, group GroupEnums) {
	for _, v := range storagestructs {
		if v.Type == v1alpha1.Ephemeral || v.Type == v1alpha1.PVC {
			if v.Usage == v1alpha1.Data {
				setMountPath(dep, strings.Replace(volumeName, "type", string(v1alpha1.Data), -1), fmt.Sprintf("%s/%s", DataMountPath, group), ctn, group)
			} else if v.Usage == v1alpha1.Log {
				setMountPath(dep, strings.Replace(volumeName, "type", string(v1alpha1.Log), -1), fmt.Sprintf("%s/%s", LogMountPath, group), ctn, group)
			} else if v.Usage == v1alpha1.Empty && group == AlamedaGroup {
				setMountPath(dep, strings.Replace(volumeName, "type", string(v1alpha1.Data), -1), fmt.Sprintf("%s/%s", DataMountPath, group), ctn, group)
				setMountPath(dep, strings.Replace(volumeName, "type", string(v1alpha1.Log), -1), fmt.Sprintf("%s/%s", LogMountPath, group), ctn, group)
			} else if v.Usage == v1alpha1.Empty && group != AlamedaGroup { // if not alameda component's then only set data
				setMountPath(dep, strings.Replace(volumeName, "type", string(v1alpha1.Data), -1), fmt.Sprintf("%s/%s", DataMountPath, group), ctn, group)
			}
		}
	}
}

func setPVCSpec(pvc *corev1.PersistentVolumeClaim, value v1alpha1.StorageSpec) {
	if value.AccessModes != "" {
		pvc.Spec.AccessModes = append(pvc.Spec.AccessModes, value.AccessModes)
	}
	if value.Size != "" {
		pvc.Spec.Resources.Requests[corev1.ResourceStorage] = resource.MustParse(value.Size)
	}
	if value.Class != nil {
		pvc.Spec.StorageClassName = value.Class
	}
}

//if user set pvc then AlamedaService set PersistentVolumeClaimSpec
func SetStorageToPersistentVolumeClaimSpec(pvc *corev1.PersistentVolumeClaim, storagestructs []v1alpha1.StorageSpec, pvctype v1alpha1.Usage) {
	for k, v := range storagestructs {
		if v.Usage == pvctype || v.Usage == v1alpha1.Empty {
			setPVCSpec(pvc, storagestructs[k])
		}
	}
}

func StringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

// StringSliceDelete removes elements in slice2 from slice1
func StringSliceDelete(slice1 []string, slice2 []string) []string {

	var diff []string

	for _, s1 := range slice1 {
		found := false
		for _, s2 := range slice2 {
			if s1 == s2 {
				found = true
				break
			}
		}
		// String not found. We add it to return slice
		if !found {
			diff = append(diff, s1)
		}
	}

	return diff
}

// GetClusterUID returns metadata.uid from "cluster-info" configMap.
func GetClusterUID(k8sClient client.Client) (string, error) {
	possibleNSList := []string{
		"kube-service-catalog", "kube-public",
	}
	var errorList []string
	clusterID := ""
	for _, possibleNS := range possibleNSList {
		clusterInfoCM := &corev1.ConfigMap{}
		err := k8sClient.Get(context.Background(), client.ObjectKey{
			Name:      "cluster-info",
			Namespace: possibleNS,
		}, clusterInfoCM)
		if err == nil {
			return string(clusterInfoCM.GetUID()), nil
		} else if !k8sErrors.IsNotFound(err) {
			errorList = append(errorList, err.Error())
		}
	}
	if len(errorList) == 0 {
		return clusterID, errors.New("no cluster info found")
	}
	return clusterID, errors.New(strings.Join(errorList, ","))
}
