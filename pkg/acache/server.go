// Copyright © 2020 github.com/ptrkrlsrd
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
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPort = 3000
)

// Server contains the dependencies and handles the logic
type Server struct {
	Storage Storage
	port    int
	router  *gin.Engine
}

// NewServer creates a new SNewServer
func NewServer(storage Storage, router *gin.Engine) Server {
	return Server{
		Storage: storage,
		router:  router,
		port:    DefaultPort,
	}
}

func (server *Server) UsePort(port int) {
	server.port = port
}

// UseStoredRoutes registers the stored routes to the server
func (server *Server) UseStoredRoutes() {
	for _, r := range server.Storage.Routes {
		server.router.GET(r.Alias, func(c *gin.Context) {
			for header, v := range r.Response.Header {
				values := strings.Join(v, ",")
				c.Header(header, values)
			}

			c.String(r.Response.StatusCode, string(r.Response.Body))
		})
	}
}

//StartServer starts the API server
func (server *Server) StartServer() error {
	return server.router.Run(fmt.Sprintf(":%d", server.port))
}
