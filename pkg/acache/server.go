package acache

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPort = 4000
)

// Server contains the dependencies and handles the logic
type Server struct {
	Store  RouteStorer
	port   int
	router *gin.Engine
}

// NewServer creates a new server
func NewServer(store RouteStorer, router *gin.Engine) Server {
	return Server{
		Store:  store,
		router: router,
		port:   DefaultPort,
	}
}

// UsePort sets the port to listen on
func (server *Server) UsePort(port int) {
	server.port = port
}

// UseStoredRoutes registers the stored routes to the server
func (server *Server) RegisterRoutes(routes Routes) {
	for _, r := range routes {
		handler := func(c *gin.Context) {
			for header, v := range r.Response.Header {
				c.Header(header, strings.Join(v, ","))
			}

			c.String(r.Response.StatusCode, string(r.Response.Body))
		}

		server.router.GET(r.Alias, handler)
	}
}

//Start starts the API server
func (server *Server) Start() error {
	return server.router.Run(fmt.Sprintf(":%d", server.port))
}
