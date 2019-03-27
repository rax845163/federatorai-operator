package alamedaserviceparamter

import (
	federatoraiv1alpha1 "github.com/containers-ai/federatorai-operator/pkg/apis/federatorai/v1alpha1"
)

type AlamedaServiceParamter struct {
	//AlmedaInstallOrUninstall bool
	EnableExecution       bool
	EnableGUI             bool
	Version               string
	PrometheusService     string
	PersistentVolumeClaim string
	Guicomponent          []string
	Excutioncomponent     []string
}

func NewAlamedaServiceParamter(instance *federatoraiv1alpha1.AlamedaService) *AlamedaServiceParamter {
	asp := &AlamedaServiceParamter{
		//AlmedaInstallOrUninstall: instance.Spec.AlmedaInstallOrUninstall,
		EnableExecution:       instance.Spec.EnableExecution,
		EnableGUI:             instance.Spec.EnableGUI,
		Version:               instance.Spec.Version,
		PrometheusService:     instance.Spec.PrometheusService,
		PersistentVolumeClaim: instance.Spec.PersistentVolumeClaim,
	}
	var guicomponent = make([]string, 0)
	if !instance.Spec.EnableGUI {
		//guicomponent = append(guicomponent, "PersistentVolumeClaim/my-alamedagrafanaPVC.yaml")
		guicomponent = append(guicomponent, "ConfigMap/grafana-datasources.yaml")
		guicomponent = append(guicomponent, "Deployment/alameda-grafanaDM.yaml")
		guicomponent = append(guicomponent, "Service/alameda-grafanaSV.yaml")

		asp.Guicomponent = guicomponent
	} else {
		asp.Guicomponent = nil
	}
	var excutioncomponent = make([]string, 0)
	if !instance.Spec.EnableExecution {
		//excutioncomponent = append(excutioncomponent, "ClusterRoleBinding/admission-controllerCRB.yaml")
		//excutioncomponent = append(excutioncomponent, "ClusterRoleBinding/alameda-evictionerCRB.yaml")
		//excutioncomponent = append(excutioncomponent, "ClusterRole/admission-controllerCR.yaml")
		//excutioncomponent = append(excutioncomponent, "ClusterRole/alameda-evictionerCR.yaml")
		excutioncomponent = append(excutioncomponent, "Deployment/admission-controllerDM.yaml")
		excutioncomponent = append(excutioncomponent, "Deployment/alameda-evictionerDM.yaml")
		excutioncomponent = append(excutioncomponent, "Service/admission-controllerSV.yaml")
		//excutioncomponent = append(excutioncomponent, "ServiceAccount/admission-controllerSA.yaml")
		//excutioncomponent = append(excutioncomponent, "ServiceAccount/alameda-evictionerSA.yaml")
		asp.Excutioncomponent = excutioncomponent
	} else {
		asp.Excutioncomponent = nil
	}
	return asp
}
