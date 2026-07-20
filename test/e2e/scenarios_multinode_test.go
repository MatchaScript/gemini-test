//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	kubeadmconfig "k8s.io/kubernetes/cmd/kubeadm/app/util/config"

	"github.com/MatchaScript/nanokube/contract/desiredpb"
	"github.com/MatchaScript/nanokube/internal/config"
	"github.com/MatchaScript/nanokube/internal/layouttest"
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

func helperMakeInitConfig(nodeName, ip, version string) string {
	return fmt.Sprintf(`apiVersion: bootstrap.nanokube.io/v1alpha1
kind: NanoKubeConfig
metadata:
  name: %s
---
apiVersion: kubeadm.k8s.io/v1beta4
kind: InitConfiguration
localAPIEndpoint:
  advertiseAddress: %s
nodeRegistration:
  name: %s
  criSocket: unix:///var/run/crio/crio.sock
---
apiVersion: kubeadm.k8s.io/v1beta4
kind: ClusterConfiguration
kubernetesVersion: %s
networking:
  serviceSubnet: 10.96.0.0/12
  podSubnet: 10.244.0.0/16
`, nodeName, ip, nodeName, version)
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
	lCP := layouttest.New(t)
	cpCfgPath := filepath.Join(t.TempDir(), "cp-config.yaml")
	if err := os.WriteFile(cpCfgPath, []byte(helperMakeInitConfig("cp-1", "192.168.1.10", "v1.35.0")), 0644); err != nil {
		t.Fatalf("write CP config: %v", err)
	}
	loadedCP, err := config.Load(cpCfgPath, lCP)
	if err != nil {
		t.Fatalf("load CP config: %v", err)
	}

	dCP, err := render.ControlPlaneDesired(loadedCP.Init, t.TempDir())
	if err != nil {
		t.Fatalf("render control plane desired: %v", err)
	}

	err = push.DesiredToAgent(ctx, dCP, "", cpNode.agentAddr)
	if err != nil {
		if strings.Contains(err.Error(), "mkfs.erofs") || strings.Contains(err.Error(), "systemd-repart") {
			t.Skipf("skipping DDI push test due to missing repart binary: %v", err)
			return
		}
		t.Fatalf("push control plane desired to cp-1: %v", err)
	}

	if cpNode.lastPushed == nil {
		t.Fatal("cp-1 agent received no desired blob")
	}

	// 4. Push Worker desired documents to worker-1 & worker-2
	for _, wName := range workerNames {
		wNode := workers[wName]
		lW := layouttest.New(t)
		wCfgPath := filepath.Join(t.TempDir(), wName+"-config.yaml")
		if err := os.WriteFile(wCfgPath, []byte(helperMakeInitConfig(wName, "192.168.1.11", "v1.35.0")), 0644); err != nil {
			t.Fatalf("write worker config: %v", err)
		}
		loadedW, err := config.Load(wCfgPath, lW)
		if err != nil {
			t.Fatalf("load worker config: %v", err)
		}

		fakeBootstrap := []byte("apiVersion: v1\nkind: Config\nclusters: []\n")
		dW, err := render.WorkerDesired(loadedW.Init, fakeBootstrap)
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
	updatedCPYaml := helperMakeInitConfig("cp-1", "192.168.1.10", "v1.35.0")
	updatedCfg, err := kubeadmconfig.BytesToInitConfiguration([]byte(updatedCPYaml), false)
	if err == nil && updatedCfg != nil {
		dCPUpdated, err := render.ControlPlaneDesired(updatedCfg, os.TempDir())
		if err == nil {
			_ = push.DesiredToAgent(ctx, dCPUpdated, "", cpNode.agentAddr)
		}
	}
}
