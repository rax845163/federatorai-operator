package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AlamedaServiceSpec defines the desired state of AlamedaService
// +k8s:openapi-gen=true
type AlamedaServiceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	//AlmedaInstallOrUninstall bool   `json:"almedainstalloruninstall"`
	EnableExecution   bool          `json:"enableexecution"`
	EnableGUI         bool          `json:"enablegui"`
	Version           string        `json:"version"`
	PrometheusService string        `json:"prometheusservice"`
	Storages          []StorageSpec `json:"storages"`
	//Component Section Schema
	InfluxdbSectionSet            AlamedaComponentSpec `json:"alameda-influxdb"`
	GrafanaSectionSet             AlamedaComponentSpec `json:"alameda-grafana"`
	AlamedaAISectionSet           AlamedaComponentSpec `json:"alameda-ai"`
	AlamedaOperatorSectionSet     AlamedaComponentSpec `json:"alameda-operator"`
	AlamedaDatahubSectionSet      AlamedaComponentSpec `json:"alameda-datahub"`
	AlamedaEvictionerSectionSet   AlamedaComponentSpec `json:"alameda-evictioner"`
	AdmissionControllerSectionSet AlamedaComponentSpec `json:"alameda-admission-controller"`
	AlamedaRecommenderSectionSet  AlamedaComponentSpec `json:"alameda-recommender"`
}
type AlamedaComponentSpec struct {
	Image              string            `json:"image"`
	Version            string            `json:"version"`
	ImagePullPolicy    corev1.PullPolicy `json:"imagepullpolicy"`
	Storages           []StorageSpec     `json:"storages"`
	BootStrapContainer Imagestruct       `json:"bootstrap"`
}
type Imagestruct struct {
	Image           string            `json:"image"`
	Version         string            `json:"version"`
	ImagePullPolicy corev1.PullPolicy `json:"imagepullpolicy"`
}
type Usage string
type Type string

const (
	Empty     Usage = ""
	Log       Usage = "log"
	Data      Usage = "data"
	PVC       Type  = "pvc"
	Ephemeral Type  = "ephemeral"
)

var (
	PvcUsage = []Usage{Data, Log}
)

type StorageSpec struct {
	Type        Type                              `json:"type"`
	Usage       Usage                             `json:"usage"`
	Size        string                            `json:"size"`
	Class       *string                           `json:"class"`
	AccessModes corev1.PersistentVolumeAccessMode `json:"accessmode"`
}

//check StorageStruct
func (storageStruct StorageSpec) StorageIsEmpty() bool {
	if storageStruct.Size != "" && storageStruct.Type == PVC {
		return false
	}
	return true
}

// AlamedaServiceStatus defines the observed state of AlamedaService
// +k8s:openapi-gen=true
type AlamedaServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html

	Conditions []AlamedaServiceStatusCondition `json:"conditions"`
}

type AlamedaServiceStatusCondition struct {

	// Represents whether any actions on the underlaying managed objects are
	// being performed. Only delete actions will be performed.
	Paused  bool   `json:"paused"`
	Message string `json:"message"`
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
