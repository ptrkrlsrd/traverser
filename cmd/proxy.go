package cmd

import (
	"fmt"
	"log"
	"regexp"

	"github.com/spf13/cobra"
)

// proxyCmd proxies the stored routes
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Start Acache as a proxy between you and another API and save the responses locally",
	Long: `Start Acache as a proxy between you and another API and save the responses locally. 
	Example: acache proxy https://pokeapi.co/api/v2/pokemon.`,
	Args: cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		proxyURL := args[0]
		re := regexp.MustCompile(`^(http|https)://.+`)
		if !re.MatchString(proxyURL) {
			HandleError(fmt.Errorf("invalid URL: %s", proxyURL))
		}

		server.UsePort(port)
		server.RegisterProxyRoute(proxyURL)
		log.Printf("Started server on port: %d\n", port)
		err := server.Start()
		HandleError(err)
	},
}

func init() {
	if enableExperimenalFeatures {
		rootCmd.AddCommand(proxyCmd)
	}
	proxyCmd.Flags().IntVarP(&port, "port", "p", 4000, "Port")
}
