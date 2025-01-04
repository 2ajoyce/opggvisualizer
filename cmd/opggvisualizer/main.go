// cmd/opggvisualizer/main.go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"opggvisualizer/internal/cli"
)

func main() {
	// Create a context that is cancelled on interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Initialize CLI commands
	rootCmd := cli.NewRootCommand(ctx)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
