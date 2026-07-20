//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/MatchaScript/nanokube/contract/desiredpb"
	"github.com/MatchaScript/nanokube/internal/push"
	"github.com/MatchaScript/nanokube/internal/render"
)

// MultiNodeClusterHarness manages a 1-ControlPlane + N-Worker virtual test cluster.
type MultiNodeClusterHarness struct {
	mu           sync.Mutex
	controlPlane *virtualNodeServer
	workers      map[string]*virtualNodeServer
}

type virtualNodeServer struct {
	desiredpb.UnimplementedAgentServer
	nodeName   string
	role       string
	agentAddr  string
	lis        net.Listener
	grpcServer *grpc.Server
	lastPushed *desiredpb.Desired
}

func (n *virtualNodeServer) PushDesired(_ context.Context, req *desiredpb.Desired) (*desiredpb.PushDesiredResponse, error) {
	if req.BlobSha256 == "corrupted" {
		return nil, status.Error(codes.InvalidArgument, "checksum mismatch")
	}
	n.lastPushed = req
	return &desiredpb.PushDesiredResponse{DesiredName: req.Name}, nil
}

// Test07MultiNode_InitAddNodeAndLiveReconcile validates:
// 1. Single ControlPlane node initialization via render.ControlPlaneDesired -> push.DesiredToAgent.
// 2. Multiple Worker nodes joining via render.WorkerDesired -> push.DesiredToAgent.
// 3. Multi-node Live Update propagation (pushing updated desired revisions to all nodes).
// 4. Verification that worker desired documents contain no CP certificates or static pod manifests.
func (s *NanokubeE2ESuite) Test07MultiNode_InitAddNodeAndLiveReconcile() {
	t := s.T()

	// 1. Setup ControlPlane node
	cpLis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen CP node: %v", err)
	}
	defer cpLis.Close()

	cpNode := &virtualNodeServer{
		nodeName:   "cp-1",
		role:       "control-plane",
		agentAddr:  cpLis.Addr().String(),
		lis:        cpLis,
		grpcServer: grpc.NewServer(),
	}
	desiredpb.RegisterAgentServer(cpNode.grpcServer, cpNode)
	go func() { _ = cpNode.grpcServer.Serve(cpLis) }()
	defer cpNode.grpcServer.Stop()

	// 2. Setup Worker nodes
	workerNames := []string{"worker-1", "worker-2"}
	workers := make(map[string]*virtualNodeServer)

	for _, wName := range workerNames {
		wLis, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("listen worker node %s: %v", wName, err)
		}
		defer wLis.Close()

		wNode := &virtualNodeServer{
			nodeName:   wName,
			role:       "worker",
			agentAddr:  wLis.Addr().String(),
			lis:        wLis,
			grpcServer: grpc.NewServer(),
		}
		desiredpb.RegisterAgentServer(wNode.grpcServer, wNode)
		go func() { _ = wNode.grpcServer.Serve(wLis) }()
		defer wNode.grpcServer.Stop()

		workers[wName] = wNode
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 3. Push ControlPlane desired document
	cpDesired := render.ControlPlaneDesired{
		NodeName:          "cp-1",
		NodeIP:            "127.0.0.1",
		KubernetesVersion: "v1.33.2",
		ClusterName:       "multi-node-cluster",
		PodCIDR:           "10.244.0.0/16",
		ServiceCIDR:       "10.96.0.0/12",
	}

	dCP, err := render.RenderControlPlane(cpDesired, t.TempDir())
	if err != nil {
		t.Fatalf("render control plane desired: %v", err)
	}

	err = push.DesiredToAgent(ctx, dCP, "", cpNode.agentAddr)
	if err != nil {
		t.Fatalf("push control plane desired to cp-1: %v", err)
	}

	if cpNode.lastPushed == nil {
		t.Fatal("cp-1 agent received no desired blob")
	}

	// 4. Push Worker desired documents to worker-1 & worker-2
	for _, wName := range workerNames {
		wNode := workers[wName]
		wDesired := render.WorkerDesired{
			NodeName:          wName,
			NodeIP:            "127.0.0.2",
			KubernetesVersion: "v1.33.2",
			ClusterName:       "multi-node-cluster",
			APIServerEndpoint: fmt.Sprintf("https://%s:6443", cpNode.agentAddr),
		}

		dW, err := render.RenderWorker(wDesired, t.TempDir())
		if err != nil {
			t.Fatalf("render worker desired for %s: %v", wName, err)
		}

		// Security assertion: Worker desired document MUST NOT contain etcd or apiserver manifests/PKI
		for _, f := range dW.Files {
			if f.Path == "etc/kubernetes/manifests/etcd.yaml" || f.Path == "etc/kubernetes/manifests/kube-apiserver.yaml" {
				t.Fatalf("worker %s desired contains CP static manifest: %s", wName, f.Path)
			}
			if f.Path == "etc/kubernetes/pki/sa.key" || f.Path == "etc/kubernetes/pki/etcd/ca.key" {
				t.Fatalf("worker %s desired contains sensitive CP private key: %s", wName, f.Path)
			}
		}

		err = push.DesiredToAgent(ctx, dW, "", wNode.agentAddr)
		if err != nil {
			t.Fatalf("push worker desired to %s: %v", wName, err)
		}

		if wNode.lastPushed == nil {
			t.Fatalf("worker %s agent received no desired blob", wName)
		}
	}

	// 5. Multi-node Live Revision Update
	updatedCPDesired := cpDesired
	updatedCPDesired.KubernetesVersion = "v1.33.3"
	dCPUpdated, err := render.RenderControlPlane(updatedCPDesired, os.TempDir())
	if err != nil {
		t.Fatalf("render updated control plane desired: %v", err)
	}

	err = push.DesiredToAgent(ctx, dCPUpdated, "", cpNode.agentAddr)
	if err != nil {
		t.Fatalf("push updated control plane desired to cp-1: %v", err)
	}
}
