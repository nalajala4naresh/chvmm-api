package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VMSetSpec defines the desired state of VMSet.
type VMSetSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of VMSet. Edit vmset_types.go to remove/update

	Replicas *int32                     `json:"replicas,omitempty"`
	Selector metav1.LabelSelector       `json:"selector"`
	Template VirtualMachineTemplateSpec `json:"virtualMachineTemplate"`
	// +optional
	DataVolumeTemplates []DataVolumeTemplate `json:"dataVolumeTemplates"`
	// +optional
	Strategy VMSetStrategy `json:"strategy,omitempty"`
}

type VolumeClaimRetentionPolicy struct {
	WhenDeleted string `json:"whenDeleted,omitempty"` // "Retain" or "Delete"
	WhenScaled  string `json:"whenScaled,omitempty"`  // "Retain" or "Delete"
}

type VMSetStrategy struct {
	VolumeClaimRetentionPolicy *VolumeClaimRetentionPolicy `json:"volumeClaimRetentionPolicy,omitempty"`
	MaxSurge                   *int32                      `json:"maxSurge,omitempty"`
	MaxUnavailable             *int32                      `json:"maxUnavailable,omitempty"`
}

// VMSetStatus defines the observed state of VMSet.
type VMSetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Replicas          int32  `json:"replicas,omitempty"`
	ReadyReplicas     int32  `json:"readyReplicas,omitempty"`
	AvailableReplicas int32  `json:"availableReplicas,omitempty"`
	LabelSelector     string `json:"labelSelector,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.replicas,selectorpath=.status.labelSelector
// +kubebuilder:printcolumn:name="Desired",type="integer",JSONPath=".spec.replicas",description="Number of desired VMs"
// +kubebuilder:printcolumn:name="Current",type="integer",JSONPath=".status.replicas",description="Number of managed VMs"
// +kubebuilder:printcolumn:name="Ready",type="integer",JSONPath=".status.readyReplicas",description="Number of ready VMs"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// VMSet is the Schema for the vmsets API.
type VMSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VMSetSpec   `json:"spec,omitempty"`
	Status VMSetStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VMSetList contains a list of VMSet.
type VMSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VMSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VMSet{}, &VMSetList{})
}
