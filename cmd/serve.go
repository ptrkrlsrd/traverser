package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	port int
)

// serveCmd serves the stored routes
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Load the stored routes from cache and serve the API",
	Run: func(cmd *cobra.Command, args []string) {
		server.UsePort(port)
		server.PrintRoutes()
		server.LoadRoutes()

		log.Printf("Started server on port: %d\n", port)
		err := server.Start()
		HandleError(err)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntVarP(&port, "port", "p", 4000, "Port")
}
