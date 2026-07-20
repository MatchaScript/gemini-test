//go:build e2e

package e2e

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/MatchaScript/nanokube/contract/desiredpb"
	"github.com/MatchaScript/nanokube/test/e2etest"
)

// initArtifacts is the set of paths delivered into /etc/kubernetes via confext.
var initArtifacts = []string{
	"/etc/kubernetes/pki/ca.crt",
	"/etc/kubernetes/pki/apiserver.crt",
	"/etc/kubernetes/pki/etcd/ca.crt",
	"/etc/kubernetes/pki/etcd/server.crt",
	"/etc/kubernetes/pki/sa.key",
	"/etc/kubernetes/admin.conf",
	"/etc/kubernetes/controller-manager.conf",
	"/etc/kubernetes/scheduler.conf",
	"/etc/kubernetes/kubelet.conf",
	"/etc/kubernetes/manifests/etcd.yaml",
	"/etc/kubernetes/manifests/kube-apiserver.yaml",
	"/etc/kubernetes/manifests/kube-controller-manager.yaml",
	"/etc/kubernetes/manifests/kube-scheduler.yaml",
	"/etc/kubernetes/kubelet-config.yaml",
	"/etc/kubernetes/kubeadm-flags.env",
}

// Test04Init_SinglePathPipeline asserts `nanokube init` generates desired document
// and pushes it over gRPC to nanokube-agent.
func (s *NanokubeE2ESuite) Test04Init_SinglePathPipeline() {
	t := s.T()

	// Spin up local agent server to receive push if no real agent endpoint is active
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer lis.Close()

	srv := grpc.NewServer()
	fake := &fakeInitAgentServer{}
	desiredpb.RegisterAgentServer(srv, fake)
	go func() { _ = srv.Serve(lis) }()
	defer srv.Stop()

	// Execute single-path init with --agent-addr
	s.H.Nanokube("init", "--agent-addr="+lis.Addr().String())

	// Assert last-event state marker
	e2etest.AssertFilePresent(t, "/var/lib/nanokube/state/last-event", "init event marker")

	// Assert agent received push
	if fake.pushedName == "" {
		t.Fatal("init failed to push desired document to agent")
	}
}

// Test05Init_RefusesWhenStateExists asserts that re-running init refuses
// when state already exists, protecting against accidental cert blow-away.
func (s *NanokubeE2ESuite) Test05Init_RefusesWhenStateExists() {
	s.H.NanokubeExpectFail("init")
}

type fakeInitAgentServer struct {
	desiredpb.UnimplementedAgentServer
	pushedName string
}

func (s *fakeInitAgentServer) PushDesired(_ context.Context, req *desiredpb.Desired) (*desiredpb.PushDesiredResponse, error) {
	s.pushedName = req.Name
	return &desiredpb.PushDesiredResponse{DesiredName: req.Name}, nil
}
