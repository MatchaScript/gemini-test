package kubeadm

import (
	"fmt"

	kubeletconfigv1beta1 "k8s.io/kubelet/config/v1beta1"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	"k8s.io/kubernetes/cmd/kubeadm/app/componentconfigs"
	"k8s.io/utils/ptr"
)

// KubeletResolverConfig is pinned explicitly so rendered kubelet
// config bytes never depend on the machine doing the rendering:
// kubeadm's own defaulting (componentconfigs' mutateResolverConfig)
// probes whether systemd-resolved is active on the local host and
// fills resolvConf from the answer — but only when the field is nil.
// Nodes always run systemd-resolved (homelab bootc images), so the
// resolved stub path is the correct value for every node, and setting
// it up front makes the host probe irrelevant.
const KubeletResolverConfig = "/run/systemd/resolve/resolv.conf"

// PinKubeletResolverConfig sets ResolverConfig on cfg's kubelet
// component config to KubeletResolverConfig. Every kubelet-config
// writer (e.g. render.KubeletConfig) must call this before
// kubelet.WriteConfigToDisk so outputs stay byte-identical on any host,
// systemd-resolved or not.
func PinKubeletResolverConfig(cfg *kubeadmapi.ClusterConfiguration) error {
	kubeletCfg, ok := cfg.ComponentConfigs[componentconfigs.KubeletGroup]
	if !ok {
		return fmt.Errorf("no kubelet component config found")
	}
	kc, ok := kubeletCfg.Get().(*kubeletconfigv1beta1.KubeletConfiguration)
	if !ok {
		return fmt.Errorf("unexpected kubelet component config type %T", kubeletCfg.Get())
	}
	kc.ResolverConfig = ptr.To(KubeletResolverConfig)
	kubeletCfg.Set(kc)
	return nil
}
