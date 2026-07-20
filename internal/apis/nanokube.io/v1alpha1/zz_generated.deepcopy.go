package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyObject implementations for runtime.Object interface

func (in *Cluster) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *Cluster) DeepCopy() *Cluster {
	if in == nil {
		return nil
	}
	out := new(Cluster)
	*out = *in
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	return out
}

func (in *ClusterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *ClusterList) DeepCopy() *ClusterList {
	if in == nil {
		return nil
	}
	out := new(ClusterList)
	*out = *in
	if in.Items != nil {
		out.Items = make([]Cluster, len(in.Items))
		for i := range in.Items {
			out.Items[i] = *in.Items[i].DeepCopy()
		}
	}
	return out
}

func (in *NodePool) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *NodePool) DeepCopy() *NodePool {
	if in == nil {
		return nil
	}
	out := new(NodePool)
	*out = *in
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	return out
}

func (in *NodePoolList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *NodePoolList) DeepCopy() *NodePoolList {
	if in == nil {
		return nil
	}
	out := new(NodePoolList)
	*out = *in
	if in.Items != nil {
		out.Items = make([]NodePool, len(in.Items))
		for i := range in.Items {
			out.Items[i] = *in.Items[i].DeepCopy()
		}
	}
	return out
}

func (in *NodeConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *NodeConfig) DeepCopy() *NodeConfig {
	if in == nil {
		return nil
	}
	out := new(NodeConfig)
	*out = *in
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	return out
}

func (in *NodeConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *NodeConfigList) DeepCopy() *NodeConfigList {
	if in == nil {
		return nil
	}
	out := new(NodeConfigList)
	*out = *in
	if in.Items != nil {
		out.Items = make([]NodeConfig, len(in.Items))
		for i := range in.Items {
			out.Items[i] = *in.Items[i].DeepCopy()
		}
	}
	return out
}

func (in *InventoryMachine) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *InventoryMachine) DeepCopy() *InventoryMachine {
	if in == nil {
		return nil
	}
	out := new(InventoryMachine)
	*out = *in
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	return out
}

func (in *InventoryMachineList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *InventoryMachineList) DeepCopy() *InventoryMachineList {
	if in == nil {
		return nil
	}
	out := new(InventoryMachineList)
	*out = *in
	if in.Items != nil {
		out.Items = make([]InventoryMachine, len(in.Items))
		for i := range in.Items {
			out.Items[i] = *in.Items[i].DeepCopy()
		}
	}
	return out
}

func (in *NanokubeCluster) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *NanokubeCluster) DeepCopy() *NanokubeCluster {
	if in == nil {
		return nil
	}
	out := new(NanokubeCluster)
	*out = *in
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	return out
}

func (in *NanokubeClusterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *NanokubeClusterList) DeepCopy() *NanokubeClusterList {
	if in == nil {
		return nil
	}
	out := new(NanokubeClusterList)
	*out = *in
	if in.Items != nil {
		out.Items = make([]NanokubeCluster, len(in.Items))
		for i := range in.Items {
			out.Items[i] = *in.Items[i].DeepCopy()
		}
	}
	return out
}

func (in *NanokubeControlPlane) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *NanokubeControlPlane) DeepCopy() *NanokubeControlPlane {
	if in == nil {
		return nil
	}
	out := new(NanokubeControlPlane)
	*out = *in
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	return out
}

func (in *NanokubeControlPlaneList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *NanokubeControlPlaneList) DeepCopy() *NanokubeControlPlaneList {
	if in == nil {
		return nil
	}
	out := new(NanokubeControlPlaneList)
	*out = *in
	if in.Items != nil {
		out.Items = make([]NanokubeControlPlane, len(in.Items))
		for i := range in.Items {
			out.Items[i] = *in.Items[i].DeepCopy()
		}
	}
	return out
}

func (in *NanokubeMachine) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *NanokubeMachine) DeepCopy() *NanokubeMachine {
	if in == nil {
		return nil
	}
	out := new(NanokubeMachine)
	*out = *in
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	return out
}

func (in *NanokubeMachineList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *NanokubeMachineList) DeepCopy() *NanokubeMachineList {
	if in == nil {
		return nil
	}
	out := new(NanokubeMachineList)
	*out = *in
	if in.Items != nil {
		out.Items = make([]NanokubeMachine, len(in.Items))
		for i := range in.Items {
			out.Items[i] = *in.Items[i].DeepCopy()
		}
	}
	return out
}

func (in *NanokubeMachineTemplate) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *NanokubeMachineTemplate) DeepCopy() *NanokubeMachineTemplate {
	if in == nil {
		return nil
	}
	out := new(NanokubeMachineTemplate)
	*out = *in
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	return out
}

func (in *NanokubeMachineTemplateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *NanokubeMachineTemplateList) DeepCopy() *NanokubeMachineTemplateList {
	if in == nil {
		return nil
	}
	out := new(NanokubeMachineTemplateList)
	*out = *in
	if in.Items != nil {
		out.Items = make([]NanokubeMachineTemplate, len(in.Items))
		for i := range in.Items {
			out.Items[i] = *in.Items[i].DeepCopy()
		}
	}
	return out
}
