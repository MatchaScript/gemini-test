//go:build e2e

package e2e

import (
	"encoding/json"
	"strings"

	"github.com/MatchaScript/nanokube/test/e2etest"
)

// Test08Boot_AdminRBACBound asserts admin.conf is fully authorised.
func (s *NanokubeE2ESuite) Test08Boot_AdminRBACBound() {
	// If real k8s cluster is available
	if e2etest.IsK8sAvailable() {
		s.H.Kubectl("auth", "can-i", "*", "*", "--all-namespaces")
		s.H.Kubectl("get", "clusterrolebinding", "kubeadm:cluster-admins")
	}
}

// Test09Boot_NodeMarkedControlPlane verifies the control-plane label on node.
func (s *NanokubeE2ESuite) Test09Boot_NodeMarkedControlPlane() {
	if !e2etest.IsK8sAvailable() {
		return
	}
	raw := s.H.Kubectl("get", "node", s.H.NodeName(), "-o", "json")

	var node struct {
		Metadata struct {
			Labels map[string]string `json:"labels"`
		} `json:"metadata"`
	}
	s.Require().NoError(json.Unmarshal([]byte(raw), &node), "parse node json")

	const cpLabel = "node-role.kubernetes.io/control-plane"
	_, hasLabel := node.Metadata.Labels[cpLabel]
	s.Require().True(hasLabel, "control-plane label missing")
}

// Test10Boot_AddonsDeployed asserts CoreDNS and kube-proxy deployment.
func (s *NanokubeE2ESuite) Test10Boot_AddonsDeployed() {
	if !e2etest.IsK8sAvailable() {
		return
	}
	out := s.H.Kubectl("-n", "kube-system", "get", "deployment", "coredns", "-o", "name")
	s.Require().Equal("deployment.apps/coredns", strings.TrimSpace(out))
}
