package operator

import (
	"context"
	"testing"

	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func TestEvaluateInPlaceUpdate_VersionUpgrade(t *testing.T) {
	m := &clusterv1.Machine{
		Spec: clusterv1.MachineSpec{
			Version: ptr.To("v1.33.1"),
		},
	}

	req := InPlaceUpdateRequest{
		Machine:       m,
		TargetVersion: "v1.33.2",
	}

	resp := EvaluateInPlaceUpdate(context.Background(), req)
	if !resp.CanInPlaceUpdate {
		t.Errorf("CanInPlaceUpdate = false, want true for version upgrade")
	}
}

func TestEvaluateInPlaceUpdate_SameVersionConfigDiff(t *testing.T) {
	m := &clusterv1.Machine{
		Spec: clusterv1.MachineSpec{
			Version: ptr.To("v1.33.2"),
		},
	}

	req := InPlaceUpdateRequest{
		Machine:        m,
		TargetVersion:  "v1.33.2",
		TargetRevision: "rev-new-hash",
	}

	resp := EvaluateInPlaceUpdate(context.Background(), req)
	if !resp.CanInPlaceUpdate {
		t.Errorf("CanInPlaceUpdate = false, want true for config diff")
	}
}
