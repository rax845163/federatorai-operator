package alamedaserviceparamter

import (
	"strings"

	federatoraiv1alpha1 "github.com/containers-ai/federatorai-operator/pkg/apis/federatorai/v1alpha1"
)

var (
	crdList = []string{
		"CustomResourceDefinition/alamedarecommendationsCRD.yaml",
		"CustomResourceDefinition/alamedascalersCRD.yaml",
	}
	cmList = []string{}

	svList = []string{"Service/alameda-datahubSV.yaml",
		"Service/alameda-influxdbSV.yaml"}

	depList = []string{"Deployment/alameda-datahubDM.yaml",
		"Deployment/alameda-operatorDM.yaml",
		"Deployment/alameda-influxdbDM.yaml",
		"Deployment/alameda-aiDM.yaml"}

	guiList = []string{"ConfigMap/grafana-datasources.yaml",
		"Deployment/alameda-grafanaDM.yaml",
		"Service/alameda-grafanaSV.yaml",
	}
	excutionList = []string{"Deployment/admission-controllerDM.yaml",
		"Deployment/alameda-evictionerDM.yaml",
		"Service/admission-controllerSV.yaml",
	}
)

type AlamedaServiceParamter struct {
	//AlmedaInstallOrUninstall bool
	NameSpace             string
	EnableExecution       bool
	EnableGUI             bool
	Version               string
	PrometheusService     string
	PersistentVolumeClaim string
	GuiFlag               bool
	ExcutionFlag          bool
	Guicomponent          []string
	Excutioncomponent     []string
}

type Resource struct {
	CustomResourceDefinitionList []string
	ConfigMapList                []string
	ServiceList                  []string
	DeploymentList               []string
}

func GetExcutionResource() *Resource {
	var excDep = make([]string, 0)
	var excCM = make([]string, 0)
	var excSV = make([]string, 0)
	for _, str := range excutionList {
		if len(strings.Split(str, "/")) > 0 {
			switch resource := strings.Split(str, "/")[0]; resource {
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
		ConfigMapList:  excCM,
		ServiceList:    excSV,
		DeploymentList: excDep,
	}
}

func GetGUIResource() *Resource {
	var guiDep = make([]string, 0)
	var guiCM = make([]string, 0)
	var guiSV = make([]string, 0)
	for _, str := range guiList {
		if len(strings.Split(str, "/")) > 0 {
			switch resource := strings.Split(str, "/")[0]; resource {
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
		ConfigMapList:  guiCM,
		ServiceList:    guiSV,
		DeploymentList: guiDep,
	}
}

func GetUnInstallResource() *Resource {
	return &Resource{
		CustomResourceDefinitionList: crdList,
		ConfigMapList:                cmList,
		ServiceList:                  svList,
		DeploymentList:               depList,
	}
}

func (asp AlamedaServiceParamter) GetInstallResource() *Resource {
	crd := crdList
	cm := cmList
	sv := svList
	dep := depList
	if asp.GuiFlag {
		cm = append(cm, "ConfigMap/grafana-datasources.yaml")
		sv = append(sv, "Service/alameda-grafanaSV.yaml")
		dep = append(dep, "Deployment/alameda-grafanaDM.yaml")
	}
	if asp.ExcutionFlag {
		sv = append(sv, "Service/admission-controllerSV.yaml")
		dep = append(dep, "Deployment/admission-controllerDM.yaml")
		dep = append(dep, "Deployment/alameda-evictionerDM.yaml")
	}
	return &Resource{
		CustomResourceDefinitionList: crd,
		ConfigMapList:                cm,
		ServiceList:                  sv,
		DeploymentList:               dep,
	}
}

func NewAlamedaServiceParamter(instance *federatoraiv1alpha1.AlamedaService) *AlamedaServiceParamter {
	asp := &AlamedaServiceParamter{
		//AlmedaInstallOrUninstall: instance.Spec.AlmedaInstallOrUninstall,
		NameSpace:             instance.Namespace,
		EnableExecution:       instance.Spec.EnableExecution,
		EnableGUI:             instance.Spec.EnableGUI,
		Version:               instance.Spec.Version,
		PrometheusService:     instance.Spec.PrometheusService,
		PersistentVolumeClaim: instance.Spec.PersistentVolumeClaim,
		GuiFlag:               instance.Spec.EnableExecution,
		ExcutionFlag:          instance.Spec.EnableExecution,
	}
	return asp
}
