package cmd

import (
	"fmt"
	"net/url"
	"os"
	"regexp"

	"github.com/ptrkrlsrd/acache/pkg/acache"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use: "add",
	Short: `Add a new route. 
		Example: "acache add https://pokeapi.co/api/v2/pokemon/ditto /ditto"
		Here the first argument is the path to the endpoint you want to cache, 
		and the last is the alias. Note that you can also add from a json file by replacing 
        the first URL with a relative path to a json file.`,

	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		alias := args[1]

        if isURL(path) {
            route, err := acache.NewRouteFromURL(path, alias)
            HandleError(err)

            err = server.AddRoute(route)
            HandleError(err)
        } else if isFile(path) {
            route, err := acache.NewRouteFromFile(path, alias)
            HandleError(err)

            err = server.AddRoute(route)
            HandleError(err)
        }
	},
}


func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func isFile(str string) bool {
	match, _ := regexp.MatchString(`^(?:(?:[a-zA-Z]:)?[/\\]{0,2})?(?:\.{1,2}[/\\])?[\w\s-]+(?:[/\\][\w\s-]+)*(?:\.[\w-]+)?$`, str)
	if !match {
		return false
	}

	_, err := os.Stat(str)
	return !os.IsNotExist(err)
}
