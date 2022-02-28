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
		routes, err := server.Storage.LoadRoutes()
		if err != nil {
			HandleError(err)
		}
		routes.Print()
		server.MapRoutes(routes)
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
