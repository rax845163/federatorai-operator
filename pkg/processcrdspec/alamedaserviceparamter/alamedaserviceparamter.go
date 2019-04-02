package alamedaserviceparamter

import (
	"strings"

	federatoraiv1alpha1 "github.com/containers-ai/federatorai-operator/pkg/apis/federatorai/v1alpha1"
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
	ClusterRoleBinding           []string
	ClusterRole                  []string
	ServiceAccount               []string
	CustomResourceDefinitionList []string
	ConfigMapList                []string
	ServiceList                  []string
	DeploymentList               []string
}

func GetExcutionResource() *Resource {
	var guicrb = make([]string, 0)
	var guicr = make([]string, 0)
	var guisa = make([]string, 0)
	var excDep = make([]string, 0)
	var excCM = make([]string, 0)
	var excSV = make([]string, 0)
	for _, str := range excutionList {
		if len(strings.Split(str, "/")) > 0 {
			switch resource := strings.Split(str, "/")[0]; resource {
			case "ClusterRoleBinding":
				guicrb = append(guicrb, str)
			case "ClusterRole":
				guicr = append(guicr, str)
			case "ServiceAccount":
				guisa = append(guisa, str)
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
		ClusterRoleBinding: guicrb,
		ClusterRole:        guicr,
		ServiceAccount:     guisa,
		ConfigMapList:      excCM,
		ServiceList:        excSV,
		DeploymentList:     excDep,
	}
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
		ClusterRoleBinding: guicrb,
		ClusterRole:        guicr,
		ServiceAccount:     guisa,
		ConfigMapList:      guiCM,
		ServiceList:        guiSV,
		DeploymentList:     guiDep,
	}
}

func GetUnInstallResource() *Resource {
	return &Resource{
		ClusterRoleBinding:           crbList,
		ClusterRole:                  crList,
		ServiceAccount:               saList,
		CustomResourceDefinitionList: crdList,
		ConfigMapList:                cmList,
		ServiceList:                  svList,
		DeploymentList:               depList,
	}
}

func (asp AlamedaServiceParamter) GetInstallResource() *Resource {
	crb := crbList
	cr := crList
	sa := saList
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
		crb = append(crb, "ClusterRoleBinding/alameda-evictionerCRB.yaml")
		crb = append(crb, "ClusterRoleBinding/admission-controllerCRB.yaml")
		cr = append(cr, "ClusterRole/alameda-evictionerCR.yaml")
		cr = append(cr, "ClusterRole/admission-controllerCR.yaml")
		sa = append(sa, "ServiceAccount/alameda-evictionerSA.yaml")
		sa = append(sa, "ServiceAccount/admission-controllerSA.yaml")
		sv = append(sv, "Service/admission-controllerSV.yaml")
		dep = append(dep, "Deployment/admission-controllerDM.yaml")
		dep = append(dep, "Deployment/alameda-evictionerDM.yaml")
	}
	return &Resource{
		ClusterRoleBinding:           crb,
		ClusterRole:                  cr,
		ServiceAccount:               sa,
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
		GuiFlag:               instance.Spec.EnableGUI,
		ExcutionFlag:          instance.Spec.EnableExecution,
	}
	return asp
}
