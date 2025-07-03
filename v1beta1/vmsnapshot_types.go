package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VMSnapShotSpec defines the desired state of VMSnapShot.
type VMSnapShotSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	VmName             string `json:"vm,omitempty"`
	SkipMemorySnapshot bool   `json:"skipMemorySnapshot"`
	Bucket             string `json:"bucket,omitempty"`
	Region             string `json:"region,omitempty"`
	AccessKey          string `json:"accessKey,omitempty"`
	SecretKey          string `json:"secretKey,omitempty"`
	SessionToken       string `json:"sessionToken,omitempty"`
	Destination        string `json:"path,omitempty"`
}

// VMSnapShotStatus defines the observed state of VMSnapShot.
type VMSnapShotStatus struct {
	VMSnapShotStatusPhase VMSnapShotStatusPhase `json:"phase,omitempty"`
	DiskStatus            []DiskStatus          `json:"diskStatus,omitempty"`
	MemoryStatus          MemoryStatus          `json:"memoryStatus,omitempty"`
}

type DiskStatus struct {
	Name               string      `json:"name"`
	Ready              bool        `json:"ready"`
	SnapshotName       string      `json:"snapName,omitempty"`
	Message            string      `json:"message,omitempty"`
	Reason             string      `json:"reason,omitempty"`
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
}

type MemoryStatus struct {
	Name               string      `json:"name"`
	Ready              bool        `json:"ready"`
	Bucket             string      `json:"bucket"`
	ObjectKey          string      `json:"objectKey"`
	Message            string      `json:"message,omitempty"`
	Reason             string      `json:"reason,omitempty"`
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
}

// +kubebuilder:validation:Enum=Created;Failed;Scheduled;InProgress;VMRestoreSpecFailed;VMRestoreSpecCreated
type VMSnapShotStatusPhase string

const (
	VMSnapshotCreated    VMSnapShotStatusPhase = "Created"
	VMSnapshotFailed     VMSnapShotStatusPhase = "Failed"
	VMRestoreSpecFailed  VMSnapShotStatusPhase = "VMRestoreSpecFailed"
	VMRestoreSpecCreated VMSnapShotStatusPhase = "VMRestoreSpecCreated"
	VMSnapShotScheduled  VMSnapShotStatusPhase = "Scheduled"
	VMSnapShotInProgress VMSnapShotStatusPhase = "InProgress"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=vmsnap

// VMSnapShot is the Schema for the vmsnapshots API.
type VMSnapShot struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VMSnapShotSpec   `json:"spec,omitempty"`
	Status VMSnapShotStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VMSnapShotList contains a list of VMSnapShot.
type VMSnapShotList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VMSnapShot `json:"items"`
}
