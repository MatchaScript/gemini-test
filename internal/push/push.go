package push

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/MatchaScript/nanokube/contract/desiredpb"
	"github.com/MatchaScript/nanokube/internal/ddi"
	"github.com/MatchaScript/nanokube/internal/operator"
	"github.com/MatchaScript/nanokube/internal/render"
)

// DesiredToAgent builds a confext DDI image from the given desired document
// and pushes it over gRPC to the specified agent address.
func DesiredToAgent(ctx context.Context, d render.Desired, fileContexts string, agentAddr string) error {
	tmpFile, err := os.CreateTemp("", "nanokube-ddi-*.raw")
	if err != nil {
		return fmt.Errorf("create temp ddi file: %w", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	os.Remove(tmpPath) // systemd-repart requires destination file not to exist yet
	defer os.Remove(tmpPath)

	input := ddi.BuildInput{
		Name:             d.Name(),
		Files:            d.Files,
		FileContextsPath: fileContexts,
	}

	if err := ddi.Build(input, tmpPath); err != nil {
		return fmt.Errorf("build ddi: %w", err)
	}

	blobBytes, err := os.ReadFile(tmpPath)
	if err != nil {
		return fmt.Errorf("read built ddi: %w", err)
	}

	h := sha256.Sum256(blobBytes)
	sha256Hex := hex.EncodeToString(h[:])

	meta := &desiredpb.DesiredMetadata{
		Name:       d.Name(),
		BlobSha256: sha256Hex,
	}

	pushFn := operator.NewGRPCPush(agentAddr)
	if err := pushFn(ctx, meta, tmpPath); err != nil {
		return fmt.Errorf("push to agent (%s): %w", agentAddr, err)
	}

	return nil
}
