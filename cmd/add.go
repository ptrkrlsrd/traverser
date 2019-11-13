// Copyright © 2018 Petter Karlsrud petterkarlsrud@me.com
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
	"log"

	"github.com/ptrkrlsrd/acache/pkg/acache"
	"github.com/spf13/cobra"
)

var (
	validHTTPMethods = []string{"GET", "POST"}
	selectedHTTPMode = "GET"
	postData         string
)

func init() {
	addCmd.PersistentFlags().StringVarP(&selectedHTTPMode, "method", "m", "GET", "")
	addCmd.PersistentFlags().StringVarP(&postData, "data", "i", "", "")
	rootCmd.AddCommand(addCmd)
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use: "add",
	Short: `Add a new route. 
		Example: "acache add https://api.coinmarketcap.com/v1/ticker/ /eth"
		Here the first argument is the path to the endpoint you want to cache, 
		and the last is the alias`,

	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		alias := args[1]

		var route acache.Route
		var err error

		switch selectedHTTPMode {
		case "GET":
			route, err = service.NewRouteFromGetRequest(url, alias)
			HandleError(err)
		case "POST":
			if postData == "" {
				log.Fatal("")
			}

			route, err = service.NewRouteFromPostRequest(url, alias, []byte(postData))
			HandleError(err)
		}

		if err := service.StoreRoute(route); err != nil {
			HandleError(fmt.Errorf("error adding route: %v", err))
		}
	},
}
