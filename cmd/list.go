package cmd

import (
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all routes",
	Aliases: []string{"ls", "l"},
	Run: func(cmd *cobra.Command, args []string) {
		server.PrintRoutes()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
