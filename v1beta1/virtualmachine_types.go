package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VirtualMachineSpec defines the desired state of VirtualMachine
type VirtualMachineSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of VirtualMachine. Edit virtualmachine_types.go to remove/update
	NodeSelector map[string]string           `json:"nodeSelector,omitempty"`
	Tolerations  []corev1.Toleration         `json:"tolerations,omitempty"`
	Affinity     *corev1.Affinity            `json:"affinity,omitempty"`
	Resources    corev1.ResourceRequirements `json:"resources,omitempty"`

	LivenessProbe   *corev1.Probe `json:"livenessProbe,omitempty"`
	ReadinessProbes *corev1.Probe `json:"readinessProbe,omitempty"`

	RunPolicy RunPolicy `json:"runPolicy,omitempty"`
	Instance  Instance  `json:"instance"`
	Volumes   []*Volume `json:"volumes,omitempty"`
	Networks  []Network `json:"networks,omitempty"`
}

type MemoryPath struct {
	Bucket string `json:"bucket,omitempty"`
	Key    string `json:"key,omitempty"`
}

type RunPolicy string

const (
	RunAlways               RunPolicy = "RunAlways"
	RunPolicyRerunOnFailure RunPolicy = "RerunOnFailure"
	RunPolicyOnce           RunPolicy = "Once"
	RunPolicyManual         RunPolicy = "Manual"
	RunPolicyHalted         RunPolicy = "Halted"
)

type Instance struct {
	CPU         CPU          `json:"cpu,omitempty"`
	Memory      Memory       `json:"memory,omitempty"`
	Kernel      *Kernel      `json:"kernel,omitempty"`
	Pmem        []Pmem       `json:"pmem,omitempty"`
	Disks       []Disk       `json:"disks,omitempty"`
	FileSystems []FileSystem `json:"fileSystems,omitempty"`
	Interfaces  []Interface  `json:"interfaces,omitempty"`
}

type Pmem struct {
	Name string `json:"name,omitempty"`
}

type CPU struct {
	Sockets               uint32 `json:"sockets,omitempty"`
	CoresPerSocket        uint32 `json:"coresPerSocket,omitempty"`
	DedicatedCPUPlacement bool   `json:"dedicatedCPUPlacement,omitempty"`
}

type Memory struct {
	Size      resource.Quantity `json:"size,omitempty"`
	Hugepages *HugePages        `json:"hugepages,omitempty"`
}

type HugePages struct {
	// +kubebuilder:default="1Gi"
	// +kubebuilder:validation:Enum="2Mi";"1Gi"
	PageSize string `json:"pageSize,omitempty"`
}

type Kernel struct {
	Image           string            `json:"image"`
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
	Cmdline         string            `json:"cmdline"`
}

type Initrd struct {
	Image           string            `json:"image"`
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
}

type Disk struct {
	Name     string `json:"name"`
	ReadOnly *bool  `json:"readOnly,omitempty"`
}

type FileSystem struct {
	Name string `json:"name"`
}

type Interface struct {
	Name                   string `json:"name"`
	MAC                    string `json:"mac,omitempty"`
	InterfaceBindingMethod `json:",inline"`
}

type InterfaceBindingMethod struct {
	Bridge     *InterfaceBridge     `json:"bridge,omitempty"`
	Masquerade *InterfaceMasquerade `json:"masquerade,omitempty"`
	SRIOV      *InterfaceSRIOV      `json:"sriov,omitempty"`
	VDPA       *InterfaceVDPA       `json:"vdpa,omitempty"`
	VhostUser  *InterfaceVhostUser  `json:"vhostUser,omitempty"`
}

type InterfaceBridge struct {
}

type InterfaceMasquerade struct {
	CIDR string `json:"cidr,omitempty"`
}

type InterfaceSRIOV struct {
}

type InterfaceVDPA struct {
	NumQueues int  `json:"numQueues,omitempty"`
	IOMMU     bool `json:"iommu,omitempty"`
}

type InterfaceVhostUser struct {
}

type Volume struct {
	Name         string `json:"name"`
	VolumeSource `json:",inline"`
}

func (v *Volume) IsHotpluggable() bool {
	return v.PersistentVolumeClaim != nil && v.PersistentVolumeClaim.Hotpluggable ||
		v.DataVolume != nil && v.DataVolume.Hotpluggable
}

func (v *Volume) PVCName() string {
	switch {
	case v.PersistentVolumeClaim != nil:
		return v.PersistentVolumeClaim.ClaimName
	case v.DataVolume != nil:
		return v.DataVolume.VolumeName
	case v.VirtualDisk != nil:
		return v.VirtualDisk.VirtualDiskName
	default:
		return ""
	}
}

type VolumeSource struct {
	ContainerDisk         *ContainerDiskVolumeSource         `json:"containerDisk,omitempty"`
	CloudInit             *CloudInitVolumeSource             `json:"cloudInit,omitempty"`
	ContainerRootfs       *ContainerRootfsVolumeSource       `json:"containerRootfs,omitempty"`
	PersistentVolumeClaim *PersistentVolumeClaimVolumeSource `json:"persistentVolumeClaim,omitempty"`
	DataVolume            *DataVolumeVolumeSource            `json:"dataVolume,omitempty"`
	MemorySnapshot        *MemorySnapshotSource              `json:"memorySnapshot,omitempty"`
	VirtualDisk           *VirtualDiskVMSource               `json:"virtualDisk,omitempty"`
}

type VirtualDiskVMSource struct {
	Hotpluggable    bool   `json:"hotpluggable,omitempty"`
	VirtualDiskName string `json:"virtualDiskName"`
}

type ContainerDiskVolumeSource struct {
	Image           string            `json:"image"`
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
}

type MemorySnapshotSource struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}

type CloudInitVolumeSource struct {
	UserData              string `json:"userData,omitempty"`
	UserDataBase64        string `json:"userDataBase64,omitempty"`
	UserDataSecretName    string `json:"userDataSecretName,omitempty"`
	NetworkData           string `json:"networkData,omitempty"`
	NetworkDataBase64     string `json:"networkDataBase64,omitempty"`
	NetworkDataSecretName string `json:"networkDataSecretName,omitempty"`
}

type ContainerRootfsVolumeSource struct {
	Image           string            `json:"image"`
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
	Size            resource.Quantity `json:"size"`
}

type PersistentVolumeClaimVolumeSource struct {
	Hotpluggable bool   `json:"hotpluggable,omitempty"`
	ClaimName    string `json:"claimName"`
}

type DataVolumeVolumeSource struct {
	Hotpluggable bool   `json:"hotpluggable,omitempty"`
	VolumeName   string `json:"volumeName"`
}

type Network struct {
	Name          string `json:"name"`
	NetworkSource `json:",inline"`
}

type NetworkSource struct {
	Pod    *PodNetworkSource    `json:"pod,omitempty"`
	Multus *MultusNetworkSource `json:"multus,omitempty"`
}

type PodNetworkSource struct {
}

type MultusNetworkSource struct {
	NetworkName string `json:"networkName"`
}

type VirtualMachineStatus struct {
	Phase        VirtualMachinePhase            `json:"phase,omitempty"`
	VMPodName    string                         `json:"vmPodName,omitempty"`
	VMPodUID     types.UID                      `json:"vmPodUID,omitempty"`
	NodeName     string                         `json:"nodeName,omitempty"`
	PowerAction  VirtualMachinePowerAction      `json:"powerAction,omitempty"`
	Migration    *VirtualMachineStatusMigration `json:"migration,omitempty"`
	Conditions   []metav1.Condition             `json:"conditions,omitempty"`
	VolumeStatus []VolumeStatus                 `json:"volumeStatus,omitempty"`
}

// +kubebuilder:validation:Enum=Pending;Scheduling;Scheduled;Running;Succeeded;Failed;Unknown;SnapshotInprogress;VMRestoreInProgress;Paused;Resumed

type VirtualMachinePhase string

const (
	VirtualMachinePending            VirtualMachinePhase = "Pending"
	VirtualMachineScheduling         VirtualMachinePhase = "Scheduling"
	VirtualMachineScheduled          VirtualMachinePhase = "Scheduled"
	VirtualMachineRunning            VirtualMachinePhase = "Running"
	VirtualMachineSucceeded          VirtualMachinePhase = "Succeeded"
	VirtualMachineFailed             VirtualMachinePhase = "Failed"
	VirtualMachineUnknown            VirtualMachinePhase = "Unknown"
	VirtualMachinePaused             VirtualMachinePhase = "Paused"
	VirtualMachineRestoreInProgress  VirtualMachinePhase = "VMRestoreInProgress"
	VirtualMachineResumed            VirtualMachinePhase = "Resumed"
	VirtualMachineSnapShotInprogress VirtualMachinePhase = "SnapshotInprogress"
)

// +kubebuilder:validation:Enum=PowerOn;PowerOff;Shutdown;Reset;Reboot

type VirtualMachinePowerAction string

const (
	VirtualMachinePowerOn  VirtualMachinePowerAction = "PowerOn"
	VirtualMachinePowerOff VirtualMachinePowerAction = "PowerOff"
	VirtualMachineShutdown VirtualMachinePowerAction = "Shutdown"
	VirtualMachineReset    VirtualMachinePowerAction = "Reset"
	VirtualMachineReboot   VirtualMachinePowerAction = "Reboot"
)

type VirtualMachineStatusMigration struct {
	UID                types.UID                    `json:"uid,omitempty"`
	Phase              VirtualMachineMigrationPhase `json:"phase,omitempty"`
	TargetNodeName     string                       `json:"targetNodeName,omitempty"`
	TargetNodeIP       string                       `json:"targetNodeIP,omitempty"`
	TargetNodePort     int                          `json:"targetNodePort,omitempty"`
	TargetVMPodName    string                       `json:"targetVMPodName,omitempty"`
	TargetVMPodUID     types.UID                    `json:"targetVMPodUID,omitempty"`
	TargetVolumePodUID types.UID                    `json:"targetVolumePodUID,omitempty"`
}

type VirtualMachineConditionType string

const (
	VirtualMachineMigratable VirtualMachineConditionType = "Migratable"
	VirtualMachineReady      VirtualMachineConditionType = "Ready"
)

type VolumeStatus struct {
	Name          string               `json:"name"`
	Phase         VolumePhase          `json:"phase,omitempty"`
	HotplugVolume *HotplugVolumeStatus `json:"hotplugVolume,omitempty"`
}

type VolumePhase string

const (
	VolumePending        VolumePhase = "Pending"
	VolumeAttachedToNode VolumePhase = "AttachedToNode"
	VolumeMountedToPod   VolumePhase = "MountedToPod"
	VolumeReady          VolumePhase = "Ready"
	VolumeDetaching      VolumePhase = "Detaching"
)

type HotplugVolumeStatus struct {
	VolumePodName string    `json:"volumePodName,omitempty"`
	VolumePodUID  types.UID `json:"volumePodUID,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=vm
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="VM phase"
// +kubebuilder:printcolumn:name="Node",type="string",JSONPath=".status.nodeName",description="VM phase"
// +kubebuilder:printcolumn:name="Power",type="string",JSONPath=".status.powerAction",description="VM phase"
// VirtualMachine is the Schema for the virtualmachines API
type VirtualMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualMachineSpec   `json:"spec,omitempty"`
	Status VirtualMachineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VirtualMachineList contains a list of VirtualMachine
type VirtualMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VirtualMachine{}, &VirtualMachineList{}, &VirtualMachineMigration{},
		&VirtualMachineMigrationList{}, &VMSnapShot{}, &VMSnapShotList{},
		&VMRestoreSpec{}, &VMRestoreSpecList{},
		&VMRollback{}, &VMRollbackList{})
}
