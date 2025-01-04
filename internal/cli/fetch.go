// internal/cli/fetch.go
package cli

import (
	"log"

	"opggvisualizer/internal/client"

	"github.com/spf13/cobra"
)

func newFetchChampionsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "champions fetch",
		Short: "Fetch and store champion data",
		Run: func(cmd *cobra.Command, args []string) {
			if err := client.FetchAndStoreChampionData(); err != nil {
				log.Fatalf("Error fetching and storing champion data: %v", err)
			}
			log.Println("Data fetching and insertion completed successfully.")
		},
	}
	return cmd
}

func newFetchGamesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "games fetch",
		Short: "Fetch and store game data",
		Run: func(cmd *cobra.Command, args []string) {
			if err := client.FetchAndStoreGameData(); err != nil {
				log.Fatalf("Error fetching and storing game data: %v", err)
			}
			log.Println("Data fetching and insertion completed successfully.")
		},
	}
	return cmd
}
