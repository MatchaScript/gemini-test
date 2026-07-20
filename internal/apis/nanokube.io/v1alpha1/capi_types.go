package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=nanokubeclusters,scope=Namespaced

type NanokubeCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NanokubeClusterSpec   `json:"spec,omitempty"`
	Status NanokubeClusterStatus `json:"status,omitempty"`
}

type NanokubeClusterSpec struct {
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`
}

type NanokubeClusterStatus struct {
	Initialization NanokubeInitializationStatus `json:"initialization,omitempty"`
	Ready          bool                         `json:"ready,omitempty"`
}

type NanokubeInitializationStatus struct {
	Provisioned             bool `json:"provisioned,omitempty"`
	ControlPlaneInitialized bool `json:"controlPlaneInitialized,omitempty"`
}

// +kubebuilder:object:root=true

type NanokubeClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NanokubeCluster `json:"items"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=nanokubecontrolplanes,scope=Namespaced

type NanokubeControlPlane struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NanokubeControlPlaneSpec   `json:"spec,omitempty"`
	Status NanokubeControlPlaneStatus `json:"status,omitempty"`
}

type NanokubeControlPlaneSpec struct {
	Replicas        *int32                       `json:"replicas,omitempty"`
	Version         string                       `json:"version"`
	MachineTemplate NanokubeControlPlaneTemplate `json:"machineTemplate"`
}

type NanokubeControlPlaneTemplate struct {
	InfrastructureRef corev1.ObjectReference `json:"infrastructureRef"`
}

type NanokubeControlPlaneStatus struct {
	Initialization NanokubeInitializationStatus `json:"initialization,omitempty"`
	Replicas       int32                        `json:"replicas,omitempty"`
	ReadyReplicas  int32                        `json:"readyReplicas,omitempty"`
	Ready          bool                         `json:"ready,omitempty"`
}

// +kubebuilder:object:root=true

type NanokubeControlPlaneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NanokubeControlPlane `json:"items"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=nanokubemachines,scope=Namespaced

type NanokubeMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NanokubeMachineSpec   `json:"spec,omitempty"`
	Status NanokubeMachineStatus `json:"status,omitempty"`
}

type NanokubeMachineSpec struct {
	HostSelector metav1.LabelSelector `json:"hostSelector"`
	ProviderID   string               `json:"providerID,omitempty"`
}

type NanokubeMachineStatus struct {
	InventoryMachineRef string                       `json:"inventoryMachineRef,omitempty"`
	Initialization      NanokubeInitializationStatus `json:"initialization,omitempty"`
	Ready               bool                         `json:"ready,omitempty"`
	Addresses           []clusterv1.MachineAddress   `json:"addresses,omitempty"`
}

// +kubebuilder:object:root=true

type NanokubeMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NanokubeMachine `json:"items"`
}

// +kubebuilder:object:root=true

type NanokubeMachineTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec NanokubeMachineTemplateSpec `json:"spec,omitempty"`
}

type NanokubeMachineTemplateSpec struct {
	Template NanokubeMachineTemplateResource `json:"template"`
}

type NanokubeMachineTemplateResource struct {
	Spec NanokubeMachineSpec `json:"spec"`
}

// +kubebuilder:object:root=true

type NanokubeMachineTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NanokubeMachineTemplate `json:"items"`
}
