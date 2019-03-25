package alamedaserviceparamter

import (
	federatoraiv1alpha1 "github.com/containers-ai/federatorai-operator/pkg/apis/federatorai/v1alpha1"
)

type Alamedaserviceparamter struct {
	AlmedaInstallOrUninstall bool
	EnableExecution          bool
	EnableGUI                bool
	Version                  string
	Guicomponent             []string
	Excutioncomponent        []string
}

func NewAlamedaServiceParamter(instance *federatoraiv1alpha1.AlamedaService) *Alamedaserviceparamter {
	asp := &Alamedaserviceparamter{
		AlmedaInstallOrUninstall: instance.Spec.AlmedaInstallOrUninstall,
		EnableExecution:          instance.Spec.EnableExecution,
		EnableGUI:                instance.Spec.EnableGUI,
		Version:                  instance.Spec.Version,
	}
	var guicomponent = make([]string, 4)
	if !instance.Spec.EnableGUI {
		guicomponent = append(guicomponent, "../../assets/PersistentVolumeClaim/my-alamedagrafanaPVC.yaml")
		guicomponent = append(guicomponent, "../../assets/ConfigMap/grafana-datasources.yaml")
		guicomponent = append(guicomponent, "../../assets/Deployment/alameda-grafanaDM.yaml")
		guicomponent = append(guicomponent, "../../assets/Service/alameda-grafanaSV.yaml")

		asp.Guicomponent = guicomponent
	} else {
		asp.Guicomponent = nil
	}
	var excutioncomponent = make([]string, 4)
	if !instance.Spec.EnableExecution {
		excutioncomponent = append(excutioncomponent, "../../assets/ClusterRoleBinding/admission-controllerCRB.yaml")
		excutioncomponent = append(excutioncomponent, "../../assets/ClusterRoleBinding/alameda-evictionerCRB.yaml")
		excutioncomponent = append(excutioncomponent, "../../assets/ClusterRole/admission-controllerCR.yaml")
		excutioncomponent = append(excutioncomponent, "../../assets/ClusterRole/alameda-evictionerCR.yaml")
		excutioncomponent = append(excutioncomponent, "../../assets/Deployment/admission-controllerDM.yaml")
		excutioncomponent = append(excutioncomponent, "../../assets/Deployment/alameda-evictionerDM.yaml")
		excutioncomponent = append(excutioncomponent, "../../assets/Service/admission-controllerSV.yaml")
		excutioncomponent = append(excutioncomponent, "../../assets/ServiceAccount/admission-controllerSA.yaml")
		excutioncomponent = append(excutioncomponent, "../../assets/ServiceAccount/alameda-evictionerSA.yaml")
		asp.Excutioncomponent = excutioncomponent
	} else {
		asp.Excutioncomponent = nil
	}
	return asp
}
