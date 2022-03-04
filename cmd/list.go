package cmd

import (
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all routes(aliases)",
	Aliases: []string{"ls", "l"},
	Run: func(cmd *cobra.Command, args []string) {
		routes, err := server.Store.GetRoutes()
		HandleError(err)
		routes.Print()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
