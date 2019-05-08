package util

import (
	"fmt"
	"strings"

	"github.com/containers-ai/federatorai-operator/pkg/apis/federatorai/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

type GroupEnums string

const (
	AlamedaGroup  GroupEnums = "alameda"
	GrafanaGroup  GroupEnums = "grafana"
	InfluxDBGroup GroupEnums = "influxdb"
	//deployment name
	AlamedaaiDPN           = "alameda-ai"
	AlamedaoperatorDPN     = "alameda-operator"
	AlamedadatahubDPN      = "alameda-datahub"
	AlamedaevictionerDPN   = "alameda-evictioner"
	AdmissioncontrollerDPN = "admission-controller"
	GrafanaDPN             = "alameda-grafana"
	InfluxdbDPN            = "alameda-influxdb"
	//container name
	AlamedaaiCTN           = "alameda-ai-engine"
	AlamedaoperatorCTN     = "alameda-operator"
	AlamedadatahubCTN      = "alameda-datahub"
	AlamedaevictionerCTN   = "alameda-evictioner"
	AdmissioncontrollerCTN = "admission-controller"
	GetTokenCTN            = "gettoken"
	GrafanaCTN             = "grafana"
	InfluxdbCTN            = "influxdb"
	//AlamedaService modify Prometheus's var
	OriginPrometheus_URL               = "https://prometheus-k8s.openshift-monitoring.svc:9091"
	OriginDeploymentPrometheusLocation = "ALAMEDA_DATAHUB_PROMETHEUS_URL"
	OriginComfigMapPrometheusLocation  = "prometheus.yaml"
	NamespaceService                   = "federatorai.svc"
	//MountPath
	DataMountPath = "/var/lib"
	LogMountPath  = "/var/log"
)

var (
	log = logf.Log.WithName("controller_alamedaservice")
)

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
			log.V(1).Info("SetImageVersion", dep.Spec.Template.Spec.Containers[index].Name, newImage)
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
			log.V(1).Info("SetImagePullPolicy", dep.Spec.Template.Spec.Containers[index].Name, imagePullPolicy)
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
			if value.Name == "influxdb-data-storage" {
				return index
			}
			if value.Name == "grafana-data-storage" {
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
	log.V(1).Info("SetVolumeSource", dep.Name)
}

//if user set pvc then AlamedaService set Deployment VolumeSource is PersistentVolumeClaim
func setVolumeSource(dep *appsv1.Deployment, index int, claimName string) {
	pvcs := &corev1.PersistentVolumeClaimVolumeSource{ClaimName: claimName}
	vs := corev1.VolumeSource{PersistentVolumeClaim: pvcs}
	dep.Spec.Template.Spec.Volumes[index].VolumeSource = vs
	log.V(1).Info("SetVolumeSource", dep.Name, pvcs)
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
