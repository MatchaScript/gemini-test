package operator

import (
	"reflect"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	nanokubev1alpha1 "github.com/MatchaScript/nanokube/internal/apis/nanokube.io/v1alpha1"
)

func TestFoldNodeConfig_PrecedenceOrder(t *testing.T) {
	clusterDefaults := nanokubev1alpha1.NodeConfigSpec{
		Kubelet: nanokubev1alpha1.KubeletConfigSpec{MaxPods: 100},
		Files: []nanokubev1alpha1.FileSourceSpec{
			{Path: "/etc/kubernetes/audit.yaml", Mode: "0644", Source: nanokubev1alpha1.FileSourceRefSpec{Inline: "default"}},
		},
	}

	poolConfig := nanokubev1alpha1.NodeConfigSpec{
		Kubelet: nanokubev1alpha1.KubeletConfigSpec{MaxPods: 250},
		Files: []nanokubev1alpha1.FileSourceSpec{
			{Path: "/etc/kubernetes/audit.yaml", Mode: "0644", Source: nanokubev1alpha1.FileSourceRefSpec{Inline: "pool"}},
			{Path: "/etc/kubernetes/pool.conf", Mode: "0600", Source: nanokubev1alpha1.FileSourceRefSpec{Inline: "pool-file"}},
		},
	}

	nodeConfigs := []nanokubev1alpha1.NodeConfig{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "b-override"},
			Spec: nanokubev1alpha1.NodeConfigObjectSpec{
				Config: nanokubev1alpha1.NodeConfigSpec{
					Kubelet: nanokubev1alpha1.KubeletConfigSpec{MaxPods: 500},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "a-override"},
			Spec: nanokubev1alpha1.NodeConfigObjectSpec{
				Config: nanokubev1alpha1.NodeConfigSpec{
					Kubelet: nanokubev1alpha1.KubeletConfigSpec{MaxPods: 300},
				},
			},
		},
	}

	got := FoldNodeConfig(clusterDefaults, poolConfig, nodeConfigs)

	// Alphabetical order means "b-override" wins over "a-override"
	if got.Kubelet.MaxPods != 500 {
		t.Errorf("Kubelet.MaxPods = %d, want 500 (b-override)", got.Kubelet.MaxPods)
	}

	wantFiles := []nanokubev1alpha1.FileSourceSpec{
		{Path: "/etc/kubernetes/audit.yaml", Mode: "0644", Source: nanokubev1alpha1.FileSourceRefSpec{Inline: "pool"}},
		{Path: "/etc/kubernetes/pool.conf", Mode: "0600", Source: nanokubev1alpha1.FileSourceRefSpec{Inline: "pool-file"}},
	}

	if !reflect.DeepEqual(got.Files, wantFiles) {
		t.Errorf("Files = %v, want %v", got.Files, wantFiles)
	}
}
