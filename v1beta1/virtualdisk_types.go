package v1beta1

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VirtualDiskSpec defines the desired state of VirtualDisk
type VirtualDiskSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Storage configuration
	Storage *VirtualDiskStorage `json:"storage,omitempty"`

	// Source configuration
	Source *VirtualDiskSource `json:"source,omitempty"`
}

// VirtualDiskStorage defines the storage configuration for a virtual disk
type VirtualDiskStorage struct {
	// Access modes for the disk
	AccessModes []string `json:"accessModes,omitempty"`

	// Volume mode (Filesystem or Block)
	VolumeMode string `json:"volumeMode,omitempty"`

	// Storage resources
	Resources *VirtualDiskResources `json:"resources,omitempty"`
}

// VirtualDiskResources defines the storage resources
type VirtualDiskResources struct {
	// Storage requests
	Requests map[string]resource.Quantity `json:"requests,omitempty"`
}

// VirtualDiskSource defines the source configuration for a virtual disk
type VirtualDiskSource struct {
	// Registry source
	Registry *RegistrySource `json:"registry,omitempty"`

	// HTTP source
	HTTP *HTTPSource `json:"http,omitempty"`

	DiskSnapshot *DiskSnapshotSource `json:"diskSnapshot,omitempty"`

	Empty *EmptySource `json:"empty,omitempty"`
}

// RegistrySource defines a registry disk source
type RegistrySource struct {
	// Registry URL
	URL string `json:"url"`

	// Pull method
	PullMethod string `json:"pullMethod,omitempty"`
}

type EmptySource struct {
}

type DiskSnapshotSource struct {
	DiskSnapshotName string `json:"snapshotName"`
}

// HTTPSource defines an HTTP disk source
type HTTPSource struct {
	// HTTP URL
	URL string `json:"url"`
}

// VirtualDiskStatus defines the observed state of VirtualDisk
type VirtualDiskStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Phase of the disk
	Phase VirtualDiskPhase `json:"phase,omitempty"`

	// Conditions for the disk
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Capacity of the disk
	Capacity resource.Quantity `json:"capacity,omitempty"`

	// Used space on the disk
	Used resource.Quantity `json:"used,omitempty"`

	// Available space on the disk
	Available resource.Quantity `json:"available,omitempty"`

	// Disk path
	Path string `json:"path,omitempty"`

	// Whether the disk is attached
	Attached bool `json:"attached,omitempty"`

	// Name of the VM the disk is attached to
	AttachedTo string `json:"attachedTo,omitempty"`
}

// VirtualDiskPhase represents the phase of a virtual disk
// +kubebuilder:validation:Enum=Pending;Bound;Ready;Failed;Terminating
type VirtualDiskPhase string

const (
	// VirtualDiskPending means the disk is pending creation
	VirtualDiskPending VirtualDiskPhase = "Pending"

	// VirtualDiskBound means the disk is bound and ready for use
	VirtualDiskBound VirtualDiskPhase = "Bound"

	VirtualDiskReady VirtualDiskPhase = "Ready"

	// VirtualDiskFailed means the disk creation failed
	VirtualDiskFailed VirtualDiskPhase = "Failed"

	// VirtualDiskTerminating means the disk is being deleted
	VirtualDiskTerminating VirtualDiskPhase = "Terminating"
)

// VirtualDiskConditionType represents the type of a virtual disk condition
type VirtualDiskConditionType string

const (
	// VirtualDiskReady means the disk is ready for use
	VirtualDiskReadyCondition VirtualDiskConditionType = "Ready"

	// VirtualDiskBound means the disk is bound to a storage backend
	VirtualDiskBoundCondition VirtualDiskConditionType = "Bound"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=vdisk

// VirtualDisk is the Schema for the virtualdisks API
type VirtualDisk struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualDiskSpec   `json:"spec,omitempty"`
	Status VirtualDiskStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VirtualDiskList contains a list of VirtualDisk
type VirtualDiskList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualDisk `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VirtualDisk{}, &VirtualDiskList{})
}
