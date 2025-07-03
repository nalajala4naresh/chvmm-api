package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VMRollbackSpec defines the desired state of VMRollback.
type VMRollbackSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of VMRollback. Edit vmrollback_types.go to remove/update
	Snapshot string `json:"snapshot,omitempty"`
}

type VMRollbackPhase string

const (
	RollbackScheduled   VMRollbackPhase = "Scheduled"
	RollbackSucceeded   VMRollbackPhase = "Succeeded"
	RollbackFailed      VMRollbackPhase = "Failed"
	SnapShotToPVCFailed VMRollbackPhase = "SnapShotToPVCFailed"
)

// VMRollbackStatus defines the observed state of VMRollback.
type VMRollbackStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Phase VMRollbackPhase `json:"phase,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// VMRollback is the Schema for the vmrollbacks API.
type VMRollback struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VMRollbackSpec   `json:"spec,omitempty"`
	Status VMRollbackStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VMRollbackList contains a list of VMRollback.
type VMRollbackList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VMRollback `json:"items"`
}
