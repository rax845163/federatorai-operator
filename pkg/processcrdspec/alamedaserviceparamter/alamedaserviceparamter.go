package alamedaserviceparamter

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/containers-ai/federatorai-operator/pkg/apis/federatorai/v1alpha1"
	"github.com/containers-ai/federatorai-operator/pkg/util"
)

var (
	defaultInstallLists = [][]string{
		[]string{
			"ClusterRole/aggregate-alameda-admin-edit-alamedaCR.yaml",
		},
		customResourceDefinitionList,
		datahubList,
		operatorList,
		influxDBList,
		aiEngineList,
		recommenderList,
		analyzerList,
		rabbitmqList,
		notifierList,
		federatoraiAgentList,
		fedemeterList,
	}

	datahubList = []string{
		"ClusterRoleBinding/alameda-datahubCRB.yaml",
		"ClusterRole/alameda-datahubCR.yaml",
		"ServiceAccount/alameda-datahubSA.yaml",
		"Service/alameda-datahubSV.yaml",
		"Deployment/alameda-datahubDM.yaml",
		"Role/alameda-datahub.yaml",
		"RoleBinding/alameda-datahub.yaml",
	}

	operatorList = []string{
		"ClusterRoleBinding/alameda-operatorCRB.yaml",
		"ClusterRole/alameda-operatorCR.yaml",
		"ServiceAccount/alameda-operatorSA.yaml",
		"Deployment/alameda-operatorDM.yaml",
	}

	influxDBList = []string{
		"Service/alameda-influxdbSV.yaml",
		"Deployment/alameda-influxdbDM.yaml",
		"Secret/alameda-influxdb.yaml",
	}

	aiEngineList = []string{
		"ServiceAccount/alameda-aiSA.yaml",
		"Service/alameda-ai-metricsSV.yaml",
		"Deployment/alameda-aiDM.yaml",
	}

	rabbitmqList = []string{
		"Deployment/alameda-rabbitmqDM.yaml",
		"Service/alameda-rabbitmqSV.yaml",
		"ServiceAccount/alameda-rabbitmqSA.yaml",
		"ClusterRole/alameda-rabbitmqCR.yaml",
		"ClusterRoleBinding/alameda-rabbitmqCRB.yaml",
	}

	recommenderList = []string{
		"ConfigMap/alameda-recommender-config.yaml",
		"Deployment/alameda-recommenderDM.yaml",
	}

	analyzerList = []string{
		"ClusterRoleBinding/alameda-analyzerCRB.yaml",
		"ClusterRole/alameda-analyzerCR.yaml",
		"ServiceAccount/alameda-analyzerSA.yaml",
		"Deployment/alameda-analyzerDM.yaml",
	}

	guiList = []string{
		"ClusterRoleBinding/alameda-grafanaCRB.yaml",
		"ClusterRole/alameda-grafanaCR.yaml",
		"ServiceAccount/alameda-grafanaSA.yaml",
		"ConfigMap/grafana-datasources.yaml",
		"Deployment/alameda-grafanaDM.yaml",
		"Service/alameda-grafanaSV.yaml",
		"Route/alameda-grafanaRT.yaml",
	}

	excutionList = []string{
		"ClusterRole/alameda-vpa-executorCR.yaml",
		"ClusterRoleBinding/alameda-vpa-executorCRB.yaml",
		"Deployment/alameda-vpa-executorDM.yaml",
		"MutatingWebhookConfiguration/alameda-vpa-executor-mutating-webhook-configuration.yaml",
		"Secret/alameda-vpa-executor-server-cert.yaml",
		"Service/alameda-vpa-executorSV.yaml",
		"ServiceAccount/alameda-vpa-executorSA.yaml",
		"Deployment/alameda-executorDM.yaml",
		"ServiceAccount/alameda-executorSA.yaml",
		"ClusterRole/alameda-executorCR.yaml",
		"ClusterRoleBinding/alameda-executorCRB.yaml",
		"ConfigMap/alameda-executor-config.yaml",
	}

	fedemeterList = []string{
		"Deployment/fedemeterDM.yaml",
		"Service/fedemeterSV.yaml",
		"ConfigMap/fedemeter-config.yaml",
		"Service/fedemeter-influxdbSV.yaml",
		"StatefulSet/fedemeter-influxdbSS.yaml",
		"Ingress/fedemeterIG.yaml",
		"Ingress/fedemeterSwaggerIG.yaml",
		"Secret/fedemeter-tls.yaml",
	}

	aiDispatcherList = []string{
		"Deployment/alameda-ai-dispatcherDM.yaml",
	}

	weavescopeList = []string{
		"ClusterRole/alameda-weavescopeCR.yaml",
		"ClusterRoleBinding/alameda-weavescopeCRB.yaml",
		"DaemonSet/alamdea-weavescopeDS.yaml",
		"Deployment/alameda-weavescope-probeDM.yaml",
		"Deployment/alameda-weavescopeDM.yaml",
		"PodSecurityPolicy/alameda-weavescopePSP.yaml",
		"Service/alameda-weavescopeSV.yaml",
		"ServiceAccount/alameda-weavescopeSA.yaml",
	}

	notifierList = []string{
		"ClusterRole/alameda-notifier.yaml",
		"ClusterRoleBinding/alameda-notifier.yaml",
		"Deployment/alameda-notifier.yaml",
		"MutatingWebhookConfiguration/alameda-notifier-mutating-webhook-configuration.yaml",
		"Service/alameda-notifier-webhook-service.yaml",
		"ServiceAccount/alameda-notifier.yaml",
		"ValidatingWebhookConfiguration/alameda-notifier-validating-webhook-configuration.yaml",
		"AlamedaNotificationChannel/default.yaml",
		"AlamedaNotificationTopic/default.yaml",
		"Secret/alameda-notifier-webhook-server-cert.yaml",
	}

	federatoraiAgentList = []string{
		"Deployment/federatorai-agent.yaml",
		"ConfigMap/federatorai-agent-config.yaml",
	}

	selfDrivingList = []string{
		"AlamedaScaler/alamedaScaler-alameda.yaml",
	}

	customResourceDefinitionList = []string{
		"CustomResourceDefinition/alamedascalersCRD.yaml",
		"CustomResourceDefinition/alamedascalersV2CRD.yaml",
		"CustomResourceDefinition/alamedarecommendationsCRD.yaml",
		"CustomResourceDefinition/alamedanotificationchannels.yaml",
		"CustomResourceDefinition/alamedanotificationtopics.yaml",
	}

	alamedaScalerCRD = []string{
		"CustomResourceDefinition/alamedascalersCRD.yaml",
	}

	alamedaScalerCRDV2 = []string{
		"CustomResourceDefinition/alamedascalersV2CRD.yaml",
	}

	logPVCList = []string{
		"PersistentVolumeClaim/alameda-ai-log.yaml",
		"PersistentVolumeClaim/alameda-operator-log.yaml",
		"PersistentVolumeClaim/alameda-datahub-log.yaml",
		"PersistentVolumeClaim/alameda-evictioner-log.yaml",
		"PersistentVolumeClaim/admission-controller-log.yaml",
		"PersistentVolumeClaim/alameda-recommender-log.yaml",
		"PersistentVolumeClaim/alameda-executor-log.yaml",
		"PersistentVolumeClaim/alameda-analyzer-log.yaml",
		"PersistentVolumeClaim/alameda-dispatcher-log.yaml",
		"PersistentVolumeClaim/fedemeter-log.yaml",
	}

	dataPVCList = []string{
		"PersistentVolumeClaim/my-alamedainfluxdbPVC.yaml",
		"PersistentVolumeClaim/my-alamedagrafanaPVC.yaml",
		"PersistentVolumeClaim/alameda-ai-data.yaml",
		"PersistentVolumeClaim/alameda-operator-data.yaml",
		"PersistentVolumeClaim/alameda-datahub-data.yaml",
		"PersistentVolumeClaim/alameda-evictioner-data.yaml",
		"PersistentVolumeClaim/alameda-analyzer-data.yaml",
		"PersistentVolumeClaim/admission-controller-data.yaml",
		"PersistentVolumeClaim/alameda-recommender-data.yaml",
		"PersistentVolumeClaim/alameda-executor-data.yaml",
		"PersistentVolumeClaim/alameda-dispatcher-data.yaml",
		"PersistentVolumeClaim/fedemeter-data.yaml",
	}
)

// GetDispatcherResource returns resource that needs to be installed for Alameda-Dispathcer
func GetDispatcherResource() *Resource {
	r, _ := getResourceFromList(aiDispatcherList)
	return &r
}

// GetExcutionResource returns resource that needs to be installed for Execution
func GetExcutionResource() *Resource {
	r, _ := getResourceFromList(excutionList)
	return &r
}

// GetGUIResource returns resource that needs to be installed for GUI
func GetGUIResource() *Resource {
	r, _ := getResourceFromList(guiList)
	return &r
}

// GetFedemeterResource returns resource that needs to be installed for Federmeter
func GetFedemeterResource() *Resource {
	r, _ := getResourceFromList(fedemeterList)
	return &r
}

// GetWeavescopeResource returns resource that needs to be installed for weavescope
func GetWeavescopeResource() Resource {
	r, _ := getResourceFromList(weavescopeList)
	return r
}

// GetSelfDrivingRsource returns resource that needs to be installed for Alameda self driving
func GetSelfDrivingRsource() *Resource {
	r, _ := getResourceFromList(selfDrivingList)
	return &r
}

// GetAlamedaDatahubService returns service asset name of alameda-dathub
func GetAlamedaDatahubService() string {
	return "Service/alameda-datahubSV.yaml"
}

// GetAlamedaNotifierWebhookService returns service asset name of alameda-notifier-webhook-service
func GetAlamedaNotifierWebhookService() string {
	return "Service/alameda-notifier-webhook-service.yaml"
}

// GetAlamedaVPAExecutorService returns service asset name of alameda-vpa-executorSV
func GetAlamedaVPAExecutorService() string {
	return "Service/alameda-vpa-executorSV.yaml"
}

// GetAlamedaNotifierWebhookServerCertSecret returns secret asset name of alameda-notifier-webhook-server-cert
func GetAlamedaNotifierWebhookServerCertSecret() string {
	return "Secret/alameda-notifier-webhook-server-cert.yaml"
}

// GetAlamedaNotifierWebhookServerCertSecret returns secret asset name of alameda-notifier-webhook-server-cert
func GetAlamedaVPAExecutorServerCertSecret() string {
	return "Secret/alameda-vpa-executor-server-cert.yaml"
}

// GetCustomResourceDefinitions returns crd assets' name
func GetCustomResourceDefinitions() []string {
	return customResourceDefinitionList
}

type AlamedaServiceParamter struct {
	NameSpace                     string
	SelfDriving                   bool
	Platform                      string
	EnableExecution               bool
	EnableGUI                     bool
	EnableDispatcher              bool
	EnableWeavescope              bool
	Version                       string
	PrometheusService             string
	Storages                      []v1alpha1.StorageSpec
	ServiceExposures              []v1alpha1.ServiceExposureSpec
	InfluxdbSectionSet            v1alpha1.AlamedaComponentSpec
	GrafanaSectionSet             v1alpha1.AlamedaComponentSpec
	AlamedaAISectionSet           v1alpha1.AlamedaComponentSpec
	AlamedaOperatorSectionSet     v1alpha1.AlamedaComponentSpec
	AlamedaDatahubSectionSet      v1alpha1.AlamedaComponentSpec
	AlamedaEvictionerSectionSet   v1alpha1.AlamedaComponentSpec
	AdmissionControllerSectionSet v1alpha1.AlamedaComponentSpec
	AlamedaRecommenderSectionSet  v1alpha1.AlamedaComponentSpec
	AlamedaExecutorSectionSet     v1alpha1.AlamedaComponentSpec
	AlamedaDispatcherSectionSet   v1alpha1.AlamedaComponentSpec
	AlamedaFedemeterSectionSet    v1alpha1.AlamedaComponentSpec
	AlamedaWeavescopeSectionSet   v1alpha1.AlamedaComponentSpec
	AlamedaAnalyzerSectionSet     v1alpha1.AlamedaComponentSpec
	AlamedaNotifierSectionSet     v1alpha1.AlamedaComponentSpec
	AlamedaRabbitMQSectionSet     v1alpha1.AlamedaComponentSpec
	FederatoraiAgentSectionSet    v1alpha1.AlamedaComponentSpec
	CurrentCRDVersion             v1alpha1.AlamedaServiceStatusCRDVersion
	previousCRDVersion            v1alpha1.AlamedaServiceStatusCRDVersion
}

func NewAlamedaServiceParamter(instance *v1alpha1.AlamedaService) *AlamedaServiceParamter {
	asp := &AlamedaServiceParamter{
		NameSpace:                     instance.Namespace,
		SelfDriving:                   instance.Spec.SelfDriving,
		Platform:                      instance.Spec.Platform,
		EnableExecution:               instance.Spec.EnableExecution,
		EnableGUI:                     instance.Spec.EnableGUI,
		EnableDispatcher:              *instance.Spec.EnableDispatcher,
		EnableWeavescope:              instance.Spec.EnableWeavescope,
		Version:                       instance.Spec.Version,
		PrometheusService:             instance.Spec.PrometheusService,
		Storages:                      instance.Spec.Storages,
		ServiceExposures:              instance.Spec.ServiceExposures,
		InfluxdbSectionSet:            instance.Spec.InfluxdbSectionSet,
		GrafanaSectionSet:             instance.Spec.GrafanaSectionSet,
		AlamedaAISectionSet:           instance.Spec.AlamedaAISectionSet,
		AlamedaOperatorSectionSet:     instance.Spec.AlamedaOperatorSectionSet,
		AlamedaDatahubSectionSet:      instance.Spec.AlamedaDatahubSectionSet,
		AlamedaEvictionerSectionSet:   instance.Spec.AlamedaEvictionerSectionSet,
		AdmissionControllerSectionSet: instance.Spec.AdmissionControllerSectionSet,
		AlamedaRecommenderSectionSet:  instance.Spec.AlamedaRecommenderSectionSet,
		AlamedaExecutorSectionSet:     instance.Spec.AlamedaExecutorSectionSet,
		AlamedaDispatcherSectionSet:   instance.Spec.AlamedaDispatcherSectionSet,
		AlamedaRabbitMQSectionSet:     instance.Spec.AlamedaRabbitMQSectionSet,
		AlamedaAnalyzerSectionSet:     instance.Spec.AlamedaAnalyzerSectionSet,
		AlamedaFedemeterSectionSet:    instance.Spec.AlamedaFedemeterSectionSet,
		AlamedaWeavescopeSectionSet:   instance.Spec.AlamedaWeavescopeSectionSet,
		FederatoraiAgentSectionSet:    instance.Spec.FederatoraiAgentSectionSet,
		AlamedaNotifierSectionSet:     instance.Spec.AlamedaNotifierSectionSet,
		CurrentCRDVersion:             instance.Status.CRDVersion,
		previousCRDVersion:            instance.Status.CRDVersion,
	}
	asp.changeScalerCRDVersion()
	return asp
}

// GetInstallResource returns resources that the AlamedaServiceParamter needs to install
func (asp *AlamedaServiceParamter) GetInstallResource() *Resource {

	var resource Resource

	for _, defaultInstallList := range defaultInstallLists {
		defaultResource, _ := getResourceFromList(defaultInstallList)
		resource.add(defaultResource)
	}

	pvcList := asp.getInstallPersistentVolumeClaimSource()
	pvcResource, _ := getResourceFromList(pvcList)
	resource.add(pvcResource)

	if asp.SelfDriving {
		r, _ := getResourceFromList(selfDrivingList)
		resource.add(r)
	}
	if asp.EnableGUI {
		r, _ := getResourceFromList(guiList)
		resource.add(r)
	}
	if asp.EnableExecution {
		r, _ := getResourceFromList(excutionList)
		resource.add(r)
	}
	if asp.EnableDispatcher {
		r := GetDispatcherResource()
		resource.add(*r)
	}
	if asp.EnableWeavescope {
		weavescopeResource := GetWeavescopeResource()
		resource.add(weavescopeResource)
	}

	if asp.hasToInstallAlamedaAcalerV2() {
		v2Resource, _ := getResourceFromList(alamedaScalerCRDV2)
		resource.add(v2Resource)

		defaultResource, _ := getResourceFromList(alamedaScalerCRD)
		resource.delete(defaultResource)
	}

	return &resource
}

func (asp *AlamedaServiceParamter) GetUninstallPersistentVolumeClaimSource() *Resource {
	pvc := []string{}

	appendLogPVC := false
	appendDataPVC := false
	for _, v := range asp.Storages {
		if v.Type != v1alpha1.PVC {
			if v.Usage == v1alpha1.Log {
				appendLogPVC = true
			} else if v.Usage == v1alpha1.Data {
				appendDataPVC = true
			} else if v.Usage == v1alpha1.Empty {
				appendLogPVC = true
				appendDataPVC = true
			}
		}
	}
	if appendLogPVC {
		pvc = append(pvc, logPVCList...)
	}
	if appendDataPVC {
		pvc = append(pvc, dataPVCList...)
	}

	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaAISectionSet.Storages, "PersistentVolumeClaim/alameda-ai-log.yaml", v1alpha1.Log)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaOperatorSectionSet.Storages, "PersistentVolumeClaim/alameda-operator-log.yaml", v1alpha1.Log)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaDatahubSectionSet.Storages, "PersistentVolumeClaim/alameda-datahub-log.yaml", v1alpha1.Log)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaEvictionerSectionSet.Storages, "PersistentVolumeClaim/alameda-evictioner-log.yaml", v1alpha1.Log)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AdmissionControllerSectionSet.Storages, "PersistentVolumeClaim/admission-controller-log.yaml", v1alpha1.Log)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaRecommenderSectionSet.Storages, "PersistentVolumeClaim/alameda-recommender-log.yaml", v1alpha1.Log)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaExecutorSectionSet.Storages, "PersistentVolumeClaim/alameda-executor-log.yaml", v1alpha1.Log)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaDispatcherSectionSet.Storages, "PersistentVolumeClaim/alameda-dispatcher-log.yaml", v1alpha1.Log)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaAnalyzerSectionSet.Storages, "PersistentVolumeClaim/alameda-analyzer-log.yaml", v1alpha1.Log)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaFedemeterSectionSet.Storages, "PersistentVolumeClaim/fedemeter-log.yaml", v1alpha1.Log)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaNotifierSectionSet.Storages, "PersistentVolumeClaim/alameda-notifier-log.yaml", v1alpha1.Log)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.FederatoraiAgentSectionSet.Storages, "PersistentVolumeClaim/federatorai-agent-log.yaml", v1alpha1.Log)

	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.InfluxdbSectionSet.Storages, "PersistentVolumeClaim/my-alamedainfluxdbPVC.yaml", v1alpha1.Data)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.GrafanaSectionSet.Storages, "PersistentVolumeClaim/my-alamedagrafanaPVC.yaml", v1alpha1.Data)

	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaAISectionSet.Storages, "PersistentVolumeClaim/alameda-ai-data.yaml", v1alpha1.Data)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaOperatorSectionSet.Storages, "PersistentVolumeClaim/alameda-operator-data.yaml", v1alpha1.Data)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaDatahubSectionSet.Storages, "PersistentVolumeClaim/alameda-datahub-data.yaml", v1alpha1.Data)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaEvictionerSectionSet.Storages, "PersistentVolumeClaim/alameda-evictioner-data.yaml", v1alpha1.Data)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AdmissionControllerSectionSet.Storages, "PersistentVolumeClaim/admission-controller-data.yaml", v1alpha1.Data)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaRecommenderSectionSet.Storages, "PersistentVolumeClaim/alameda-recommender-data.yaml", v1alpha1.Data)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaExecutorSectionSet.Storages, "PersistentVolumeClaim/alameda-executor-data.yaml", v1alpha1.Data)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaDispatcherSectionSet.Storages, "PersistentVolumeClaim/alameda-dispatcher-data.yaml", v1alpha1.Data)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaAnalyzerSectionSet.Storages, "PersistentVolumeClaim/alameda-analyzer-data.yaml", v1alpha1.Data)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaFedemeterSectionSet.Storages, "PersistentVolumeClaim/fedemeter-data.yaml", v1alpha1.Data)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.AlamedaNotifierSectionSet.Storages, "PersistentVolumeClaim/alameda-notifier-data.yaml", v1alpha1.Data)
	pvc = sectionUninstallPersistentVolumeClaimSource(pvc, asp.FederatoraiAgentSectionSet.Storages, "PersistentVolumeClaim/federatorai-agent-data.yaml", v1alpha1.Data)

	return &Resource{
		PersistentVolumeClaimList: pvc,
	}
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

func (asp *AlamedaServiceParamter) getInstallPersistentVolumeClaimSource() []string {

	pvc := make([]string, 0)

	// get install resource
	gloabalLogFlag := false
	gloabalDataFlag := false
	for _, storage := range asp.Storages {
		if storage.Type == v1alpha1.PVC {
			switch storage.Usage {
			case v1alpha1.Empty:
				gloabalLogFlag = true
				gloabalDataFlag = true
			case v1alpha1.Log:
				gloabalLogFlag = true
			case v1alpha1.Data:
				gloabalDataFlag = true
			}
		}
	}
	if gloabalLogFlag { //Gloabal append
		pvc = append(pvc, logPVCList...)
	}
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AlamedaAISectionSet.Storages, "PersistentVolumeClaim/alameda-ai-log.yaml", v1alpha1.Log)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AlamedaOperatorSectionSet.Storages, "PersistentVolumeClaim/alameda-operator-log.yaml", v1alpha1.Log)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AlamedaDatahubSectionSet.Storages, "PersistentVolumeClaim/alameda-datahub-log.yaml", v1alpha1.Log)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AlamedaEvictionerSectionSet.Storages, "PersistentVolumeClaim/alameda-evictioner-log.yaml", v1alpha1.Log)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AdmissionControllerSectionSet.Storages, "PersistentVolumeClaim/admission-controller-log.yaml", v1alpha1.Log)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AlamedaRecommenderSectionSet.Storages, "PersistentVolumeClaim/alameda-recommender-log.yaml", v1alpha1.Log)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AlamedaExecutorSectionSet.Storages, "PersistentVolumeClaim/alameda-executor-log.yaml", v1alpha1.Log)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AlamedaDispatcherSectionSet.Storages, "PersistentVolumeClaim/alameda-dispatcher-log.yaml", v1alpha1.Log)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AlamedaAnalyzerSectionSet.Storages, "PersistentVolumeClaim/alameda-analyzer-log.yaml", v1alpha1.Log)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AlamedaFedemeterSectionSet.Storages, "PersistentVolumeClaim/fedemeter-log.yaml", v1alpha1.Log)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AlamedaNotifierSectionSet.Storages, "PersistentVolumeClaim/alameda-notifier-log.yaml", v1alpha1.Log)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.FederatoraiAgentSectionSet.Storages, "PersistentVolumeClaim/federatorai-agent-log.yaml", v1alpha1.Log)
	if gloabalDataFlag {
		pvc = append(pvc, dataPVCList...)
	}
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.InfluxdbSectionSet.Storages, "PersistentVolumeClaim/my-alamedainfluxdbPVC.yaml", v1alpha1.Data)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.GrafanaSectionSet.Storages, "PersistentVolumeClaim/my-alamedagrafanaPVC.yaml", v1alpha1.Data)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AlamedaAISectionSet.Storages, "PersistentVolumeClaim/alameda-ai-data.yaml", v1alpha1.Data)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AlamedaOperatorSectionSet.Storages, "PersistentVolumeClaim/alameda-operator-data.yaml", v1alpha1.Data)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AlamedaDatahubSectionSet.Storages, "PersistentVolumeClaim/alameda-datahub-data.yaml", v1alpha1.Data)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AlamedaEvictionerSectionSet.Storages, "PersistentVolumeClaim/alameda-evictioner-data.yaml", v1alpha1.Data)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AdmissionControllerSectionSet.Storages, "PersistentVolumeClaim/admission-controller-data.yaml", v1alpha1.Data)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AlamedaRecommenderSectionSet.Storages, "PersistentVolumeClaim/alameda-recommender-data.yaml", v1alpha1.Data)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AlamedaExecutorSectionSet.Storages, "PersistentVolumeClaim/alameda-executor-data.yaml", v1alpha1.Data)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AlamedaDispatcherSectionSet.Storages, "PersistentVolumeClaim/alameda-dispatcher-data.yaml", v1alpha1.Data)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AlamedaAnalyzerSectionSet.Storages, "PersistentVolumeClaim/alameda-analyzer-data.yaml", v1alpha1.Data)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AlamedaFedemeterSectionSet.Storages, "PersistentVolumeClaim/fedemeter-data.yaml", v1alpha1.Data)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.AlamedaNotifierSectionSet.Storages, "PersistentVolumeClaim/alameda-notifier-data.yaml", v1alpha1.Data)
	pvc = sectioninstallPersistentVolumeClaimSource(pvc, asp.FederatoraiAgentSectionSet.Storages, "PersistentVolumeClaim/federatorai-agent-data.yaml", v1alpha1.Data)
	return pvc

}

func (asp *AlamedaServiceParamter) changeScalerCRDVersion() {

	alamedaOperatorVersion := util.OriAlamedaOperatorVersion
	if asp.Version != "" {
		alamedaOperatorVersion = asp.Version
	}
	if asp.AlamedaOperatorSectionSet.Version != "" {
		alamedaOperatorVersion = asp.AlamedaOperatorSectionSet.Version
	}
	if util.StringInSlice(alamedaOperatorVersion, util.V1scalerOperatorVersionList) { //check current operatorVersion used scaler version is scaler V1
		asp.CurrentCRDVersion.ScalerVersion = util.AlamedaScalerVersion[0]
		asp.CurrentCRDVersion.CRDName = util.AlamedaScalerName
	} else {
		asp.CurrentCRDVersion.ScalerVersion = util.AlamedaScalerVersion[1]
		asp.CurrentCRDVersion.CRDName = util.AlamedaScalerName
	}
	if asp.CurrentCRDVersion.ScalerVersion != asp.previousCRDVersion.ScalerVersion {
		asp.SetCurrentCRDChangeVersionToTrue()
	}
}

func (asp *AlamedaServiceParamter) hasToInstallAlamedaAcalerV2() bool {

	alamedaOperatorVersion := util.OriAlamedaOperatorVersion
	if asp.Version != "" {
		alamedaOperatorVersion = asp.Version
	}
	if asp.AlamedaOperatorSectionSet.Version != "" {
		alamedaOperatorVersion = asp.AlamedaOperatorSectionSet.Version
	}
	if util.StringInSlice(alamedaOperatorVersion, util.V1scalerOperatorVersionList) { //check current operatorVersion used scaler version is scaler V1
		return false
	}
	return true
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

type Resource struct {
	ClusterRoleBindingList             []string
	ClusterRoleList                    []string
	ServiceAccountList                 []string
	CustomResourceDefinitionList       []string
	ConfigMapList                      []string
	ServiceList                        []string
	DeploymentList                     []string
	SecretList                         []string
	PersistentVolumeClaimList          []string
	AlamedaScalerList                  []string
	RouteList                          []string
	StatefulSetList                    []string
	IngressList                        []string
	PodSecurityPolicyList              []string
	DaemonSetList                      []string
	SecurityContextConstraintsList     []string
	RoleBindingList                    []string
	RoleList                           []string
	MutatingWebhookConfigurationList   []string
	ValidatingWebhookConfigurationList []string
	AlamedaNotificationChannelList     []string
	AlamedaNotificationTopic           []string
}

func (r *Resource) GetAll() []string {

	files := make([]string, 0)

	for _, file := range r.CustomResourceDefinitionList {
		files = append(files, file)
	}
	for _, file := range r.PodSecurityPolicyList {
		files = append(files, file)
	}
	for _, file := range r.SecurityContextConstraintsList {
		files = append(files, file)
	}
	for _, file := range r.ClusterRoleList {
		files = append(files, file)
	}
	for _, file := range r.ServiceAccountList {
		files = append(files, file)
	}
	for _, file := range r.ClusterRoleBindingList {
		files = append(files, file)
	}
	for _, file := range r.RoleList {
		files = append(files, file)
	}
	for _, file := range r.RoleBindingList {
		files = append(files, file)
	}
	for _, file := range r.SecretList {
		files = append(files, file)
	}
	for _, file := range r.PersistentVolumeClaimList {
		files = append(files, file)
	}
	for _, file := range r.ConfigMapList {
		files = append(files, file)
	}
	for _, file := range r.ServiceList {
		files = append(files, file)
	}
	for _, file := range r.MutatingWebhookConfigurationList {
		files = append(files, file)
	}
	for _, file := range r.ValidatingWebhookConfigurationList {
		files = append(files, file)
	}
	for _, file := range r.DeploymentList {
		files = append(files, file)
	}
	for _, file := range r.StatefulSetList {
		files = append(files, file)
	}
	for _, file := range r.IngressList {
		files = append(files, file)
	}
	for _, file := range r.RouteList {
		files = append(files, file)
	}
	for _, file := range r.DaemonSetList {
		files = append(files, file)
	}
	for _, file := range r.AlamedaNotificationChannelList {
		files = append(files, file)
	}
	for _, file := range r.AlamedaNotificationTopic {
		files = append(files, file)
	}

	for _, file := range r.AlamedaScalerList {
		files = append(files, file)
	}

	return files
}

func (r *Resource) add(in Resource) {
	r.ClusterRoleBindingList = append(r.ClusterRoleBindingList, in.ClusterRoleBindingList...)
	r.ClusterRoleList = append(r.ClusterRoleList, in.ClusterRoleList...)
	r.ServiceAccountList = append(r.ServiceAccountList, in.ServiceAccountList...)
	r.CustomResourceDefinitionList = append(r.CustomResourceDefinitionList, in.CustomResourceDefinitionList...)
	r.ConfigMapList = append(r.ConfigMapList, in.ConfigMapList...)
	r.ServiceList = append(r.ServiceList, in.ServiceList...)
	r.DeploymentList = append(r.DeploymentList, in.DeploymentList...)
	r.SecretList = append(r.SecretList, in.SecretList...)
	r.PersistentVolumeClaimList = append(r.PersistentVolumeClaimList, in.PersistentVolumeClaimList...)
	r.AlamedaScalerList = append(r.AlamedaScalerList, in.AlamedaScalerList...)
	r.RouteList = append(r.RouteList, in.RouteList...)
	r.StatefulSetList = append(r.StatefulSetList, in.StatefulSetList...)
	r.IngressList = append(r.IngressList, in.IngressList...)
	r.PodSecurityPolicyList = append(r.PodSecurityPolicyList, in.PodSecurityPolicyList...)
	r.DaemonSetList = append(r.DaemonSetList, in.DaemonSetList...)
	r.SecurityContextConstraintsList = append(r.SecurityContextConstraintsList, in.SecurityContextConstraintsList...)
	r.RoleBindingList = append(r.RoleBindingList, in.RoleBindingList...)
	r.RoleList = append(r.RoleList, in.RoleList...)
	r.MutatingWebhookConfigurationList = append(r.MutatingWebhookConfigurationList, in.MutatingWebhookConfigurationList...)
	r.ValidatingWebhookConfigurationList = append(r.ValidatingWebhookConfigurationList, in.ValidatingWebhookConfigurationList...)
	r.AlamedaNotificationChannelList = append(r.AlamedaNotificationChannelList, in.AlamedaNotificationChannelList...)
	r.AlamedaNotificationTopic = append(r.AlamedaNotificationTopic, in.AlamedaNotificationTopic...)
}

func (r *Resource) delete(in Resource) {
	r.ClusterRoleBindingList = util.StringSliceDelete(r.ClusterRoleBindingList, in.ClusterRoleBindingList)
	r.ClusterRoleList = util.StringSliceDelete(r.ClusterRoleList, in.ClusterRoleList)
	r.ServiceAccountList = util.StringSliceDelete(r.ServiceAccountList, in.ServiceAccountList)
	r.CustomResourceDefinitionList = util.StringSliceDelete(r.CustomResourceDefinitionList, in.CustomResourceDefinitionList)
	r.ConfigMapList = util.StringSliceDelete(r.ConfigMapList, in.ConfigMapList)
	r.ServiceList = util.StringSliceDelete(r.ServiceList, in.ServiceList)
	r.DeploymentList = util.StringSliceDelete(r.DeploymentList, in.DeploymentList)
	r.SecretList = util.StringSliceDelete(r.SecretList, in.SecretList)
	r.PersistentVolumeClaimList = util.StringSliceDelete(r.PersistentVolumeClaimList, in.PersistentVolumeClaimList)
	r.AlamedaScalerList = util.StringSliceDelete(r.AlamedaScalerList, in.AlamedaScalerList)
	r.RouteList = util.StringSliceDelete(r.RouteList, in.RouteList)
	r.StatefulSetList = util.StringSliceDelete(r.StatefulSetList, in.StatefulSetList)
	r.IngressList = util.StringSliceDelete(r.IngressList, in.IngressList)
	r.PodSecurityPolicyList = util.StringSliceDelete(r.PodSecurityPolicyList, in.PodSecurityPolicyList)
	r.DaemonSetList = util.StringSliceDelete(r.DaemonSetList, in.DaemonSetList)
	r.SecurityContextConstraintsList = util.StringSliceDelete(r.SecurityContextConstraintsList, in.SecurityContextConstraintsList)
	r.RoleBindingList = util.StringSliceDelete(r.RoleBindingList, in.RoleBindingList)
	r.RoleList = util.StringSliceDelete(r.RoleList, in.RoleList)
	r.MutatingWebhookConfigurationList = util.StringSliceDelete(r.MutatingWebhookConfigurationList, in.MutatingWebhookConfigurationList)
	r.ValidatingWebhookConfigurationList = util.StringSliceDelete(r.ValidatingWebhookConfigurationList, in.ValidatingWebhookConfigurationList)
	r.AlamedaNotificationChannelList = util.StringSliceDelete(r.AlamedaNotificationChannelList, in.AlamedaNotificationChannelList)
	r.AlamedaNotificationTopic = util.StringSliceDelete(r.AlamedaNotificationTopic, in.AlamedaNotificationTopic)
}

func getResourceFromList(resourceList []string) (Resource, error) {

	var clusterRoleBindingList = make([]string, 0)
	var clusterRoleList = make([]string, 0)
	var serviceAccountList = make([]string, 0)
	var customResourceDefinitionList = make([]string, 0)
	var configMapList = make([]string, 0)
	var serviceList = make([]string, 0)
	var deploymentList = make([]string, 0)
	var secretList = make([]string, 0)
	var persistentVolumeClaimList = make([]string, 0)
	var alamedaScalerList = make([]string, 0)
	var routeList = make([]string, 0)
	var statefulSetList = make([]string, 0)
	var ingressList = make([]string, 0)
	var podSecurityPolicyList = make([]string, 0)
	var daemonSetList = make([]string, 0)
	var securityContextConstraintsList = make([]string, 0)
	var roleBindingList = make([]string, 0)
	var roleList = make([]string, 0)
	var mutatingWebhookConfigurationList = make([]string, 0)
	var validatingWebhookConfigurationList = make([]string, 0)
	var alamedaNotificationChannelList = make([]string, 0)
	var alamedaNotificationTopicList = make([]string, 0)

	for _, assetFile := range resourceList {
		if len(strings.Split(assetFile, "/")) > 0 {
			switch kind := strings.Split(assetFile, "/")[0]; kind {
			case "ClusterRoleBinding":
				clusterRoleBindingList = append(clusterRoleBindingList, assetFile)
			case "ClusterRole":
				clusterRoleList = append(clusterRoleList, assetFile)
			case "ServiceAccount":
				serviceAccountList = append(serviceAccountList, assetFile)
			case "CustomResourceDefinition":
				customResourceDefinitionList = append(customResourceDefinitionList, assetFile)
			case "ConfigMap":
				configMapList = append(configMapList, assetFile)
			case "Service":
				serviceList = append(serviceList, assetFile)
			case "Deployment":
				deploymentList = append(deploymentList, assetFile)
			case "Secret":
				secretList = append(secretList, assetFile)
			case "PersistentVolumeClaim":
				persistentVolumeClaimList = append(persistentVolumeClaimList, assetFile)
			case "AlamedaScaler":
				alamedaScalerList = append(alamedaScalerList, assetFile)
			case "Route":
				routeList = append(routeList, assetFile)
			case "StatefulSet":
				statefulSetList = append(statefulSetList, assetFile)
			case "Ingress":
				ingressList = append(ingressList, assetFile)
			case "PodSecurityPolicy":
				podSecurityPolicyList = append(podSecurityPolicyList, assetFile)
			case "DaemonSet":
				daemonSetList = append(daemonSetList, assetFile)
			case "SecurityContextConstraint":
				securityContextConstraintsList = append(securityContextConstraintsList, assetFile)
			case "RoleBinding":
				roleBindingList = append(roleBindingList, assetFile)
			case "Role":
				roleList = append(roleList, assetFile)
			case "MutatingWebhookConfiguration":
				mutatingWebhookConfigurationList = append(mutatingWebhookConfigurationList, assetFile)
			case "ValidatingWebhookConfiguration":
				validatingWebhookConfigurationList = append(validatingWebhookConfigurationList, assetFile)
			case "AlamedaNotificationChannel":
				alamedaNotificationChannelList = append(alamedaNotificationChannelList, assetFile)
			case "AlamedaNotificationTopic":
				alamedaNotificationTopicList = append(alamedaNotificationTopicList, assetFile)
			default:
				return Resource{}, errors.Errorf("unknown kind \"%s\"", kind)
			}
		} else {
			return Resource{}, errors.Errorf("invalid format \"%s\"", assetFile)
		}
	}

	return Resource{
		ClusterRoleBindingList:             clusterRoleBindingList,
		ClusterRoleList:                    clusterRoleList,
		ServiceAccountList:                 serviceAccountList,
		CustomResourceDefinitionList:       customResourceDefinitionList,
		ConfigMapList:                      configMapList,
		ServiceList:                        serviceList,
		DeploymentList:                     deploymentList,
		SecretList:                         secretList,
		PersistentVolumeClaimList:          persistentVolumeClaimList,
		AlamedaScalerList:                  alamedaScalerList,
		RouteList:                          routeList,
		StatefulSetList:                    statefulSetList,
		IngressList:                        ingressList,
		PodSecurityPolicyList:              podSecurityPolicyList,
		DaemonSetList:                      daemonSetList,
		SecurityContextConstraintsList:     securityContextConstraintsList,
		RoleBindingList:                    roleBindingList,
		RoleList:                           roleList,
		MutatingWebhookConfigurationList:   mutatingWebhookConfigurationList,
		ValidatingWebhookConfigurationList: validatingWebhookConfigurationList,
		AlamedaNotificationChannelList:     alamedaNotificationChannelList,
		AlamedaNotificationTopic:           alamedaNotificationTopicList,
	}, nil

}
