// Copyright © 2020 Petter Karlsrud
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

package acache

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Server contains the dependencies and handles the logic
type Server struct {
	Storage Storage
}

// NewRoute creates a new route from an URL and an alias
func (service *Server) NewRoute(url string, alias string) (Route, error) {
	res, err := http.Get(url)
	if err != nil {
		return Route{}, err
	}

	route := NewRoute(url, alias, http.MethodGet, res)
	return route, err
}

// StoreRoute stores a route to the Bolt DB
func (server *Server) StoreRoute(route Route) error {
	jsonData, err := json.Marshal(route)
	if err != nil {
		return fmt.Errorf("failed marshaling JSON: %v", err)
	}

	return server.Storage.Add(route.ID, jsonData)
}

//StartServer starts the API server
func (service *Server) StartServer(addr string) error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	for _, r := range service.Storage.Routes {
		router.GET(r.Alias, func(c *gin.Context) {
			for header, v := range r.Response.Header {
				values := strings.Join(v, ",")
				c.Header(header, values)
			}

			c.String(r.Response.StatusCode, string(r.Response.Body))
		})
	}

	return router.Run(addr)
}
