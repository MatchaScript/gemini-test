//go:build e2e

package e2e

import (
	"os"
	"regexp"
	"strings"

	"github.com/MatchaScript/nanokube/test/e2etest"
)

// Test06JoinPrereqs_ClusterObjectsExist asserts that EnsureJoinPrereqs
// seeded all kubeadm join-path objects when cluster is running.
func (s *NanokubeE2ESuite) Test06JoinPrereqs_ClusterObjectsExist() {
	if !e2etest.IsK8sAvailable() {
		return
	}
	s.H.Kubectl("get", "configmap", "kubeadm-config", "-n", "kube-system")
	s.H.Kubectl("get", "configmap", "kubelet-config", "-n", "kube-system")
	s.H.Kubectl("get", "configmap", "cluster-info", "-n", "kube-public")
}

// Test06JoinPrereqs_LastBootRecordsRole asserts state contains last event.
func (s *NanokubeE2ESuite) Test06JoinPrereqs_LastBootRecordsRole() {
	if !e2etest.IsK8sAvailable() {
		return
	}
	b, err := os.ReadFile("/var/lib/nanokube/state/last-event")
	s.Require().NoError(err)
	s.Contains(string(b), `"init"`)
}

// Test06JoinPrereqs_TokenCreate mints a join token via `nanokube token create`.
func (s *NanokubeE2ESuite) Test06JoinPrereqs_TokenCreate() {
	if !e2etest.IsK8sAvailable() {
		return
	}
	out, _ := s.H.Nanokube("token", "create")
	s.Contains(out, "token: ")
	s.Contains(out, "ca-cert-hash: sha256:")
	tokenLine := regexp.MustCompile(`token: (\S+)`).FindStringSubmatch(out)
	s.Require().Len(tokenLine, 2)
	id := strings.SplitN(tokenLine[1], ".", 2)[0]
	s.H.Kubectl("get", "secret", "-n", "kube-system", "bootstrap-token-"+id)
}
