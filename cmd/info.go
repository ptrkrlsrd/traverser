package cmd

import (
	"github.com/spf13/cobra"
)

// infoCmd prints route information
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print route information",
	Run: func(cmd *cobra.Command, args []string) {
		routes, err := server.Store.LoadRoutes()
		if err != nil {
			HandleError(err)
		}
		routes.PrintInfo()
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
