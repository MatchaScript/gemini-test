package push

import (
	"context"
	"errors"
	"net"
	"strings"
	"testing"

	"google.golang.org/grpc"

	"github.com/MatchaScript/nanokube/contract/desiredpb"
	"github.com/MatchaScript/nanokube/internal/ddi"
	"github.com/MatchaScript/nanokube/internal/render"
)

type fakeAgentServer struct {
	desiredpb.UnimplementedAgentServer
	pushedName string
	pushedBlob []byte
}

func (s *fakeAgentServer) PushDesired(_ context.Context, req *desiredpb.Desired) (*desiredpb.PushDesiredResponse, error) {
	s.pushedName = req.Name
	s.pushedBlob = req.Blob
	return &desiredpb.PushDesiredResponse{DesiredName: req.Name}, nil
}

func skipIfRepartUnusable(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		return
	}
	if errors.Is(err, ddi.ErrSystemdRepartNotFound) {
		t.Skipf("systemd-repart not found in PATH: skipping real DDI build")
	}
	if strings.Contains(err.Error(), "mkfs.erofs") {
		t.Skipf("systemd-repart present but mkfs.erofs unavailable: skipping real DDI build: %v", err)
	}
	t.Fatalf("Build: %v", err)
}

func TestDesiredToAgent_PushesSuccessfully(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer lis.Close()

	srv := grpc.NewServer()
	fake := &fakeAgentServer{}
	desiredpb.RegisterAgentServer(srv, fake)
	go func() {
		_ = srv.Serve(lis)
	}()
	defer srv.Stop()

	d := render.Desired{
		Files: []render.File{
			{Path: "etc/test.txt", Content: []byte("hello"), Mode: 0644},
		},
	}

	ctx := context.Background()
	err = DesiredToAgent(ctx, d, "", lis.Addr().String())
	skipIfRepartUnusable(t, err)

	if fake.pushedName != d.Name() {
		t.Errorf("pushedName = %q, want %q", fake.pushedName, d.Name())
	}
}
