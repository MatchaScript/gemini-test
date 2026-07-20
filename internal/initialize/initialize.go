// Package initialize implements `nanokube init`: the discrete,
// node-local control-plane bootstrap that renders the control-plane desired
// document and pushes it to nanokube-agent over gRPC.
package initialize

import (
	"context"
	"fmt"
	"io"

	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"

	"github.com/MatchaScript/nanokube/internal/layout"
	"github.com/MatchaScript/nanokube/internal/push"
	"github.com/MatchaScript/nanokube/internal/render"
)

// Run executes the one-time init: renders the control-plane desired document
// (kubelet config + flags env, static pod manifests, PKI + kubeconfigs)
// and pushes it directly to nanokube-agent over gRPC.
func Run(ctx context.Context, cfg *kubeadmapi.InitConfiguration, l layout.Layout, fileContexts, agentAddr string, out io.Writer) error {
	logf := func(format string, a ...any) { fmt.Fprintf(out, "[init] "+format+"\n", a...) }
	nodeName := cfg.NodeRegistration.Name

	logf("rendering control-plane desired document for node %s", nodeName)
	d, err := render.ControlPlaneDesired(cfg, l.NanoKubeVarDir)
	if err != nil {
		return fmt.Errorf("render control-plane desired: %w", err)
	}

	logf("pushing desired document (%s) to agent at %s", d.Name(), agentAddr)
	if err := push.DesiredToAgent(ctx, d, fileContexts, agentAddr); err != nil {
		return fmt.Errorf("push to agent: %w", err)
	}

	logf("init push complete (revision=%s)", d.Name())
	return nil
}
