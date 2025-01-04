// internal/cli/cli.go
package cli

import (
	"context"

	"github.com/spf13/cobra"
)

func NewRootCommand(ctx context.Context) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "opggvisualizer",
		Short: "A tool to visualize League of Legends game data",
	}

	// Add subcommands
	rootCmd.AddCommand(newFetchChampionsCommand()) // TODO: Add context to other commands
	rootCmd.AddCommand(newFetchGamesCommand())
	rootCmd.AddCommand(newStartAPICmd(ctx))
	rootCmd.AddCommand(newStopAPICmd(ctx))
	rootCmd.AddCommand(newDBClearChampionsCmd())
	rootCmd.AddCommand(newDBClearGamesCmd())

	return rootCmd
}
