package main

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/MatchaScript/nanokube/internal/config"
	"github.com/MatchaScript/nanokube/internal/initialize"
	"github.com/MatchaScript/nanokube/internal/state"
)

type initOpts struct {
	agentAddr string
}

func newInitCmd(g *globalOpts) *cobra.Command {
	opts := &initOpts{
		agentAddr: "127.0.0.1:50051",
	}
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialise a fresh node (run once per install)",
		Long: "Mirrors `kubeadm init`'s scope: renders control-plane configuration, " +
			"PKI, kubeconfigs, and static pod manifests into a confext DDI, " +
			"and pushes it to nanokube-agent over gRPC.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			existed, err := state.Exists(g.layout)
			if err != nil {
				return err
			}
			if existed {
				return errors.New("nanokube state already exists; run `nanokube reset --yes` first to re-initialise")
			}
			loaded, err := config.Load(g.configPath, g.layout)
			if err != nil {
				return err
			}
			if loaded.HasJoin {
				return fmt.Errorf("config %s describes a joined node (JoinConfiguration); `nanokube init` bootstraps a new cluster — use `nanokube add-node`", g.configPath)
			}
			return initialize.Run(cmd.Context(), loaded.Init, g.layout, g.fileContexts, opts.agentAddr, cmd.OutOrStdout())
		},
	}
	cmd.Flags().StringVar(&opts.agentAddr, "agent-addr", opts.agentAddr, "nanokube-agent gRPC address")
	return cmd
}
