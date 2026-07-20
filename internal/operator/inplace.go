package operator

import (
	"context"
	"fmt"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// InPlaceUpdateRequest describes a request to evaluate if a machine can be updated in-place.
type InPlaceUpdateRequest struct {
	Machine        *clusterv1.Machine
	TargetVersion  string
	TargetRevision string
}

// InPlaceUpdateResponse represents the evaluation result.
type InPlaceUpdateResponse struct {
	CanInPlaceUpdate bool
	Reason           string
}

// EvaluateInPlaceUpdate determines whether a machine version/config diff can be applied in-place
// without replacing the machine instance.
func EvaluateInPlaceUpdate(ctx context.Context, req InPlaceUpdateRequest) InPlaceUpdateResponse {
	if req.Machine == nil {
		return InPlaceUpdateResponse{CanInPlaceUpdate: false, Reason: "machine is nil"}
	}

	currentVersion := req.Machine.Spec.Version
	if currentVersion != nil && *currentVersion != "" && req.TargetVersion != "" {
		if *currentVersion != req.TargetVersion {
			// Version diff (e.g. v1.33.1 -> v1.33.2) is supported in-place by nanokube agent
			return InPlaceUpdateResponse{
				CanInPlaceUpdate: true,
				Reason:           fmt.Sprintf("in-place version upgrade supported (%s -> %s)", *currentVersion, req.TargetVersion),
			}
		}
	}

	// Config diff (revision change) is supported live via confext refresh
	return InPlaceUpdateResponse{
		CanInPlaceUpdate: true,
		Reason:           "config update supported live via confext refresh",
	}
}
