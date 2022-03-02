package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// proxyCmd proxies the stored routes
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Start the server as a proxy between you and another API",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
		proxyURL := args[0]

		routes, err := server.Store.LoadRoutes()
		if err != nil {
			HandleError(err)
		}
		routes.Print()
		server.UsePort(port)
		server.RegisterProxyRoute(proxyURL)
		log.Printf("Started server on port: %d\n", port)
		if err := server.StartServer(); err != nil {
			HandleError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)
	proxyCmd.Flags().IntVarP(&port, "port", "p", 4000, "Port")
}
