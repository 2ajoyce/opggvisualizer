// internal/cli/db.go
package cli

import (
	"opggvisualizer/internal/db"

	"github.com/spf13/cobra"
)

func newDBClearChampionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "champions wipe",
		Short: "Removes all champion data from the database",
		Run: func(cmd *cobra.Command, args []string) {
			database := db.GetDatabaseConnection()
			err := database.ClearChampionData()
			if err != nil {
				cmd.PrintErrf("Error clearing champion data: %v\n", err)
				return
			}
		},
	}
	return cmd
}

func newDBClearGamesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "games wipe",
		Short: "Removes all Game and Participant data from the database",
		Run: func(cmd *cobra.Command, args []string) {
			database := db.GetDatabaseConnection()
			err := database.ClearGameData()
			if err != nil {
				cmd.PrintErrf("Error clearing game data: %v\n", err)
				return
			}
		},
	}
	return cmd
}
