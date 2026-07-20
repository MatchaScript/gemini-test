// Package addnode implements `nanokube add-node`: worker node join flow.
// It discovers cluster configuration via bootstrap token, renders the worker
// desired document, and pushes it to nanokube-agent over gRPC.
package addnode

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"strings"

	"k8s.io/client-go/tools/clientcmd"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	kubeadmapiv1 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta4"
	"k8s.io/kubernetes/cmd/kubeadm/app/discovery"
	kubeadmconfig "k8s.io/kubernetes/cmd/kubeadm/app/util/config"

	"github.com/MatchaScript/nanokube/internal/layout"
	"github.com/MatchaScript/nanokube/internal/push"
	"github.com/MatchaScript/nanokube/internal/render"
)

// Options are the operator-supplied join credentials.
type Options struct {
	Server       string   // reachable apiserver, host:port or https://host:port
	Token        string   // bootstrap token (id.secret)
	CACertHashes []string // sha256:... pins; required
	AgentAddr    string   // gRPC endpoint of nanokube-agent
	FileContexts string   // SELinux file_contexts database
}

func normalizeServer(s string) (fullURL, hostPort, host string, err error) {
	if !strings.Contains(s, "://") {
		s = "https://" + s
	}
	u, err := url.Parse(s)
	if err != nil || u.Host == "" {
		return "", "", "", fmt.Errorf("invalid --server %q: want https://host:port", s)
	}
	h, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		return "", "", "", fmt.Errorf("invalid --server %q: missing port", s)
	}
	return "https://" + u.Host, u.Host, h, nil
}

func Run(ctx context.Context, opts Options, l layout.Layout, out io.Writer) error {
	logf := func(format string, a ...any) { fmt.Fprintf(out, "[add-node] "+format+"\n", a...) }

	// Idempotency: a node that completed TLS bootstrap is joined.
	if _, err := os.Stat(l.KubeletKubeconfig); err == nil {
		logf("kubelet.conf already present — node already joined; nothing to do")
		return nil
	}

	if len(opts.CACertHashes) == 0 {
		return fmt.Errorf("--ca-cert-hash is required (printed by `nanokube token create`)")
	}

	_, hostPort, _, err := normalizeServer(opts.Server)
	if err != nil {
		return err
	}

	versionedJoin := &kubeadmapiv1.JoinConfiguration{
		Discovery: kubeadmapiv1.Discovery{
			BootstrapToken: &kubeadmapiv1.BootstrapTokenDiscovery{
				Token:             opts.Token,
				APIServerEndpoint: hostPort,
				CACertHashes:      opts.CACertHashes,
			},
			TLSBootstrapToken: opts.Token,
		},
	}
	joinCfg, err := kubeadmconfig.DefaultedJoinConfiguration(versionedJoin, kubeadmconfig.LoadOrDefaultConfigurationOptions{})
	if err != nil {
		return fmt.Errorf("default join configuration: %w", err)
	}

	logf("discovering cluster via %s", hostPort)
	tlsBootstrapCfg, err := discovery.For(nil, joinCfg)
	if err != nil {
		return fmt.Errorf("discovery: %w", err)
	}

	bootstrapBytes, err := clientcmd.Write(*tlsBootstrapCfg)
	if err != nil {
		return fmt.Errorf("serialize bootstrap kubeconfig: %w", err)
	}

	initCfg := &kubeadmapi.InitConfiguration{}
	initCfg.NodeRegistration.Name = joinCfg.NodeRegistration.Name

	d, err := render.WorkerDesired(initCfg, bootstrapBytes)
	if err != nil {
		return fmt.Errorf("render worker desired: %w", err)
	}

	agentAddr := opts.AgentAddr
	if agentAddr == "" {
		agentAddr = "127.0.0.1:50051"
	}

	logf("pushing worker desired document (%s) to agent at %s", d.Name(), agentAddr)
	if err := push.DesiredToAgent(ctx, d, opts.FileContexts, agentAddr); err != nil {
		return fmt.Errorf("push to agent: %w", err)
	}

	logf("add-node push complete (revision=%s)", d.Name())
	return nil
}
