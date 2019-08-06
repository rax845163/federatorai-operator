package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type Platform = string

const (
	PlatformOpenshift3_9 Platform = "openshift3.9"
)

// AlamedaServiceSpec defines the desired state of AlamedaService
// +k8s:openapi-gen=true
type AlamedaServiceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	// +kubebuilder:validation:Enum=openshift3.9
	Platform          Platform              `json:"platform,omitempty"`
	EnableExecution   bool                  `json:"enableExecution"`
	EnableGUI         bool                  `json:"enableGui"`
	SelfDriving       bool                  `json:"selfDriving"`
	EnableFedemeter   bool                  `json:"enableFedemeter"`
	Version           string                `json:"version"`
	PrometheusService string                `json:"prometheusService"`
	Storages          []StorageSpec         `json:"storages"`
	ServiceExposures  []ServiceExposureSpec `json:"serviceExposures"`

	//Component Section Schema
	InfluxdbSectionSet            AlamedaComponentSpec `json:"alamedaInfluxdb"`
	GrafanaSectionSet             AlamedaComponentSpec `json:"alamedaGrafana"`
	AlamedaAISectionSet           AlamedaComponentSpec `json:"alamedaAi"`
	AlamedaOperatorSectionSet     AlamedaComponentSpec `json:"alamedaOperator"`
	AlamedaDatahubSectionSet      AlamedaComponentSpec `json:"alamedaDatahub"`
	AlamedaEvictionerSectionSet   AlamedaComponentSpec `json:"alamedaEvictioner"`
	AdmissionControllerSectionSet AlamedaComponentSpec `json:"alamedaAdmissionController"`
	AlamedaRecommenderSectionSet  AlamedaComponentSpec `json:"alamedaRecommender"`
	AlamedaExecutorSectionSet     AlamedaComponentSpec `json:"alamedaExecutor"`
	AlamedaFedemeterSectionSet    AlamedaComponentSpec `json:"fedemeter"`
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
	AccessModes corev1.PersistentVolumeAccessMode `json:"accessMode"`
}

//check StorageStruct
func (storageStruct StorageSpec) StorageIsEmpty() bool {
	if storageStruct.Size != "" && storageStruct.Type == PVC {
		return false
	}
	return true
}

// ServiceExposureType defines the type of the service to be exposed
type ServiceExposureType = string

var (
	// ServiceExposureTypeNodePort represents NodePort type
	ServiceExposureTypeNodePort ServiceExposureType = "NodePort"
)

// ServiceExposureSpec defines the service to be exposed
type ServiceExposureSpec struct {
	Name string `json:"name"`
	// +kubebuilder:validation:Enum=NodePort
	Type     ServiceExposureType `json:"type"`
	NodePort *NodePortSpec       `json:"nodePort,omitempty"`
}

// NodePortSpec defines the ports to be proxied from node to service
type NodePortSpec struct {
	Ports []PortSpec `json:"ports"`
}

// PortSpec defines the service port
type PortSpec struct {
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=65535
	Port int32 `json:"port"`
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=65535
	NodePort int32 `json:"nodePort"`
}

// AlamedaServiceStatus defines the observed state of AlamedaService
// +k8s:openapi-gen=true
type AlamedaServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	CRDVersion AlamedaServiceStatusCRDVersion  `json:"crdversion"`
	Conditions []AlamedaServiceStatusCondition `json:"conditions"`
}

type AlamedaServiceStatusCRDVersion struct {

	// Represents whether any actions on the underlaying managed objects are
	// being performed. Only delete actions will be performed.
	ChangeVersion bool   `json:"-"`
	ScalerVersion string `json:"scalerversion"`
	CRDName       string `json:"crdname"`
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
