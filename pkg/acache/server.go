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
	"net/http"
	"net/http/httputil"
	"net/url"
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

// NewServer creates a new NewServer
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

// RegisterStoredRoutes registers the stored routes to the server
func (server *Server) RegisterStoredRoutes() {
	for _, r := range server.Storage.Routes {
		handler := func(c *gin.Context) {
			for header, v := range r.Response.Header {
				c.Header(header, strings.Join(v, ","))
			}

			c.String(r.Response.StatusCode, string(r.Response.Body))
		}

		server.router.GET(r.Alias, handler)
	}
}

// UseProxyRoute creates a catch all route which intercepts all API calls
func (server *Server) proxyRoute(proxyURL string) func(*gin.Context) {
	proxyHandler := func(c *gin.Context) {
		remote, err := url.Parse(proxyURL)
		if err != nil {
			c.AbortWithError(500, err)
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = c.Param("proxyPath")

			if server.Storage.Routes.ContainsURL(proxyURL) {
				return
			}

			replacedURL := fmt.Sprintf("%s%s", proxyURL, c.Request.URL.Path)
			route, err := NewRouteFromRequest(replacedURL, c.Request.URL.Path)
			if err != nil {
				c.AbortWithError(500, err)
				return
			}

			server.Storage.AddRoute(route)
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}

	return proxyHandler
}

func (server *Server) RegisterProxyHandler(proxyURL string) {
	server.router.NoRoute(server.proxyRoute(proxyURL))
}

//StartServer starts the API server
func (server *Server) StartServer() error {
	return server.router.Run(fmt.Sprintf(":%d", server.port))
}
