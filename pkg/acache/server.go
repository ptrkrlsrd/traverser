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
	"strings"

	"github.com/gin-gonic/gin"
)

// Server contains the dependencies and handles the logic
type Server struct {
	Storage Storage
}

//StartServer starts the API server
func (server *Server) StartServer(addr string) error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	for _, r := range server.Storage.Routes {
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
