package v1beta1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=vmm
// +kubebuilder:printcolumn:name="VM",type=string,JSONPath=`.spec.vmName`
// +kubebuilder:printcolumn:name="Source",type=string,JSONPath=`.status.sourceNodeName`
// +kubebuilder:printcolumn:name="Target",type=string,JSONPath=`.status.targetNodeName`
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.phase`

type VirtualMachineMigration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualMachineMigrationSpec   `json:"spec,omitempty"`
	Status VirtualMachineMigrationStatus `json:"status,omitempty"`
}

type VirtualMachineMigrationSpec struct {
	VMName string `json:"vmName"`
}

type VirtualMachineMigrationStatus struct {
	Phase          VirtualMachineMigrationPhase `json:"phase,omitempty"`
	SourceNodeName string                       `json:"sourceNodeName,omitempty"`
	TargetNodeName string                       `json:"targetNodeName,omitempty"`
}

// +kubebuilder:validation:Enum=Pending;Scheduling;Scheduled;TargetReady;Running;Sent;Succeeded;Failed

type VirtualMachineMigrationPhase string

const (
	VirtualMachineMigrationPending     VirtualMachineMigrationPhase = "Pending"
	VirtualMachineMigrationScheduling  VirtualMachineMigrationPhase = "Scheduling"
	VirtualMachineMigrationScheduled   VirtualMachineMigrationPhase = "Scheduled"
	VirtualMachineMigrationTargetReady VirtualMachineMigrationPhase = "TargetReady"
	VirtualMachineMigrationRunning     VirtualMachineMigrationPhase = "Running"
	VirtualMachineMigrationSent        VirtualMachineMigrationPhase = "Sent"
	VirtualMachineMigrationSucceeded   VirtualMachineMigrationPhase = "Succeeded"
	VirtualMachineMigrationFailed      VirtualMachineMigrationPhase = "Failed"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type VirtualMachineMigrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []VirtualMachineMigration `json:"items"`
}
