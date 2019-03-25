package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AlamedaServiceSpec defines the desired state of AlamedaService
// +k8s:openapi-gen=true
type AlamedaServiceSpec struct {
	AlmedaInstallOrUninstall bool   `json:"almedainstalloruninstall"`
	EnableExecution          bool   `json:"enableexecution"`
	EnableGUI                bool   `json:"enablegui"`
	Version                  string `json:"version"`
	PrometheusService        string `json:"prometheusservice"`
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// AlamedaServiceStatus defines the observed state of AlamedaService
// +k8s:openapi-gen=true
type AlamedaServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AlamedaService is the Schema for the alamedaservices API
// +k8s:openapi-gen=true
type AlamedaService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AlamedaServiceSpec   `json:"spec,omitempty"`
	Status AlamedaServiceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AlamedaServiceList contains a list of AlamedaService
type AlamedaServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AlamedaService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AlamedaService{}, &AlamedaServiceList{})
}
