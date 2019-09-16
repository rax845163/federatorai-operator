/*
Copyright 2019 The Alameda Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"fmt"

	"github.com/containers-ai/alameda/operator/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type enableExecution = bool

const (
	defaultEnableExecution = false
)

type alamedaPolicy = string
type NamespacedName = string

const (
	RecommendationPolicySTABLE  alamedaPolicy = "stable"
	RecommendationPolicyCOMPACT alamedaPolicy = "compact"
)

type AlamedaContainer struct {
	Name      string                      `json:"name" protobuf:"bytes,1,opt,name=name"`
	Resources corev1.ResourceRequirements `json:"resources,omitempty" protobuf:"bytes,2,opt,name=resources"`
}

type AlamedaPod struct {
	Namespace  string             `json:"namespace" protobuf:"bytes,1,opt,name=namespace"`
	Name       string             `json:"name" protobuf:"bytes,2,opt,name=name"`
	UID        string             `json:"uid" protobuf:"bytes,3,opt,name=uid"`
	Containers []AlamedaContainer `json:"containers" protobuf:"bytes,4,opt,name=containers"`
}

func (p *AlamedaPod) GetNamespacedName() NamespacedName {
	return utils.GetNamespacedNameKey(p.Namespace, p.Name)
}

type AlamedaResource struct {
	Namespace    string                        `json:"namespace" protobuf:"bytes,1,opt,name=namespace"`
	Name         string                        `json:"name" protobuf:"bytes,2,opt,name=name"`
	UID          string                        `json:"uid" protobuf:"bytes,3,opt,name=uid"`
	Pods         map[NamespacedName]AlamedaPod `json:"pods" protobuf:"bytes,4,opt,name=pods"`
	SpecReplicas *int32                        `json:"specReplicas" protobuf:"varint,5,opt,name=spec_replicas"`
}

type AlamedaController struct {
	Deployments       map[NamespacedName]AlamedaResource `json:"deployments,omitempty" protobuf:"bytes,1,opt,name=deployments"`
	DeploymentConfigs map[NamespacedName]AlamedaResource `json:"deploymentConfigs,omitempty" protobuf:"bytes,2,opt,name=deployment_configs"`
	StatefulSets      map[NamespacedName]AlamedaResource `json:"statefulSets,omitempty" protobuf:"bytes,3,opt,name=stateful_sets"`
}

type AlamedaControllerType int

const (
	DeploymentController       AlamedaControllerType = 1
	DeploymentConfigController AlamedaControllerType = 2
	StatefulSetController      AlamedaControllerType = 3
)

var (
	AlamedaControllerTypeName = map[AlamedaControllerType]string{
		DeploymentController:       "deployment",
		DeploymentConfigController: "deploymentconfig",
		StatefulSetController:      "statefulset",
	}

	K8SKindToAlamedaControllerType = map[string]AlamedaControllerType{
		"Deployment":       DeploymentController,
		"DeploymentConfig": DeploymentConfigController,
		"StatefulSet":      StatefulSetController,
	}
)

type TriggerThreshold struct {
	// +kubebuilder:validation:Pattern=^\d*[1-9]+\d*%$|^\d*[1-9]+\d*\.\d*%$|^\d*\.\d*[1-9]+\d*%$
	CPU string `json:"cpu,omitempty" protobuf:"bytes,1,name=cpu"`
	// +kubebuilder:validation:Pattern=^\d*[1-9]+\d*%$|^\d*[1-9]+\d*\.\d*%$|^\d*\.\d*[1-9]+\d*%$
	Memory string `json:"memory,omitempty" protobuf:"bytes,2,name=memory"`
}

const (
	DefaultTriggerThresholdCPUPercentage    = "10%"
	DefaultTriggerThresholdMemoryPercentage = "10%"
)

func NewDefaultTriggerThreshold() TriggerThreshold {

	return TriggerThreshold{
		CPU:    DefaultTriggerThresholdCPUPercentage,
		Memory: DefaultTriggerThresholdMemoryPercentage,
	}
}

type ExecutionStrategy struct {
	// +kubebuilder:validation:Pattern=^\d*[1-9]+\d*(%?$)$|^\d*[1-9]+\d*\.\d*(%?$)$|^\d*\.\d*[1-9]+\d*(%?$)$
	MaxUnavailable   string                       `json:"maxUnavailable,omitempty" protobuf:"bytes,1,name=max_unavailable"`
	TriggerThreshold *TriggerThreshold            `json:"triggerThreshold,omitempty" protobuf:"bytes,2,name=trigger_threshold"`
	Resources        *corev1.ResourceRequirements `json:"resources,omitempty" protobuf:"bytes,3,name=resources"`
}

const (
	DefaultMaxUnavailablePercentage = "25%"
)

func NewDefaultExecutionStrategy() ExecutionStrategy {

	triggerThreshold := NewDefaultTriggerThreshold()

	return ExecutionStrategy{
		MaxUnavailable:   DefaultMaxUnavailablePercentage,
		TriggerThreshold: &triggerThreshold,
	}
}

type ScalingToolType = string

const (
	ScalingToolTypeVPA     ScalingToolType = "vpa"
	ScalingToolTypeHPA     ScalingToolType = "hpa"
	ScalingToolTypeDefault ScalingToolType = "N/A"
)

type ScalingToolSpec struct {
	// +kubebuilder:validation:Enum=,vpa,hpa,N/A
	Type              string             `json:"type,omitempty" protobuf:"bytes,1,name=type"`
	ExecutionStrategy *ExecutionStrategy `json:"executionStrategy,omitempty" protobuf:"bytes,2,name=execution_strategy"`
}

// AlamedaScalerSpec defines the desired state of AlamedaScaler
// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
type AlamedaScalerSpec struct {
	// Important: Run "make" to regenerate code after modifying this file
	Selector        *metav1.LabelSelector `json:"selector" protobuf:"bytes,1,name=selector"`
	EnableExecution *enableExecution      `json:"enableExecution,omitempty" protobuf:"bytes,2,name=enable_execution"`
	// +kubebuilder:validation:Enum=stable,compact
	Policy                alamedaPolicy   `json:"policy,omitempty" protobuf:"bytes,3,opt,name=policy"`
	CustomResourceVersion string          `json:"customResourceVersion,omitempty" protobuf:"bytes,4,opt,name=custom_resource_version"`
	ScalingTool           ScalingToolSpec `json:"scalingTool,omitempty" protobuf:"bytes,5,opt,name=scaling_tool"`
}

// AlamedaScalerStatus defines the observed state of AlamedaScaler
type AlamedaScalerStatus struct {
	AlamedaController AlamedaController `json:"alamedaController,omitempty" protobuf:"bytes,4,opt,name=alameda_controller"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AlamedaScaler is the Schema for the alamedascalers API
// +k8s:openapi-gen=true
type AlamedaScaler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AlamedaScalerSpec   `json:"spec,omitempty"`
	Status AlamedaScalerStatus `json:"status,omitempty"`
}

func (as *AlamedaScaler) SetDefaultValue() { //this function is set alamedascaler default value
	as.setDefaultEnableExecution()
	as.setDefaultScalingTool()
}

func (as *AlamedaScaler) SetCustomResourceVersion(v string) {
	as.Spec.CustomResourceVersion = v
}

func (as *AlamedaScaler) SetStatusAlamedaController(ac AlamedaController) {
	as.Status.AlamedaController = ac
}

func (as *AlamedaScaler) GenCustomResourceVersion() string {
	v := as.ResourceVersion
	return v
}

// GetMonitoredPods returns pods restoring in AlamedaScaler.Status
func (as *AlamedaScaler) GetMonitoredPods() []*AlamedaPod {
	pods := make([]*AlamedaPod, 0)

	for _, alamedaResource := range as.Status.AlamedaController.Deployments {
		for _, pod := range alamedaResource.Pods {
			cpPod := pod
			pods = append(pods, &cpPod)
		}
	}

	for _, alamedaResource := range as.Status.AlamedaController.DeploymentConfigs {
		for _, pod := range alamedaResource.Pods {
			cpPod := pod
			pods = append(pods, &cpPod)
		}
	}

	for _, alamedaResource := range as.Status.AlamedaController.StatefulSets {
		for _, pod := range alamedaResource.Pods {
			cpPod := pod
			pods = append(pods, &cpPod)
		}
	}

	return pods
}

func (as *AlamedaScaler) GetLabelMapToSetToAlamedaRecommendationLabel() map[string]string {
	m := make(map[string]string)
	m["alamedascaler"] = fmt.Sprintf("%s.%s", as.GetName(), as.GetNamespace())
	return m
}

func (as *AlamedaScaler) GetRequestCPUMilliCores() string {

	cpuMilliCores := ""

	executionStrategy := as.Spec.ScalingTool.ExecutionStrategy
	if executionStrategy != nil {
		if executionStrategy.Resources != nil {
			if executionStrategy.Resources.Requests != nil {
				if executionStrategy.Resources.Requests.Cpu() != nil {
					cpuMilliCores = fmt.Sprintf("%d", executionStrategy.Resources.Requests.Cpu().MilliValue())
				}
			}
		}
	}

	return cpuMilliCores
}

func (as *AlamedaScaler) GetRequestMemoryBytes() string {

	memoryBytes := ""

	executionStrategy := as.Spec.ScalingTool.ExecutionStrategy
	if executionStrategy != nil {
		if executionStrategy.Resources != nil {
			if executionStrategy.Resources.Requests != nil {
				if executionStrategy.Resources.Requests.Memory() != nil {
					memoryBytes = fmt.Sprintf("%d", executionStrategy.Resources.Requests.Memory().Value())
				}
			}
		}
	}

	return memoryBytes
}

func (as *AlamedaScaler) GetLimitCPUMilliCores() string {

	cpuMilliCores := ""

	executionStrategy := as.Spec.ScalingTool.ExecutionStrategy
	if executionStrategy != nil {
		if executionStrategy.Resources != nil {
			if executionStrategy.Resources.Limits != nil {
				if executionStrategy.Resources.Limits.Cpu() != nil {
					cpuMilliCores = fmt.Sprintf("%d", executionStrategy.Resources.Limits.Cpu().MilliValue())
				}
			}
		}
	}

	return cpuMilliCores
}

func (as *AlamedaScaler) GetLimitMemoryBytes() string {

	memoryBytes := ""

	executionStrategy := as.Spec.ScalingTool.ExecutionStrategy
	if executionStrategy != nil {
		if executionStrategy.Resources != nil {
			if executionStrategy.Resources.Limits != nil {
				if executionStrategy.Resources.Limits.Memory() != nil {
					memoryBytes = fmt.Sprintf("%d", executionStrategy.Resources.Limits.Memory().Value())
				}
			}
		}
	}

	return memoryBytes
}

func (as *AlamedaScaler) IsEnableExecution() bool {
	if as.Spec.EnableExecution == nil || *as.Spec.EnableExecution == false {
		return false
	}
	return true
}

func (as *AlamedaScaler) IsScalingToolTypeHPA() bool {
	return as.Spec.ScalingTool.Type == ScalingToolTypeHPA
}

func (as *AlamedaScaler) IsScalingToolTypeVPA() bool {
	return as.Spec.ScalingTool.Type == ScalingToolTypeVPA
}

// HasAlamedaPod returns true if the pod is reocording in AlamedaScaler.Status
func (as *AlamedaScaler) HasAlamedaPod(namespace, name string) bool {

	for _, deployment := range as.Status.AlamedaController.Deployments {
		deploymentNS := deployment.Namespace
		for _, pod := range deployment.Pods {
			if deploymentNS == namespace && pod.Name == name {
				return true
			}
		}
	}
	for _, deploymentConfig := range as.Status.AlamedaController.DeploymentConfigs {
		deploymentConfigNS := deploymentConfig.Namespace
		for _, pod := range deploymentConfig.Pods {
			if deploymentConfigNS == namespace && pod.Name == name {
				return true
			}
		}
	}
	for _, statefulSet := range as.Status.AlamedaController.StatefulSets {
		statefulSetNS := statefulSet.Namespace
		for _, pod := range statefulSet.Pods {
			if statefulSetNS == namespace && pod.Name == name {
				return true
			}
		}
	}

	return false
}

func (as *AlamedaScaler) setDefaultEnableExecution() {
	if as.Spec.EnableExecution == nil {
		copyDefaultEnableExecution := defaultEnableExecution
		as.Spec.EnableExecution = &copyDefaultEnableExecution
	}
}

func (as *AlamedaScaler) setDefaultScalingTool() {

	if as.Spec.ScalingTool.Type == "" {
		as.Spec.ScalingTool.Type = ScalingToolTypeDefault
	}

	if as.Spec.ScalingTool.Type == ScalingToolTypeVPA {
		if as.Spec.ScalingTool.ExecutionStrategy == nil {
			as.setDefaultExecutionStrategy()
		}
		if as.Spec.ScalingTool.ExecutionStrategy.MaxUnavailable == "" || as.Spec.ScalingTool.ExecutionStrategy.MaxUnavailable == "0" || as.Spec.ScalingTool.ExecutionStrategy.MaxUnavailable == "0%" {
			as.Spec.ScalingTool.ExecutionStrategy.MaxUnavailable = DefaultMaxUnavailablePercentage
		}

		if as.Spec.ScalingTool.ExecutionStrategy.TriggerThreshold == nil {
			as.setDefaultTriggerThreshold()
		}
		if as.Spec.ScalingTool.ExecutionStrategy.TriggerThreshold.CPU == "" {
			as.Spec.ScalingTool.ExecutionStrategy.TriggerThreshold.CPU = DefaultTriggerThresholdCPUPercentage
		}
		if as.Spec.ScalingTool.ExecutionStrategy.TriggerThreshold.Memory == "" {
			as.Spec.ScalingTool.ExecutionStrategy.TriggerThreshold.Memory = DefaultTriggerThresholdMemoryPercentage
		}
	}
}

func (as *AlamedaScaler) setDefaultExecutionStrategy() {
	defaultExecutionStrategy := NewDefaultExecutionStrategy()
	as.Spec.ScalingTool.ExecutionStrategy = &defaultExecutionStrategy
}

func (as *AlamedaScaler) setDefaultTriggerThreshold() {
	defaultTriggerThreshold := NewDefaultTriggerThreshold()
	as.Spec.ScalingTool.ExecutionStrategy.TriggerThreshold = &defaultTriggerThreshold
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AlamedaScalerList contains a list of AlamedaScaler
type AlamedaScalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AlamedaScaler `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AlamedaScaler{}, &AlamedaScalerList{})
}
