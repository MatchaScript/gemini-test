package e2etest

import (
	"os"
	"testing"
)

// Helpers holds shared state and provides cli / kubectl / systemctl
// wrappers for the e2e suite. Create one per suite (or per test) with
// New; helper methods that fail the test do so against the t passed in.
type Helpers struct {
	t          testing.TB
	bin        string
	kubeconfig string
	nodeName   string
	flannelURL string
}

// Config configures a Helpers instance.
type Config struct {
	Bin        string // /usr/bin/nanokube
	Kubeconfig string // /etc/kubernetes/admin.conf
	NodeName   string // lowercased hostname
	FlannelURL string // overridable; default = github releases /latest
}

// New returns a Helpers bound to t.
func New(t testing.TB, cfg Config) *Helpers {
	return &Helpers{
		t:          t,
		bin:        cfg.Bin,
		kubeconfig: cfg.Kubeconfig,
		nodeName:   cfg.NodeName,
		flannelURL: cfg.FlannelURL,
	}
}

// NodeName returns the resolved node name used by waits/assertions.
func (h *Helpers) NodeName() string { return h.nodeName }

// IsK8sAvailable checks if admin.conf exists on disk indicating an active or initialised cluster.
func IsK8sAvailable() bool {
	_, err := os.Stat("/etc/kubernetes/admin.conf")
	return err == nil
}
