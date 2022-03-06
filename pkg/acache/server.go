package acache

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// DefaultPort is the default port to listen on
	DefaultPort = 4000
)

// Server contains the dependencies and handles the logic
type Server struct {
	store  RouteStorer
	port   int
	router *gin.Engine
}

// NewServer creates a new server
func NewServer(store RouteStorer, router *gin.Engine) Server {
	return Server{
		store:  store,
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

// AddRoute adds a route to the server
func (server *Server) AddRoute(route Route) error {
	return server.store.AddRoute(route)
}

func (server *Server) LoadRoutes() error {
	routes, err := server.store.GetRoutes()
	if err != nil {
		return err
	}

	server.RegisterRoutes(routes)
	return nil

}

// PrintRoutes prints all routes to the console
func (server *Server) PrintRoutes() {
	routes, err := server.store.GetRoutes()
	if err != nil {
		log.Println(err)
	}
	routes.Print()
}

// PrintRouteInfo prints the route info to the console
func (server *Server) PrintRouteInfo() {
	routes, err := server.store.GetRoutes()
	if err != nil {
		log.Println(err)
	}
	routes.PrintInfo()
}

//Start starts the API server
func (server *Server) Start() error {
	return server.router.Run(fmt.Sprintf(":%d", server.port))
}
