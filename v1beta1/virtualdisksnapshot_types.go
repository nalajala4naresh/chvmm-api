package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VirtualDiskSnapshotSpec defines the desired state of VirtualDiskSnapshot
type VirtualDiskSnapshotSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Name of the virtual disk to snapshot
	VirtualDiskName string `json:"virtualDiskName"`

	// Whether to skip the snapshot creation
	SkipSnapshot bool `json:"skipSnapshot,omitempty"`

	// Storage configuration for the snapshot
	Storage *SnapshotStorage `json:"storage,omitempty"`

	// Snapshot destination path
	Destination string `json:"destination,omitempty"`
}

// SnapshotStorage defines the storage configuration for snapshots
type SnapshotStorage struct {
	// Access modes for the snapshot storage
	AccessModes []string `json:"accessModes,omitempty"`

	// Volume mode (Filesystem or Block)
	VolumeMode string `json:"volumeMode,omitempty"`

	// Storage resources
	Resources *SnapshotResources `json:"resources,omitempty"`
}

// SnapshotResources defines the storage resources for snapshots
type SnapshotResources struct {
	// Storage requests
	Requests map[string]string `json:"requests,omitempty"`
}

// VirtualDiskSnapshotStatus defines the observed state of VirtualDiskSnapshot
type VirtualDiskSnapshotStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Phase of the snapshot
	Phase VirtualDiskSnapshotPhase `json:"phase,omitempty"`

	// Snapshot details
	Snapshot *SnapshotDetails `json:"snapshot,omitempty"`

	// Conditions for the snapshot
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Progress of the snapshot operation (0-100)
	Progress int32 `json:"progress,omitempty"`

	// Error message if snapshot failed
	ErrorMessage string `json:"errorMessage,omitempty"`

	// Creation timestamp
	CreatedAt metav1.Time `json:"createdAt,omitempty"`

	// Completion timestamp
	CompletedAt *metav1.Time `json:"completedAt,omitempty"`
}

// SnapshotDetails contains information about the created snapshot
type SnapshotDetails struct {
	// Name of the snapshot
	Name string `json:"name"`

	// Size of the snapshot
	Size string `json:"size,omitempty"`

	// Format of the snapshot
	Format string `json:"format,omitempty"`

	// Path to the snapshot file
	Path string `json:"path,omitempty"`

	// Checksum of the snapshot
	Checksum string `json:"checksum,omitempty"`

	// Whether the snapshot is ready
	Ready bool `json:"ready"`

	// Parent snapshot name for incremental snapshots
	ParentSnapshot string `json:"parentSnapshot,omitempty"`
}

// VirtualDiskSnapshotPhase represents the phase of a virtual disk snapshot
// +kubebuilder:validation:Enum=Pending;InProgress;ReadyToUse;Error;Deleting
type VirtualDiskSnapshotPhase string

const (
	// VirtualDiskSnapshotPending means the snapshot is pending creation
	VirtualDiskSnapshotPending VirtualDiskSnapshotPhase = "Pending"

	// VirtualDiskSnapshotInProgress means the snapshot is being created
	VirtualDiskSnapshotInProgress VirtualDiskSnapshotPhase = "InProgress"

	// VirtualDiskSnapshotReadyToUse means the snapshot is ready for use
	VirtualDiskSnapshotReadyToUse VirtualDiskSnapshotPhase = "ReadyToUse"

	// VirtualDiskSnapshotError means the snapshot creation failed
	VirtualDiskSnapshotError VirtualDiskSnapshotPhase = "Error"

	// VirtualDiskSnapshotDeleting means the snapshot is being deleted
	VirtualDiskSnapshotDeleting VirtualDiskSnapshotPhase = "Deleting"
)

// VirtualDiskSnapshotConditionType represents the type of a virtual disk snapshot condition
type VirtualDiskSnapshotConditionType string

const (
	// VirtualDiskSnapshotReady means the snapshot is ready for use
	VirtualDiskSnapshotReady VirtualDiskSnapshotConditionType = "Ready"

	// VirtualDiskSnapshotCreated means the snapshot was created successfully
	VirtualDiskSnapshotCreated VirtualDiskSnapshotConditionType = "Created"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=vdss

// VirtualDiskSnapshot is the Schema for the virtualdisksnapshots API
type VirtualDiskSnapshot struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualDiskSnapshotSpec   `json:"spec,omitempty"`
	Status VirtualDiskSnapshotStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VirtualDiskSnapshotList contains a list of VirtualDiskSnapshot
type VirtualDiskSnapshotList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualDiskSnapshot `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VirtualDiskSnapshot{}, &VirtualDiskSnapshotList{})
}
