package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VMReplicaSetSpec defines the desired state of VMReplicaSet.
type VMReplicaSetSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Replicas *int32                     `json:"replicas,omitempty"`
	Selector metav1.LabelSelector       `json:"selector"`
	Template VirtualMachineTemplateSpec `json:"virtualMachineTemplate"`
	// +optional
	DiskTemplates []DiskTemplate `json:"diskTemplates"`
	// +optional
	Strategy VMReplicaSetStrategy `json:"strategy,omitempty"`
}

type DiskTemplate struct {
	MetaData ObjMeta         `json:"metadata,omitempty"`
	Spec     VirtualDiskSpec `json:"spec"`
}

type ObjMeta struct {
	Name        string            `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	Annotations map[string]string `json:"annotations,omitempty" protobuf:"bytes,12,rep,name=annotations"`
	Labels      map[string]string `json:"labels,omitempty" protobuf:"bytes,11,rep,name=labels"`
}

type VMReplicaSetStrategy struct {
	MaxSurge       *int32 `json:"maxSurge,omitempty"`
	MaxUnavailable *int32 `json:"maxUnavailable,omitempty"`
}

// VMReplicaSetStatus defines the observed state of VMReplicaSet.
type VMReplicaSetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Replicas          int32 `json:"replicas,omitempty"`
	ReadyReplicas     int32 `json:"readyReplicas,omitempty"`
	AvailableReplicas int32 `json:"availableReplicas,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Desired",type="integer",JSONPath=".spec.replicas",description="Number of desired VMs"
// +kubebuilder:printcolumn:name="Current",type="integer",JSONPath=".status.replicas",description="Number of managed VMs"
// +kubebuilder:printcolumn:name="Ready",type="integer",JSONPath=".status.readyReplicas",description="Number of ready VMs"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// VMReplicaSet is the Schema for the vmreplicasets API.
type VMReplicaSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VMReplicaSetSpec   `json:"spec,omitempty"`
	Status VMReplicaSetStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VMReplicaSetList contains a list of VMReplicaSet.
type VMReplicaSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VMReplicaSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VMReplicaSet{}, &VMReplicaSetList{})
}
