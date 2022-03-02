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
		routes, err := server.Store.LoadRoutes()
		if err != nil {
			HandleError(err)
		}
		routes.Print()
		server.RegisterRoutes(routes)
		log.Printf("Started server on port: %d\n", port)
		if err := server.StartServer(); err != nil {
			HandleError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntVarP(&port, "port", "p", 4000, "Port")
}
