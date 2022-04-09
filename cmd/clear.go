package cmd

import (
	"github.com/spf13/cobra"
)

// clearCmd proxies the stored routes
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clears the database containing the stored routes",
	Run: func(cmd *cobra.Command, args []string) {
		err := server.ClearDatabase()
		HandleError(err)
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
