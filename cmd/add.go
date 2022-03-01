// Copyright © 2021 github.com/ptrkrlsrd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

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

		route, err := acache.NewRouteFromRequest(url, alias)
		HandleError(err)

		if err := server.Storage.AddRoute(route); err != nil {
			HandleError(fmt.Errorf("error adding route: %v", err))
		}
	},
}
