package v1alpha1

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"

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
	EnableDispatcher  *bool                 `json:"enableDispatcher"`
	SelfDriving       bool                  `json:"selfDriving"`
	Version           string                `json:"version"`
	PrometheusService string                `json:"prometheusService"`
	Storages          []StorageSpec         `json:"storages"`
	ServiceExposures  []ServiceExposureSpec `json:"serviceExposures"`
	EnableWeavescope  bool                  `json:"enableWeavescope"`
	Keycode           KeycodeSpec           `json:"keycode"`
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
	AlamedaWeavescopeSectionSet   AlamedaComponentSpec `json:"alameda-weavescope"`
	AlamedaDispatcherSectionSet   AlamedaComponentSpec `json:"alameda-dispatcher"`
	AlamedaRabbitMQSectionSet     AlamedaComponentSpec `json:"alamedaRabbitMQ"`
	AlamedaAnalyzerSectionSet     AlamedaComponentSpec `json:"alameda-analyzer"`
	AlamedaNotifierSectionSet     AlamedaComponentSpec `json:"alamedaNotifier"`
	FederatoraiAgentSectionSet    AlamedaComponentSpec `json:"federatoraiAgent"`
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

// KeycodeSpec contains data for keycode check
type KeycodeSpec struct {
	// CodeNumber provides user to apply keycode to Federator.ai
	CodeNumber string `json:"codeNumber"`
	// SignatureData provides user to apply signature data which is download from ProphetStor to Federator.ai
	SignatureData string `json:"signatureData"`
}

// KeycodeState defines type of keycode processing state
type KeycodeState string

var (
	// KeycodeStateDefault represents default state
	KeycodeStateDefault KeycodeState
	// KeycodeStateWaitingKeycode represents state in waiting keycode to be filled in
	KeycodeStateWaitingKeycode KeycodeState = "WaitingKeycode"
	// KeycodeStatePollingRegistrationData represents in poll registration data state
	KeycodeStatePollingRegistrationData KeycodeState = "PollingRegistrationData"
	// KeycodeStateWaitingSignatureData represents state waiting user fill in signature data
	KeycodeStateWaitingSignatureData KeycodeState = "WaitingSignatureData"
	// KeycodeStateDone represents state waiting user fill in signature data
	KeycodeStateDone KeycodeState = "Done"
)

// KeycodeStatus contains current keycode information
type KeycodeStatus struct {
	// CodeNumber represents the last keycode user successfully applied
	CodeNumber string `json:"codeNumber"`
	// RegistrationData contains data that user need to send to ProphetStor to activate keycode
	RegistrationData string `json:"registrationData"`
	// State represents the state of keycode processing
	State KeycodeState `json:"state"`
	// LastErrorMessage stores the error message that happend when Federatorai-Operator handled keycode
	LastErrorMessage string `json:"lastErrorMessage"`
	// Summary stores the summary of the keycode
	Summary string `json:"summary"`
}

// AlamedaServiceStatus defines the observed state of AlamedaService
// +k8s:openapi-gen=true
type AlamedaServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	CRDVersion    AlamedaServiceStatusCRDVersion  `json:"crdversion"`
	Conditions    []AlamedaServiceStatusCondition `json:"conditions"`
	KeycodeStatus KeycodeStatus                   `json:"keycodeStatus"`
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

// IsCodeNumberEmpty returns true if keycode is empty
func (as *AlamedaService) IsCodeNumberEmpty() bool {

	if as.Spec.Keycode.CodeNumber == "" {
		return true
	}

	return false
}

// IsCodeNumberUpdated returns true if current keycode is not equal to previous keycode
func (as *AlamedaService) IsCodeNumberUpdated() bool {

	if as.Spec.Keycode.CodeNumber != as.Status.KeycodeStatus.CodeNumber {
		return true
	}

	return false
}

// SetCRDVersion sets crdVersion into AlamedaService's status
func (as *AlamedaService) SetCRDVersion(crdVer AlamedaServiceStatusCRDVersion) {
	as.Status.CRDVersion = crdVer
}

// SetStatusCodeNumber sets codeNumber into AlamedaService's status
func (as *AlamedaService) SetStatusCodeNumber(codeNumber string) {
	as.Status.KeycodeStatus.CodeNumber = codeNumber
}

// SetStatusKeycode sets keycode status into AlamedaService's status
func (as *AlamedaService) SetStatusKeycode(status KeycodeStatus) {
	as.Status.KeycodeStatus = status
}

// SetStatusRegistrationData sets registration data into AlamedaService's status
func (as *AlamedaService) SetStatusRegistrationData(registrationData string) {
	as.Status.KeycodeStatus.RegistrationData = registrationData
}

// SetStatusKeycodeState sets registration data into AlamedaService's status
func (as *AlamedaService) SetStatusKeycodeState(state KeycodeState) {
	as.Status.KeycodeStatus.State = state
}

// SetStatusKeycodeLastErrorMessage sets last error message into AlamedaService's keycode status
func (as *AlamedaService) SetStatusKeycodeLastErrorMessage(msg string) {
	as.Status.KeycodeStatus.LastErrorMessage = msg
}

// SetStatusKeycodeSummary sets keycode summary into AlamedaService's status
func (as *AlamedaService) SetStatusKeycodeSummary(summary string) {
	as.Status.KeycodeStatus.Summary = summary
}

// GetSpecAnnotationWithoutKeycode sets keycode summary into AlamedaService's status
func (as AlamedaService) GetSpecAnnotationWithoutKeycode() (string, error) {
	as.Spec.Keycode = KeycodeSpec{}
	jsonSpec, err := json.Marshal(as.Spec)
	if err != nil {
		return "", err
	}
	return string(jsonSpec), nil
}

// GetPrometheusNamespace returns prometheus running namespace by parsing PrometheusService
func (as AlamedaService) GetPrometheusNamespace() (string, error) {

	ps := as.Spec.PrometheusService
	slashes := "://"
	slashesIndex := strings.Index(ps, slashes)
	domainNameAndPath := ps[slashesIndex+len(slashes):]

	semicolon := ":"
	semicolonIndex := strings.Index(domainNameAndPath, semicolon)
	domainName := domainNameAndPath[:semicolonIndex]

	tokens := strings.Split(domainName, ".")
	if len(tokens) < 2 {
		return "", errors.New("length of slice < 2 ,seperate prometheus service by \",\" ")
	}

	return tokens[1], nil
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
