// internal/cli/api.go
package cli

import (
	"context"
	"opggvisualizer/internal/api"

	"github.com/spf13/cobra"
)

func newStartAPICmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server start",
		Short: "Start the API server",
		Run: func(cmd *cobra.Command, args []string) {
			api.Start(ctx)
		},
	}
	return cmd
}

func newStopAPICmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server stop",
		Short: "Stop the API server",
		Run: func(cmd *cobra.Command, args []string) {
			api.Stop(ctx)
		},
	}
	return cmd
}
