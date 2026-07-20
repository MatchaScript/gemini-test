package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=clusters,scope=Namespaced

// Cluster represents a nanokube cluster specification.
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterSpec   `json:"spec,omitempty"`
	Status ClusterStatus `json:"status,omitempty"`
}

type ClusterSpec struct {
	KubernetesVersion    string                `json:"kubernetesVersion"`
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`
	Networking           NetworkingSpec        `json:"networking,omitempty"`
	ControlPlane         ControlPlaneSpec      `json:"controlPlane,omitempty"`
	NodeConfigDefaults   NodeConfigSpec        `json:"nodeConfigDefaults,omitempty"`
}

type NetworkingSpec struct {
	PodCIDR     string `json:"podCIDR,omitempty"`
	ServiceCIDR string `json:"serviceCIDR,omitempty"`
	DNSDomain   string `json:"dnsDomain,omitempty"`
}

type ControlPlaneSpec struct {
	APIServer    APIServerSpec `json:"apiServer,omitempty"`
	NodePoolSpec *NodePoolSpec `json:"nodePoolSpec,omitempty"`
	NodePoolName string        `json:"nodePoolName,omitempty"`
}

type APIServerSpec struct {
	ExtraArgs     map[string]string `json:"extraArgs,omitempty"`
	CertExtraSANs []string          `json:"certExtraSANs,omitempty"`
}

type ClusterStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	Versions   map[string]int32   `json:"versions,omitempty"`
}

// +kubebuilder:object:root=true

type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cluster `json:"items"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=nodepools,scope=Namespaced

type NodePool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodePoolSpec   `json:"spec,omitempty"`
	Status NodePoolStatus `json:"status,omitempty"`
}

type NodePoolSpec struct {
	ClusterName       string               `json:"clusterName"`
	Selector          metav1.LabelSelector `json:"selector"`
	KubernetesVersion string               `json:"kubernetesVersion,omitempty"`
	NodeLabels        map[string]string    `json:"nodeLabels,omitempty"`
	Taints            []corev1.Taint       `json:"taints,omitempty"`
	UpgradeStrategy   UpgradeStrategySpec  `json:"upgradeStrategy,omitempty"`
	Config            NodeConfigSpec       `json:"config,omitempty"`
}

type UpgradeStrategySpec struct {
	MaxUnavailable int32 `json:"maxUnavailable,omitempty"`
}

type NodeConfigSpec struct {
	Kubelet KubeletConfigSpec `json:"kubelet,omitempty"`
	Files   []FileSourceSpec  `json:"files,omitempty"`
}

type KubeletConfigSpec struct {
	MaxPods int32 `json:"maxPods,omitempty"`
}

type FileSourceSpec struct {
	Path   string            `json:"path"`
	Mode   string            `json:"mode,omitempty"`
	Source FileSourceRefSpec `json:"source"`
}

type FileSourceRefSpec struct {
	Inline       string                       `json:"inline,omitempty"`
	ConfigMapRef *corev1.ConfigMapKeySelector `json:"configMapRef,omitempty"`
	SecretRef    *corev1.SecretKeySelector    `json:"secretRef,omitempty"`
}

type NodePoolStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	ReadyNodes int32              `json:"readyNodes,omitempty"`
	Versions   map[string]int32   `json:"versions,omitempty"`
}

// +kubebuilder:object:root=true

type NodePoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodePool `json:"items"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=nodeconfigs,scope=Namespaced

type NodeConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodeConfigObjectSpec `json:"spec,omitempty"`
	Status NodeConfigStatus     `json:"status,omitempty"`
}

type NodeConfigObjectSpec struct {
	ClusterName  string               `json:"clusterName"`
	NodeSelector metav1.LabelSelector `json:"nodeSelector"`
	Config       NodeConfigSpec       `json:"config,omitempty"`
}

type NodeConfigStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true

type NodeConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeConfig `json:"items"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=inventorymachines,scope=Cluster

type InventoryMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InventoryMachineSpec   `json:"spec,omitempty"`
	Status InventoryMachineStatus `json:"status,omitempty"`
}

type InventoryMachineSpec struct {
	Address             string                  `json:"address"`
	K8sIP               string                  `json:"k8sIP,omitempty"`
	AgentPort           int32                   `json:"agentPort,omitempty"`
	CredentialSecretRef *corev1.SecretReference `json:"credentialSecretRef,omitempty"`
	NodeLabels          map[string]string       `json:"nodeLabels,omitempty"`
	Maintenance         bool                    `json:"maintenance,omitempty"`
	ProviderID          string                  `json:"providerID,omitempty"`
}

type InventoryMachineStatus struct {
	Provisioned       bool               `json:"provisioned,omitempty"`
	MachineRef        string             `json:"machineRef,omitempty"`
	AppliedRevision   string             `json:"appliedRevision,omitempty"`
	BootedImageDigest string             `json:"bootedImageDigest,omitempty"`
	LastBootHealth    string             `json:"lastBootHealth,omitempty"`
	Conditions        []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true

type InventoryMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []InventoryMachine `json:"items"`
}
