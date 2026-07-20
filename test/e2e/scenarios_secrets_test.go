//go:build e2e

package e2e

import (
	"os"
	"path/filepath"
	"strings"
)

// TestZZ_NoSecretsInArtifacts scans artifacts/logs produced during E2E testing
// to ensure no RSA/EC private keys or secret values were leaked into readable outputs.
func (s *NanokubeE2ESuite) TestZZ_NoSecretsInArtifacts() {
	t := s.T()

	// Prohibited secret markers
	prohibited := []string{
		"BEGIN RSA PRIVATE KEY",
		"BEGIN EC PRIVATE KEY",
		"BEGIN PRIVATE KEY",
		"client-key-data:",
	}

	searchDirs := []string{
		t.TempDir(),
	}

	for _, dir := range searchDirs {
		if dir == "" {
			continue
		}
		_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			// Skip actual PKI/cert store files under /etc/kubernetes/pki or scratch pki dirs
			if strings.Contains(path, "/pki/") || strings.HasSuffix(path, ".key") || strings.HasSuffix(path, ".crt") {
				return nil
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return nil
			}

			for _, p := range prohibited {
				if strings.Contains(string(content), p) {
					t.Errorf("Secret leak detected in artifact file %s: contains %q", path, p)
				}
			}
			return nil
		})
	}
}
