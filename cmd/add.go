package cmd

import (
	"github.com/ptrkrlsrd/acache/pkg/acache"
	"github.com/spf13/cobra"
)

var (
	postData string
)

func init() {
	addCmd.PersistentFlags().StringVarP(&postData, "data", "i", "", "")
	rootCmd.AddCommand(addCmd)
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use: "add",
	Short: `Add a new route. 
		Example: "acache add https://pokeapi.co/api/v2/pokemon/ditto /ditto"
		Here the first argument is the path to the endpoint you want to cache, 
		and the last is the alias`,

	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		alias := args[1]

		route, err := acache.NewRouteFromURL(url, alias)
		HandleError(err)

		err = server.AddRoute(route)
		HandleError(err)
	},
}
