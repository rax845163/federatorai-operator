package alamedaserviceparamter

import (
	"strings"

	admission_controller "github.com/containers-ai/alameda/admission-controller"
	"github.com/containers-ai/federatorai-operator/pkg/apis/federatorai/v1alpha1"
	"github.com/containers-ai/federatorai-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
)

var (
	crbList = []string{"ClusterRoleBinding/alameda-datahubCRB.yaml",
		"ClusterRoleBinding/alameda-operatorCRB.yaml",
	}
	crList = []string{"ClusterRole/alameda-datahubCR.yaml",
		"ClusterRole/alameda-operatorCR.yaml",
		"ClusterRole/aggregate-alameda-admin-edit-alamedaCR.yaml",
	}

	saList = []string{"ServiceAccount/alameda-datahubSA.yaml",
		"ServiceAccount/alameda-operatorSA.yaml",
		"ServiceAccount/alameda-aiSA.yaml",
	}
	crdList = []string{
		"CustomResourceDefinition/alamedarecommendationsCRD.yaml",
	}
	cmList = []string{
		"ConfigMap/alameda-recommender-config.yaml",
	}

	svList = []string{"Service/alameda-datahubSV.yaml",
		"Service/alameda-influxdbSV.yaml",
		"Service/alameda-ai-metricsSV.yaml"}

	depList = []string{"Deployment/alameda-datahubDM.yaml",
		"Deployment/alameda-operatorDM.yaml",
		"Deployment/alameda-influxdbDM.yaml",
		"Deployment/alameda-aiDM.yaml",
		"Deployment/alameda-recommenderDM.yaml",
	}

	guiList = []string{
		"ClusterRoleBinding/alameda-grafanaCRB.yaml",
		"ClusterRole/alameda-grafanaCR.yaml",
		"ServiceAccount/alameda-grafanaSA.yaml",
		"ConfigMap/grafana-datasources.yaml",
		"ConfigMap/dashboards-config.yaml",
		"Deployment/alameda-grafanaDM.yaml",
		"Service/alameda-grafanaSV.yaml",
	}
	excutionList = []string{
		"ClusterRoleBinding/alameda-evictionerCRB.yaml",
		"ClusterRoleBinding/admission-controllerCRB.yaml",
		"ClusterRole/alameda-evictionerCR.yaml",
		"ClusterRole/admission-controllerCR.yaml",
		"ServiceAccount/alameda-evictionerSA.yaml",
		"ServiceAccount/admission-controllerSA.yaml",
		"Secret/admission-controller-tls.yaml",
		"Deployment/admission-controllerDM.yaml",
		"Deployment/alameda-evictionerDM.yaml",
		"Service/admission-controllerSV.yaml",
		"Deployment/alameda-executorDM.yaml",
		"ServiceAccount/alameda-executorSA.yaml",
		"ClusterRole/alameda-executorCR.yaml",
		"ClusterRoleBinding/alameda-executorCRB.yaml",
	}

	secretList = []string{
		"Secret/alameda-influxdb.yaml",
	}
	pvcList = []string{
		"PersistentVolumeClaim/my-alamedainfluxdbPVC.yaml",
		"PersistentVolumeClaim/my-alamedagrafanaPVC.yaml",
		"PersistentVolumeClaim/alameda-ai-log.yaml",
		"PersistentVolumeClaim/alameda-operator-log.yaml",
		"PersistentVolumeClaim/alameda-datahub-log.yaml",
		"PersistentVolumeClaim/alameda-evictioner-log.yaml",
		"PersistentVolumeClaim/admission-controller-log.yaml",
		"PersistentVolumeClaim/alameda-ai-data.yaml",
		"PersistentVolumeClaim/alameda-operator-data.yaml",
		"PersistentVolumeClaim/alameda-datahub-data.yaml",
		"PersistentVolumeClaim/alameda-evictioner-data.yaml",
		"PersistentVolumeClaim/admission-controller-data.yaml",
	}
)

type AlamedaServiceParamter struct {
	NameSpace                     string
	SelfDriving                   bool
	Platform                      string
	EnableExecution               bool
	EnableGUI                     bool
	Version                       string
	PrometheusService             string
	Storages                      []v1alpha1.StorageSpec
	InfluxdbSectionSet            v1alpha1.AlamedaComponentSpec
	GrafanaSectionSet             v1alpha1.AlamedaComponentSpec
	AlamedaAISectionSet           v1alpha1.AlamedaComponentSpec
	AlamedaOperatorSectionSet     v1alpha1.AlamedaComponentSpec
	AlamedaDatahubSectionSet      v1alpha1.AlamedaComponentSpec
	AlamedaEvictionerSectionSet   v1alpha1.AlamedaComponentSpec
	AdmissionControllerSectionSet v1alpha1.AlamedaComponentSpec
	AlamedaRecommenderSectionSet  v1alpha1.AlamedaComponentSpec
	AlamedaExecutorSectionSet     v1alpha1.AlamedaComponentSpec
	CurrentCRDVersion             v1alpha1.AlamedaServiceStatusCRDVersion
	previousCRDVersion            v1alpha1.AlamedaServiceStatusCRDVersion
}

type Resource struct {
	ClusterRoleBindingList       []string
	ClusterRoleList              []string
	ServiceAccountList           []string
	CustomResourceDefinitionList []string
	ConfigMapList                []string
	ServiceList                  []string
	DeploymentList               []string
	SecretList                   []string
	PersistentVolumeClaimList    []string
	AlamdaScalerList             []string
}

func GetSelfDrivingRsource() *Resource {
	var alamedaScalerList = make([]string, 0)
	alamedaScalerList = append(alamedaScalerList, "AlamedaScaler/alamedaScaler-alameda.yaml")
	return &Resource{
		AlamdaScalerList: alamedaScalerList,
	}
}

func GetExcutionResource() *Resource {
	var excrb = make([]string, 0)
	var excr = make([]string, 0)
	var exsa = make([]string, 0)
	var excsec = make([]string, 0)
	var excDep = make([]string, 0)
	var excCM = make([]string, 0)
	var excSV = make([]string, 0)
	for _, str := range excutionList {
		if len(strings.Split(str, "/")) > 0 {
			switch resource := strings.Split(str, "/")[0]; resource {
			case "ClusterRoleBinding":
				excrb = append(excrb, str)
			case "ClusterRole":
				excr = append(excr, str)
			case "Secret":
				excsec = append(excsec, str)
			case "ServiceAccount":
				exsa = append(exsa, str)
			case "ConfigMap":
				excCM = append(excCM, str)
			case "Service":
				excSV = append(excSV, str)
			case "Deployment":
				excDep = append(excDep, str)
			default:
			}
		}
	}
	return &Resource{
		ClusterRoleBindingList: excrb,
		ClusterRoleList:        excr,
		ServiceAccountList:     exsa,
		SecretList:             excsec,
		ConfigMapList:          excCM,
		ServiceList:            excSV,
		DeploymentList:         excDep,
	}
}

func sectionUninstallPersistentVolumeClaimSource(pvc []string, storagestruct []v1alpha1.StorageSpec, resourceName string, resourceType v1alpha1.Usage) []string {
	for _, value := range storagestruct {
		if value.Type != v1alpha1.PVC {
			if value.Usage == resourceType || value.Usage == v1alpha1.Empty {
				pvc = append(pvc, resourceName)
			}
		} else { //component section set pvc
			if value.Usage == resourceType || value.Usage == v1alpha1.Empty {
				for k, v := range pvc {
					if v == resourceName {
						pvc = append(pvc[:k], pvc[k+1:]...)
					}
				}
			}
		}
	}
	return pvc
}

func (asp *AlamedaServiceParamter) GetUninstallPersistentVolumeClaimSource() *Resource {
	pvc := []string{}
	for _, v := range asp.Storages {
		if v.Type != v1alpha1.PVC {
			if v.Usage == v1alpha1.Log {
				pvc = append(pvc, "PersistentVolumeClaim/alameda-ai-log.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/alameda-operator-log.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/alameda-datahub-log.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/alameda-evictioner-log.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/admission-controller-log.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/alameda-recommender-log.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/alameda-executor-log.yaml")
			} else if v.Usage == v1alpha1.Data {
				pvc = append(pvc, "PersistentVolumeClaim/my-alamedainfluxdbPVC.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/my-alamedagrafanaPVC.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/alameda-ai-data.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/alameda-operator-data.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/alameda-datahub-data.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/alameda-evictioner-data.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/admission-controller-data.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/alameda-recommender-data.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/alameda-executor-data.yaml")
			} else if v.Usage == v1alpha1.Empty {
				pvc = append(pvc, "PersistentVolumeClaim/alameda-ai-log.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/alameda-operator-log.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/alameda-datahub-log.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/alameda-evictioner-log.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/admission-controller-log.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/alameda-recommender-log.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/alameda-executor-log.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/my-alamedainfluxdbPVC.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/my-alamedagrafanaPVC.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/alameda-ai-data.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/alameda-operator-data.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/alameda-datahub-data.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/alameda-evictioner-data.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/admission-controller-data.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/alameda-recommender-data.yaml")
				pvc = append(pvc, "PersistentVolumeClaim/alameda-executor-data.yaml")
			}
		}
	}
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaAISectionSet.Storages, "PersistentVolumeClaim/alameda-ai-log.yaml", v1alpha1.Log)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaOperatorSectionSet.Storages, "PersistentVolumeClaim/alameda-operator-log.yaml", v1alpha1.Log)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaDatahubSectionSet.Storages, "PersistentVolumeClaim/alameda-datahub-log.yaml", v1alpha1.Log)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaEvictionerSectionSet.Storages, "PersistentVolumeClaim/alameda-evictioner-log.yaml", v1alpha1.Log)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AdmissionControllerSectionSet.Storages, "PersistentVolumeClaim/admission-controller-log.yaml", v1alpha1.Log)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaRecommenderSectionSet.Storages, "PersistentVolumeClaim/alameda-recommender-log.yaml", v1alpha1.Log)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaRecommenderSectionSet.Storages, "PersistentVolumeClaim/alameda-executor-log.yaml", v1alpha1.Log)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.InfluxdbSectionSet.Storages, "PersistentVolumeClaim/my-alamedainfluxdbPVC.yaml", v1alpha1.Data)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.GrafanaSectionSet.Storages, "PersistentVolumeClaim/my-alamedagrafanaPVC.yaml", v1alpha1.Data)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaAISectionSet.Storages, "PersistentVolumeClaim/alameda-ai-data.yaml", v1alpha1.Data)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaOperatorSectionSet.Storages, "PersistentVolumeClaim/alameda-operator-data.yaml", v1alpha1.Data)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaDatahubSectionSet.Storages, "PersistentVolumeClaim/alameda-datahub-data.yaml", v1alpha1.Data)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaEvictionerSectionSet.Storages, "PersistentVolumeClaim/alameda-evictioner-data.yaml", v1alpha1.Data)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AdmissionControllerSectionSet.Storages, "PersistentVolumeClaim/admission-controller-data.yaml", v1alpha1.Data)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaRecommenderSectionSet.Storages, "PersistentVolumeClaim/alameda-recommender-data.yaml", v1alpha1.Log)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaRecommenderSectionSet.Storages, "PersistentVolumeClaim/alameda-executor-data.yaml", v1alpha1.Log)
	return &Resource{
		PersistentVolumeClaimList: pvc,
	}

}

func (asp *AlamedaServiceParamter) GetEnvVarsByDeployment(deploymentName string) []corev1.EnvVar {

	var envVars []corev1.EnvVar

	switch deploymentName {
	case util.AdmissioncontrollerDPN:
		envVars = asp.GetAdmissionControllerEnvVars()
	case util.AlamedaevictionerDPN:
		envVars = asp.GetAlamedaEvictionerEnvVars()
	default:
	}

	return envVars
}

func (asp *AlamedaServiceParamter) GetAdmissionControllerEnvVars() []corev1.EnvVar {

	envVars := make([]corev1.EnvVar, 0)

	switch asp.Platform {
	case v1alpha1.PlatformOpenshift3_9:
		envVars = append(envVars, corev1.EnvVar{
			Name:  "ALAMEDA_ADMCTL_JSON_PATCH_VALIDATION_FUNC",
			Value: admission_controller.JsonPatchValidationFuncOpenshift3_9,
		})
	}

	return envVars
}

func (asp *AlamedaServiceParamter) GetAlamedaEvictionerEnvVars() []corev1.EnvVar {
	envVars := make([]corev1.EnvVar, 0)

	switch asp.Platform {
	case v1alpha1.PlatformOpenshift3_9:
		envVars = append(envVars, corev1.EnvVar{
			Name:  "ALAMEDA_EVICTIONER_EVICTION_PURGE_CONTAINER_CPU_MEMORY",
			Value: "true",
		})
	}

	return envVars
}

func GetGUIResource() *Resource {
	var guicrb = make([]string, 0)
	var guicr = make([]string, 0)
	var guisa = make([]string, 0)
	var guiDep = make([]string, 0)
	var guiCM = make([]string, 0)
	var guiSV = make([]string, 0)
	for _, str := range guiList {
		if len(strings.Split(str, "/")) > 0 {
			switch resource := strings.Split(str, "/")[0]; resource {
			case "ClusterRoleBinding":
				guicrb = append(guicrb, str)
			case "ClusterRole":
				guicr = append(guicr, str)
			case "ServiceAccount":
				guisa = append(guisa, str)
			case "ConfigMap":
				guiCM = append(guiCM, str)
			case "Service":
				guiSV = append(guiSV, str)
			case "Deployment":
				guiDep = append(guiDep, str)
			default:
			}
		}
	}
	return &Resource{
		ClusterRoleBindingList: guicrb,
		ClusterRoleList:        guicr,
		ServiceAccountList:     guisa,
		ConfigMapList:          guiCM,
		ServiceList:            guiSV,
		DeploymentList:         guiDep,
	}
}

func GetUnInstallResource() *Resource {
	return &Resource{
		ClusterRoleBindingList:       crbList,
		ClusterRoleList:              crList,
		ServiceAccountList:           saList,
		CustomResourceDefinitionList: crdList,
		ConfigMapList:                cmList,
		ServiceList:                  svList,
		DeploymentList:               depList,
		SecretList:                   secretList,
	}
}

func sectioninstallPersistentVolumeClaimSource(pvc []string, storagestruct []v1alpha1.StorageSpec, resourceName string, resourceType v1alpha1.Usage) []string {
	for _, value := range storagestruct {
		if value.Type == v1alpha1.PVC {
			if value.Usage == resourceType || value.Usage == v1alpha1.Empty {
				pvc = append(pvc, resourceName)
			}
		} else if value.Type != v1alpha1.PVC {
			if value.Usage == resourceType || value.Usage == v1alpha1.Empty {
				for k, v := range pvc {
					if v == resourceName {
						pvc = append(pvc[:k], pvc[k+1:]...)
					}
				}
			}
		}
	}
	return pvc
}

func (asp *AlamedaServiceParamter) getInstallPersistentVolumeClaimSource(pvc []string) []string {
	// get install resource
	gloabalLogFlag := false
	gloabalDataFlag := false
	for _, value := range asp.Storages {
		if (value.Usage == v1alpha1.Log || value.Usage == v1alpha1.Empty) && value.Type == v1alpha1.PVC { //Gloabal append
			gloabalLogFlag = !gloabalLogFlag
		}
		if (value.Usage == v1alpha1.Data || value.Usage == v1alpha1.Empty) && value.Type == v1alpha1.PVC {
			gloabalDataFlag = !gloabalDataFlag
		}
	}
	if gloabalLogFlag { //Gloabal append
		pvc = append(pvc, "PersistentVolumeClaim/alameda-ai-log.yaml")
		pvc = append(pvc, "PersistentVolumeClaim/alameda-operator-log.yaml")
		pvc = append(pvc, "PersistentVolumeClaim/alameda-datahub-log.yaml")
		pvc = append(pvc, "PersistentVolumeClaim/alameda-evictioner-log.yaml")
		pvc = append(pvc, "PersistentVolumeClaim/admission-controller-log.yaml")
		pvc = append(pvc, "PersistentVolumeClaim/alameda-recommender-log.yaml")
		pvc = append(pvc, "PersistentVolumeClaim/alameda-executor-log.yaml")
	}
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AlamedaAISectionSet.Storages, "PersistentVolumeClaim/alameda-ai-log.yaml", v1alpha1.Log)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AlamedaOperatorSectionSet.Storages, "PersistentVolumeClaim/alameda-operator-log.yaml", v1alpha1.Log)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AlamedaDatahubSectionSet.Storages, "PersistentVolumeClaim/alameda-datahub-log.yaml", v1alpha1.Log)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AlamedaEvictionerSectionSet.Storages, "PersistentVolumeClaim/alameda-evictioner-log.yaml", v1alpha1.Log)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AdmissionControllerSectionSet.Storages, "PersistentVolumeClaim/admission-controller-log.yaml", v1alpha1.Log)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AdmissionControllerSectionSet.Storages, "PersistentVolumeClaim/alameda-recommender-log.yaml", v1alpha1.Log)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AdmissionControllerSectionSet.Storages, "PersistentVolumeClaim/alameda-executor-log.yaml", v1alpha1.Log)
	if gloabalDataFlag {
		pvc = append(pvc, "PersistentVolumeClaim/my-alamedainfluxdbPVC.yaml")
		pvc = append(pvc, "PersistentVolumeClaim/my-alamedagrafanaPVC.yaml")
		pvc = append(pvc, "PersistentVolumeClaim/alameda-ai-data.yaml")
		pvc = append(pvc, "PersistentVolumeClaim/alameda-operator-data.yaml")
		pvc = append(pvc, "PersistentVolumeClaim/alameda-datahub-data.yaml")
		pvc = append(pvc, "PersistentVolumeClaim/alameda-evictioner-data.yaml")
		pvc = append(pvc, "PersistentVolumeClaim/admission-controller-data.yaml")
		pvc = append(pvc, "PersistentVolumeClaim/alameda-recommender-data.yaml")
		pvc = append(pvc, "PersistentVolumeClaim/alameda-executor-data.yaml")
	}
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.InfluxdbSectionSet.Storages, "PersistentVolumeClaim/my-alamedainfluxdbPVC.yaml", v1alpha1.Data)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.GrafanaSectionSet.Storages, "PersistentVolumeClaim/my-alamedagrafanaPVC.yaml", v1alpha1.Data)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AlamedaAISectionSet.Storages, "PersistentVolumeClaim/alameda-ai-data.yaml", v1alpha1.Data)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AlamedaOperatorSectionSet.Storages, "PersistentVolumeClaim/alameda-operator-data.yaml", v1alpha1.Data)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AlamedaDatahubSectionSet.Storages, "PersistentVolumeClaim/alameda-datahub-data.yaml", v1alpha1.Data)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AlamedaEvictionerSectionSet.Storages, "PersistentVolumeClaim/alameda-evictioner-data.yaml", v1alpha1.Data)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AdmissionControllerSectionSet.Storages, "PersistentVolumeClaim/admission-controller-data.yaml", v1alpha1.Data)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AdmissionControllerSectionSet.Storages, "PersistentVolumeClaim/alameda-recommender-data.yaml", v1alpha1.Log)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AdmissionControllerSectionSet.Storages, "PersistentVolumeClaim/alameda-executor-data.yaml", v1alpha1.Log)
	return pvc

}

func (asp *AlamedaServiceParamter) changeScalerCRDVersion(crd []string) []string {
	alamedaOperatorVersion := util.OriAlamedaOperatorVersion
	if asp.Version != "" {
		alamedaOperatorVersion = asp.Version
	}
	if asp.AlamedaOperatorSectionSet.Version != "" {
		alamedaOperatorVersion = asp.AlamedaOperatorSectionSet.Version
	}
	if util.StringInSlice(alamedaOperatorVersion, util.V1scalerOperatorVersionList) { //check current operatorVersion used scaler version is scaler V1
		crd = append(crd, "CustomResourceDefinition/alamedascalersCRD.yaml")
		asp.CurrentCRDVersion.ScalerVersion = util.AlamedaScalerVersion[0]
		asp.CurrentCRDVersion.CRDName = util.AlamedaScalerName
	} else {
		crd = append(crd, "CustomResourceDefinition/alamedascalersV2CRD.yaml")
		asp.CurrentCRDVersion.ScalerVersion = util.AlamedaScalerVersion[1]
		asp.CurrentCRDVersion.CRDName = util.AlamedaScalerName
	}
	if asp.CurrentCRDVersion.ScalerVersion != asp.previousCRDVersion.ScalerVersion {
		asp.SetCurrentCRDChangeVersionToTrue()
	}
	return crd
}

func (asp *AlamedaServiceParamter) CheckCurrentCRDIsChangeVersion() bool {
	return asp.CurrentCRDVersion.ChangeVersion
}

func (asp *AlamedaServiceParamter) SetCurrentCRDChangeVersionToFalse() {
	asp.CurrentCRDVersion.ChangeVersion = false
}

func (asp *AlamedaServiceParamter) SetCurrentCRDChangeVersionToTrue() {
	asp.CurrentCRDVersion.ChangeVersion = true
}

func (asp *AlamedaServiceParamter) GetInstallResource() *Resource {
	crb := crbList
	cr := crList
	sa := saList
	crd := crdList
	cm := cmList
	sv := svList
	dep := depList
	secrets := secretList
	pvc := []string{}
	alamdaScalerList := []string{}
	if asp.SelfDriving {
		alamdaScalerList = append(alamdaScalerList, "AlamedaScaler/alamedaScaler-alameda.yaml")
	}
	if asp.EnableGUI {
		crb = append(crb, "ClusterRoleBinding/alameda-grafanaCRB.yaml")
		cr = append(cr, "ClusterRole/alameda-grafanaCR.yaml")
		sa = append(sa, "ServiceAccount/alameda-grafanaSA.yaml")
		cm = append(cm, "ConfigMap/grafana-datasources.yaml")
		cm = append(cm, "ConfigMap/dashboards-config.yaml")
		sv = append(sv, "Service/alameda-grafanaSV.yaml")
		dep = append(dep, "Deployment/alameda-grafanaDM.yaml")
	}
	if asp.EnableExecution {
		crb = append(crb, "ClusterRoleBinding/alameda-evictionerCRB.yaml")
		crb = append(crb, "ClusterRoleBinding/admission-controllerCRB.yaml")
		crb = append(crb, "ClusterRoleBinding/alameda-executorCRB.yaml")
		cr = append(cr, "ClusterRole/alameda-evictionerCR.yaml")
		cr = append(cr, "ClusterRole/admission-controllerCR.yaml")
		cr = append(cr, "ClusterRole/alameda-executorCR.yaml")
		secrets = append(secrets, "Secret/admission-controller-tls.yaml")
		sa = append(sa, "ServiceAccount/alameda-evictionerSA.yaml")
		sa = append(sa, "ServiceAccount/admission-controllerSA.yaml")
		sa = append(sa, "ServiceAccount/alameda-executorSA.yaml")
		cm = append(cm, "ConfigMap/alameda-executor-config.yaml")
		sv = append(sv, "Service/admission-controllerSV.yaml")
		dep = append(dep, "Deployment/admission-controllerDM.yaml")
		dep = append(dep, "Deployment/alameda-evictionerDM.yaml")
		dep = append(dep, "Deployment/alameda-executorDM.yaml")
	}
	pvc = asp.getInstallPersistentVolumeClaimSource(pvc)
	crd = asp.changeScalerCRDVersion(crd)
	return &Resource{
		ClusterRoleBindingList:       crb,
		ClusterRoleList:              cr,
		ServiceAccountList:           sa,
		CustomResourceDefinitionList: crd,
		ConfigMapList:                cm,
		ServiceList:                  sv,
		DeploymentList:               dep,
		SecretList:                   secrets,
		PersistentVolumeClaimList:    pvc,
		AlamdaScalerList:             alamdaScalerList,
	}
}

func NewAlamedaServiceParamter(instance *v1alpha1.AlamedaService) *AlamedaServiceParamter {
	asp := &AlamedaServiceParamter{
		NameSpace:                     instance.Namespace,
		SelfDriving:                   instance.Spec.SelfDriving,
		Platform:                      instance.Spec.Platform,
		EnableExecution:               instance.Spec.EnableExecution,
		EnableGUI:                     instance.Spec.EnableGUI,
		Version:                       instance.Spec.Version,
		PrometheusService:             instance.Spec.PrometheusService,
		Storages:                      instance.Spec.Storages,
		InfluxdbSectionSet:            instance.Spec.InfluxdbSectionSet,
		GrafanaSectionSet:             instance.Spec.GrafanaSectionSet,
		AlamedaAISectionSet:           instance.Spec.AlamedaAISectionSet,
		AlamedaOperatorSectionSet:     instance.Spec.AlamedaOperatorSectionSet,
		AlamedaDatahubSectionSet:      instance.Spec.AlamedaDatahubSectionSet,
		AlamedaEvictionerSectionSet:   instance.Spec.AlamedaEvictionerSectionSet,
		AdmissionControllerSectionSet: instance.Spec.AdmissionControllerSectionSet,
		AlamedaRecommenderSectionSet:  instance.Spec.AlamedaRecommenderSectionSet,
		AlamedaExecutorSectionSet:     instance.Spec.AlamedaExecutorSectionSet,
		CurrentCRDVersion:             instance.Status.CRDVersion,
		previousCRDVersion:            instance.Status.CRDVersion,
	}
	return asp
}
