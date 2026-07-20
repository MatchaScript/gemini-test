package main

import (
	"github.com/spf13/cobra"

	"github.com/MatchaScript/nanokube/internal/addnode"
)

func newAddNodeCmd(g *globalOpts) *cobra.Command {
	var opts addnode.Options
	opts.AgentAddr = "127.0.0.1:50051"

	cmd := &cobra.Command{
		Use:   "add-node",
		Short: "Join this node to an existing cluster as a worker",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.FileContexts = g.fileContexts
			return addnode.Run(cmd.Context(), opts, g.layout, cmd.OutOrStdout())
		},
	}
	cmd.Flags().StringVar(&opts.Server, "server", "", "reachable control-plane apiserver (https://host:port)")
	cmd.Flags().StringVar(&opts.Token, "token", "", "bootstrap token from `nanokube token create`")
	cmd.Flags().StringSliceVar(&opts.CACertHashes, "ca-cert-hash", nil, "CA public key pin (sha256:...) from `nanokube token create`")
	cmd.Flags().StringVar(&opts.AgentAddr, "agent-addr", opts.AgentAddr, "nanokube-agent gRPC address")
	_ = cmd.MarkFlagRequired("server")
	_ = cmd.MarkFlagRequired("token")
	_ = cmd.MarkFlagRequired("ca-cert-hash")
	return cmd
}
