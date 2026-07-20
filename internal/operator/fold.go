package operator

import (
	"sort"

	nanokubev1alpha1 "github.com/MatchaScript/nanokube/internal/apis/nanokube.io/v1alpha1"
)

// FoldNodeConfig merges configuration layers in precedence order:
// Cluster.nodeConfigDefaults -> NodePool.config -> NodeConfig.config
// Specific overrides general. If multiple NodeConfigs match, they are merged in alphabetical order by Name.
func FoldNodeConfig(clusterDefaults nanokubev1alpha1.NodeConfigSpec, poolConfig nanokubev1alpha1.NodeConfigSpec, nodeConfigs []nanokubev1alpha1.NodeConfig) nanokubev1alpha1.NodeConfigSpec {
	res := nanokubev1alpha1.NodeConfigSpec{
		Kubelet: clusterDefaults.Kubelet,
		Files:   append([]nanokubev1alpha1.FileSourceSpec(nil), clusterDefaults.Files...),
	}

	// 1. Merge NodePool config
	if poolConfig.Kubelet.MaxPods > 0 {
		res.Kubelet.MaxPods = poolConfig.Kubelet.MaxPods
	}
	res.Files = mergeFiles(res.Files, poolConfig.Files)

	// 2. Sort NodeConfigs by name for deterministic order (CiliumNodeConfig precedence rule)
	sortedNC := append([]nanokubev1alpha1.NodeConfig(nil), nodeConfigs...)
	sort.Slice(sortedNC, func(i, j int) bool {
		return sortedNC[i].Name < sortedNC[j].Name
	})

	// 3. Merge NodeConfigs
	for _, nc := range sortedNC {
		if nc.Spec.Config.Kubelet.MaxPods > 0 {
			res.Kubelet.MaxPods = nc.Spec.Config.Kubelet.MaxPods
		}
		res.Files = mergeFiles(res.Files, nc.Spec.Config.Files)
	}

	return res
}

func mergeFiles(base, override []nanokubev1alpha1.FileSourceSpec) []nanokubev1alpha1.FileSourceSpec {
	fileMap := make(map[string]nanokubev1alpha1.FileSourceSpec)
	var paths []string

	for _, f := range base {
		if _, exists := fileMap[f.Path]; !exists {
			paths = append(paths, f.Path)
		}
		fileMap[f.Path] = f
	}

	for _, f := range override {
		if _, exists := fileMap[f.Path]; !exists {
			paths = append(paths, f.Path)
		}
		fileMap[f.Path] = f
	}

	result := make([]nanokubev1alpha1.FileSourceSpec, 0, len(paths))
	for _, p := range paths {
		result = append(result, fileMap[p])
	}
	return result
}
