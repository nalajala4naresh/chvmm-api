package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VMRestoreSpecSpec defines the desired state of VMRestoreSpec.
type VMRestoreSpecSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of VMRestoreSpec. Edit vmrestorespec_types.go to remove/update
	VMName string             `json:"vmName,omitempty"`
	VMSpec VirtualMachineSpec `json:"vmSpec,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// VMRestoreSpec is the Schema for the vmrestorespecs API.
type VMRestoreSpec struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec VMRestoreSpecSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// VMRestoreSpecList contains a list of VMRestoreSpec.
type VMRestoreSpecList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VMRestoreSpec `json:"items"`
}
