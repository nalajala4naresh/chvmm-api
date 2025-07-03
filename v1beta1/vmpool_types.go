package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VMPoolSpec defines the desired state of VMPool.
type VMPoolSpec struct {
	Replicas *int32                     `json:"replicas,omitempty"`
	Selector *metav1.LabelSelector      `json:"selector"`
	Template VirtualMachineTemplateSpec `json:"virtualMachineTemplate"`
	// +optional
	DataVolumeTemplates []DataVolumeTemplate `json:"dataVolumeTemplates"`
}

type VirtualMachineTemplateSpec struct {
	// Standard object's metadata
	// +optional
	Metadata TemplateMetadata `json:"metadata,omitempty"`
	// Spec contains the VM specification to be created
	Spec VirtualMachineSpec `json:"spec"`
}

type TemplateMetadata struct {
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
}

// VMRolloutStrategy defines strategies for VM replacement
type VMRolloutStrategy struct {
	// LiveUpdate defines a strategy where VM definitions are updated while running
	// +optional
	LiveUpdate *LiveUpdateStrategy `json:"liveUpdate,omitempty"`
}

// LiveUpdateStrategy defines parameters for in-place VM updates
type LiveUpdateStrategy struct {
	// MaxSurge is the maximum number of VMs that can be scheduled above the desired number
	// +optional
	MaxSurge *intstr.IntOrString `json:"maxSurge,omitempty"`
}

// NameGeneration provides options for configuring resource name generation
type NameGeneration struct {
	// AppendIndexToConfigMapRefs determines if ConfigMap references should be indexed
	// +optional
	AppendIndexToConfigMapRefs bool `json:"appendIndexToConfigMapRefs,omitempty"`

	// AppendIndexToSecretRefs determines if Secret references should be indexed
	// +optional
	AppendIndexToSecretRefs bool `json:"appendIndexToSecretRefs,omitempty"`
}

// VMPoolStatus defines the observed state of VMPool.
type VMPoolStatus struct {
	Replicas          int32  `json:"replicas,omitempty"`
	ReadyReplicas     int32  `json:"readyReplicas,omitempty"`
	AvailableReplicas int32  `json:"availableReplicas,omitempty"`
	ActiveReplicaSet  string `json:"activeReplicaSet,omitempty"`
	LabelSelector     string `json:"labelSelector,omitempty"`
}

// VirtualMachinePoolCondition defines a condition of the pool
type VirtualMachinePoolCondition struct {
	// Type is the type of the condition
	Type VirtualMachinePoolConditionType `json:"type"`

	// Status is the status of the condition
	Status corev1.ConditionStatus `json:"status"`

	// LastProbeTime is the last time we probed the condition
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty"`

	// LastTransitionTime is the last time the condition transitioned
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`

	// Reason is a unique, one-word, CamelCase reason for the condition's last transition
	// +optional
	Reason string `json:"reason,omitempty"`

	// Message is a human readable message indicating details about the last transition
	// +optional
	Message string `json:"message,omitempty"`
}

// VirtualMachinePoolConditionType is the type of VirtualMachinePool condition
type VirtualMachinePoolConditionType string

const (
	// VirtualMachinePoolReplicaFailure indicates that the pool has replica failures
	VirtualMachinePoolReplicaFailure VirtualMachinePoolConditionType = "ReplicaFailure"

	// VirtualMachinePoolReplicaPaused indicates the controller is paused for this pool
	VirtualMachinePoolReplicaPaused VirtualMachinePoolConditionType = "Paused"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.replicas,selectorpath=.status.labelSelector
// +kubebuilder:printcolumn:name="Desired",type="integer",JSONPath=".spec.replicas",description="Number of desired VMs"
// +kubebuilder:printcolumn:name="Current",type="integer",JSONPath=".status.replicas",description="Number of managed VMs"
// +kubebuilder:printcolumn:name="Ready",type="integer",JSONPath=".status.readyReplicas",description="Number of ready VMs"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type VMPool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VMPoolSpec   `json:"spec,omitempty"`
	Status VMPoolStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VMPoolList contains a list of VMPool.
type VMPoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VMPool `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VMPool{}, &VMPoolList{})
}
