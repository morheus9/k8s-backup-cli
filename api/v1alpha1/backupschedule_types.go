package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BackupScheduleSpec struct {
	// ⭐️ Твои поля здесь:
	Schedule          string   `json:"schedule"` // Cron schedule
	S3Bucket          string   `json:"s3Bucket"` // S3 bucket
	IncludeNamespaces []string `json:"includeNamespaces,omitempty"`
	// Add labels for CRD
}

type BackupScheduleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	LastBackupTime   *metav1.Time `json:"lastBackupTime,omitempty"`
	LastBackupStatus string       `json:"lastBackupStatus,omitempty"`
	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the BackupSchedule resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	//
	// Standard condition types include:
	// - "Available": the resource is fully functional
	// - "Progressing": the resource is being created or updated
	// - "Degraded": the resource failed to reach or maintain its desired state
	//
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// BackupSchedule is the Schema for the backupschedules API
type BackupSchedule struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of BackupSchedule
	// +required
	Spec BackupScheduleSpec `json:"spec"`

	// status defines the observed state of BackupSchedule
	// +optional
	Status BackupScheduleStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// BackupScheduleList contains a list of BackupSchedule
type BackupScheduleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []BackupSchedule `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BackupSchedule{}, &BackupScheduleList{})
}
