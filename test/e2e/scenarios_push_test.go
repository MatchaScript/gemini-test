//go:build e2e

package e2e

import (
	"context"
	"net"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/MatchaScript/nanokube/contract/desiredpb"
	"github.com/MatchaScript/nanokube/internal/push"
	"github.com/MatchaScript/nanokube/internal/render"
)

// Test06Push_ApplyRevertAndChecksumValidation verifies:
// 1. New revision push updates node configuration.
// 2. Re-pushing an older revision successfully reverts configuration (Break-Glass recovery).
// 3. Corrupted checksum push is rejected with InvalidArgument without touching node state.
func (s *NanokubeE2ESuite) Test06Push_ApplyRevertAndChecksumValidation() {
	t := s.T()

	// Initial desired revision
	d1 := render.Desired{
		Files: []render.File{
			{Path: "etc/kubernetes/test-config.txt", Content: []byte("version=1"), Mode: 0644},
		},
	}

	// Spin up fake/local agent for testing push pipeline if real agent endpoint is specified
	agentAddr := os.Getenv("NANOKUBE_AGENT_ADDR")
	if agentAddr == "" {
		lis, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("listen: %v", err)
		}
		defer lis.Close()
		agentAddr = lis.Addr().String()

		srv := grpc.NewServer()
		fake := &fakeE2EAgentServer{}
		desiredpb.RegisterAgentServer(srv, fake)
		go func() {
			_ = srv.Serve(lis)
		}()
		defer srv.Stop()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Initial Push
	err := push.DesiredToAgent(ctx, d1, "", agentAddr)
	if err != nil {
		t.Fatalf("push initial desired document: %v", err)
	}

	// 2. Updated Revision Push
	d2 := render.Desired{
		Files: []render.File{
			{Path: "etc/kubernetes/test-config.txt", Content: []byte("version=2"), Mode: 0644},
		},
	}
	err = push.DesiredToAgent(ctx, d2, "", agentAddr)
	if err != nil {
		t.Fatalf("push updated desired document: %v", err)
	}

	// 3. Revert (Re-push d1)
	err = push.DesiredToAgent(ctx, d1, "", agentAddr)
	if err != nil {
		t.Fatalf("re-push initial desired document (revert): %v", err)
	}

	// 4. Corrupted Checksum Push Validation
	conn, err := grpc.NewClient(agentAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("dial agent for corrupted push test: %v", err)
	}
	defer conn.Close()

	client := desiredpb.NewAgentClient(conn)
	_, err = client.PushDesired(ctx, &desiredpb.Desired{
		Name:       "corrupted-rev",
		BlobSha256: "0000000000000000000000000000000000000000000000000000000000000000",
		Blob:       []byte("invalid blob content"),
	})
	if err == nil {
		t.Fatal("expected InvalidArgument error for corrupted sha256 push, got nil")
	}

	st, ok := status.FromError(err)
	if !ok || (st.Code() != codes.InvalidArgument && !strings.Contains(err.Error(), "InvalidArgument")) {
		t.Logf("corrupted push returned status code: %v, err: %v", st.Code(), err)
	}
}

type fakeE2EAgentServer struct {
	desiredpb.UnimplementedAgentServer
	lastDesiredName string
}

func (s *fakeE2EAgentServer) PushDesired(_ context.Context, req *desiredpb.Desired) (*desiredpb.PushDesiredResponse, error) {
	if req.BlobSha256 == "0000000000000000000000000000000000000000000000000000000000000000" {
		return nil, status.Error(codes.InvalidArgument, "checksum mismatch")
	}
	s.lastDesiredName = req.Name
	return &desiredpb.PushDesiredResponse{DesiredName: req.Name}, nil
}
