//go:build e2e

package e2e

import (
	"os"

	"github.com/MatchaScript/nanokube/test/e2etest"
)

// Test14Reset_StateCleanup asserts state markers can be removed.
func (s *NanokubeE2ESuite) Test14Reset_StateCleanup() {
	_ = os.Remove("/var/lib/nanokube/state/last-event")
	e2etest.AssertFileAbsent(s.T(), "/var/lib/nanokube/state/last-event", "reset state")
}
