package componentsectionset

import (
	"fmt"

	"github.com/containers-ai/federatorai-operator/pkg/apis/federatorai/v1alpha1"

	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/alamedaserviceparamter"
	"github.com/containers-ai/federatorai-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func SectionSetParamterToDeployment(dep *appsv1.Deployment, asp *alamedaserviceparamter.AlamedaServiceParamter) {
	switch dep.Name {
	case util.AlamedaaiDPN:
		{
			//set imageStuct
			util.SetImageStruct(dep, asp.AlamedaAISectionSet, util.AlamedaaiCTN)
			//set imagePullPolicy
			util.SetImagePullPolicy(dep, util.AlamedaaiCTN, asp.AlamedaAISectionSet.ImagePullPolicy)
			//set persistentVolumeClaim to mountPath
			util.SetStorageToVolumeSource(dep, asp.AlamedaAISectionSet.Storages, "alameda-ai-type.pvc", util.AlamedaGroup)
			util.SetStorageToMountPath(dep, asp.AlamedaAISectionSet.Storages, util.AlamedaaiCTN, "alameda-ai-type-storage", util.AlamedaGroup)
		}
	case util.AlamedaoperatorDPN:
		{
			util.SetImageStruct(dep, asp.AlamedaOperatorSectionSet, util.AlamedaoperatorCTN)
			util.SetImagePullPolicy(dep, util.AlamedaoperatorCTN, asp.AlamedaOperatorSectionSet.ImagePullPolicy)
			util.SetStorageToVolumeSource(dep, asp.AlamedaOperatorSectionSet.Storages, "alameda-operator-type.pvc", util.AlamedaGroup)
			util.SetStorageToMountPath(dep, asp.AlamedaOperatorSectionSet.Storages, util.AlamedaoperatorCTN, "alameda-operator-type-storage", util.AlamedaGroup)
		}
	case util.AlamedadatahubDPN:
		{
			util.SetImageStruct(dep, asp.AlamedaDatahubSectionSet, util.AlamedadatahubCTN)
			util.SetImagePullPolicy(dep, util.AlamedadatahubCTN, asp.AlamedaDatahubSectionSet.ImagePullPolicy)
			util.SetStorageToVolumeSource(dep, asp.AlamedaDatahubSectionSet.Storages, "alameda-datahub-type.pvc", util.AlamedaGroup)
			util.SetStorageToMountPath(dep, asp.AlamedaDatahubSectionSet.Storages, util.AlamedadatahubCTN, "alameda-datahub-type-storage", util.AlamedaGroup)
		}
	case util.AlamedaevictionerDPN:
		{
			util.SetImageStruct(dep, asp.AlamedaEvictionerSectionSet, util.AlamedaevictionerCTN)
			util.SetImagePullPolicy(dep, util.AlamedaevictionerCTN, asp.AlamedaEvictionerSectionSet.ImagePullPolicy)
			util.SetStorageToVolumeSource(dep, asp.AlamedaEvictionerSectionSet.Storages, "alameda-evictioner-type.pvc", util.AlamedaGroup)
			util.SetStorageToMountPath(dep, asp.AlamedaEvictionerSectionSet.Storages, util.AlamedaevictionerCTN, "alameda-evictioner-type-storage", util.AlamedaGroup)
		}
	case util.AdmissioncontrollerDPN:
		{
			util.SetImageStruct(dep, asp.AdmissionControllerSectionSet, util.AdmissioncontrollerCTN)
			util.SetImagePullPolicy(dep, util.AdmissioncontrollerCTN, asp.AdmissionControllerSectionSet.ImagePullPolicy)
			util.SetStorageToVolumeSource(dep, asp.AdmissionControllerSectionSet.Storages, "admission-controller-type.pvc", util.AlamedaGroup)
			util.SetStorageToMountPath(dep, asp.AdmissionControllerSectionSet.Storages, util.AdmissioncontrollerCTN, "admission-controller-type-storage", util.AlamedaGroup)
		}
	case util.AlamedarecommenderDPN:
		{
			util.SetImageStruct(dep, asp.AlamedaRecommenderSectionSet, util.AlamedarecommenderCTN)
			util.SetImagePullPolicy(dep, util.AlamedarecommenderCTN, asp.AlamedaRecommenderSectionSet.ImagePullPolicy)
			util.SetStorageToVolumeSource(dep, asp.AlamedaRecommenderSectionSet.Storages, "alameda-recommender-type.pvc", util.AlamedaGroup)
			util.SetStorageToMountPath(dep, asp.AlamedaRecommenderSectionSet.Storages, util.AlamedarecommenderCTN, "alameda-recommender-type.pvc", util.AlamedaGroup)
		}
	case util.AlamedaexecutorDPN:
		{
			util.SetImageStruct(dep, asp.AlamedaExecutorSectionSet, util.AlamedaexecutorCTN)
			util.SetImagePullPolicy(dep, util.AlamedaexecutorCTN, asp.AlamedaExecutorSectionSet.ImagePullPolicy)
			util.SetStorageToVolumeSource(dep, asp.AlamedaExecutorSectionSet.Storages, "alameda-executor-type.pvc", util.AlamedaGroup)
			util.SetStorageToMountPath(dep, asp.AlamedaExecutorSectionSet.Storages, util.AlamedaexecutorCTN, "alameda-executor-type.pvc", util.AlamedaGroup)
		}
	case util.AlamedadispatcherDPN:
		{
			util.SetImageStruct(dep, asp.AlamedaDispatcherSectionSet, util.AlamedadispatcherCTN)
			util.SetImagePullPolicy(dep, util.AlamedadispatcherCTN, asp.AlamedaDispatcherSectionSet.ImagePullPolicy)
			util.SetStorageToVolumeSource(dep, asp.AlamedaDispatcherSectionSet.Storages, "alameda-dispatcher-type.pvc", util.AlamedaGroup)
			util.SetStorageToMountPath(dep, asp.AlamedaDispatcherSectionSet.Storages, util.AlamedadispatcherCTN, "alameda-dispatcher-type.pvc", util.AlamedaGroup)
		}
	case util.AlamedaanalyzerDPN:
		{
			util.SetImageStruct(dep, asp.AlamedaAnalyzerSectionSet, util.AlamedaanalyzerCTN)
			util.SetImagePullPolicy(dep, util.AlamedaanalyzerCTN, asp.AlamedaAnalyzerSectionSet.ImagePullPolicy)
			util.SetStorageToVolumeSource(dep, asp.AlamedaAnalyzerSectionSet.Storages, "alameda-analyzer-type.pvc", util.AlamedaGroup)
			util.SetStorageToMountPath(dep, asp.AlamedaAnalyzerSectionSet.Storages, util.AlamedaanalyzerCTN, "alameda-analyzer-type.pvc", util.AlamedaGroup)
		}
	case util.FedemeterDPN:
		{
			util.SetImageStruct(dep, asp.AlamedaFedemeterSectionSet, util.FedemeterCTN)
			util.SetImagePullPolicy(dep, util.FedemeterCTN, asp.AlamedaFedemeterSectionSet.ImagePullPolicy)
			util.SetStorageToVolumeSource(dep, asp.AlamedaFedemeterSectionSet.Storages, "fedemeter-type.pvc", util.FedemeterGroup)
			util.SetStorageToMountPath(dep, asp.AlamedaFedemeterSectionSet.Storages, util.FedemeterCTN, "fedemeter-type.pvc", util.FedemeterGroup)
		}
	case util.InfluxdbDPN:
		{
			util.SetImageStruct(dep, asp.InfluxdbSectionSet, util.InfluxdbCTN)
			util.SetImagePullPolicy(dep, util.InfluxdbCTN, asp.InfluxdbSectionSet.ImagePullPolicy)
			util.SetStorageToVolumeSource(dep, asp.InfluxdbSectionSet.Storages, "my-alameda.influxdb-type.pvc", util.InfluxDBGroup)
			util.SetStorageToMountPath(dep, asp.InfluxdbSectionSet.Storages, util.InfluxdbCTN, "influxdb-type-storage", util.InfluxDBGroup)
		}
	case util.GrafanaDPN:
		{
			util.SetBootStrapImageStruct(dep, asp.GrafanaSectionSet, util.GetTokenCTN)
			util.SetImageStruct(dep, asp.GrafanaSectionSet, util.GrafanaCTN)
			util.SetImagePullPolicy(dep, util.GrafanaCTN, asp.GrafanaSectionSet.ImagePullPolicy)
			util.SetStorageToVolumeSource(dep, asp.GrafanaSectionSet.Storages, "my-alameda.grafana-type.pvc", util.GrafanaGroup)
			util.SetStorageToMountPath(dep, asp.GrafanaSectionSet.Storages, util.GrafanaCTN, "grafana-type-storage", util.GrafanaGroup)
		}
	case util.AlamedaweavescopeDPN:
		{
			util.SetImageStruct(dep, asp.AlamedaWeavescopeSectionSet, util.AlamedaweavescopeCTN)
			util.SetImagePullPolicy(dep, util.AlamedaweavescopeCTN, asp.AlamedaWeavescopeSectionSet.ImagePullPolicy)
		}
	case util.AlamedaweavescopeProbeDPN:
		{
			util.SetImageStruct(dep, asp.AlamedaWeavescopeSectionSet, util.AlamedaweavescopeProbeCTN)
			util.SetImagePullPolicy(dep, util.AlamedaweavescopeProbeCTN, asp.AlamedaWeavescopeSectionSet.ImagePullPolicy)
		}
	}
}

func SectionSetParamterToDaemonSet(ds *appsv1.DaemonSet, asp *alamedaserviceparamter.AlamedaServiceParamter) {
	switch ds.Name {
	case util.AlamedaweavescopeAgentDS:
		{
			util.SetDaemonSetImageStruct(ds, asp.AlamedaWeavescopeSectionSet, util.AlamedaweavescopeAgentCTN)
			util.SetDaemonSetImagePullPolicy(ds, util.AlamedaweavescopeAgentCTN, asp.AlamedaWeavescopeSectionSet.ImagePullPolicy)
		}
	}
}

func SectionSetParamterToPersistentVolumeClaim(pvc *corev1.PersistentVolumeClaim, asp *alamedaserviceparamter.AlamedaServiceParamter) {
	for _, pvcusage := range v1alpha1.PvcUsage {
		switch pvc.Name {
		case fmt.Sprintf("alameda-ai-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.AlamedaAISectionSet.Storages, pvcusage)
			}
		case fmt.Sprintf("alameda-operator-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.AlamedaOperatorSectionSet.Storages, pvcusage)
			}
		case fmt.Sprintf("alameda-datahub-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.AlamedaDatahubSectionSet.Storages, pvcusage)
			}
		case fmt.Sprintf("alameda-evictioner-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.AlamedaEvictionerSectionSet.Storages, pvcusage)
			}
		case fmt.Sprintf("admission-controller-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.AdmissionControllerSectionSet.Storages, pvcusage)
			}
		case fmt.Sprintf("alameda-recommender-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.AlamedaRecommenderSectionSet.Storages, pvcusage)
			}
		case fmt.Sprintf("alameda-executor-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.AlamedaExecutorSectionSet.Storages, pvcusage)
			}
		case fmt.Sprintf("alameda-dispatcher-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.AlamedaDispatcherSectionSet.Storages, pvcusage)
			}
		case fmt.Sprintf("alameda-analyzer-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.AlamedaAnalyzerSectionSet.Storages, pvcusage)
			}
		case fmt.Sprintf("fedemeter-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.AlamedaFedemeterSectionSet.Storages, pvcusage)
			}
		case fmt.Sprintf("my-alameda.influxdb-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.InfluxdbSectionSet.Storages, pvcusage)
			}
		case fmt.Sprintf("my-alameda.grafana-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.GrafanaSectionSet.Storages, pvcusage)
			}
		}
	}
}
