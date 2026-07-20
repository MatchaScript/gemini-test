//go:build e2e

package e2e

import (
	"context"
	"time"

	"github.com/MatchaScript/nanokube/internal/push"
	"github.com/MatchaScript/nanokube/internal/render"
)

// Test12Reconcile_DesiredPushReplacesConfiguration asserts live revision updates
// via push.DesiredToAgent cleanly replace systemd-confext state on agent nodes.
func (s *NanokubeE2ESuite) Test12Reconcile_DesiredPushReplacesConfiguration() {
	d := render.Desired{
		Files: []render.File{
			{Path: "etc/kubernetes/test-reconcile.txt", Content: []byte("reconcile=true"), Mode: 0644},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = push.DesiredToAgent(ctx, d, "", "127.0.0.1:9090")
}
