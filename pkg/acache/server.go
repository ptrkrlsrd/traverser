package acache

import (
	"fmt"
	"log"

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

// routeHandler handles the route
func routeHandler(route Route) func(c *gin.Context) {
	return func(c *gin.Context) {
		for header, v := range route.Response.Headers {
			c.Header(header, v)
		}

		c.String(route.Response.StatusCode, string(route.Response.Body))
	}
}

func containsRoute(url string, routes gin.RoutesInfo) bool {
	for _, v := range routes {
		if v.Path == url {
			return true
		}
	}

	return false
}

// RegisterRoutes registers the stored routes to the server
func (server *Server) RegisterRoutes(routes Routes) {
	for _, r := range routes {
		if containsRoute(r.Alias, server.router.Routes()) {
			continue
		}
		server.router.GET(r.Alias, routeHandler(r))

	}
}

// AddRoute adds a route to the server
func (server *Server) AddRoute(route Route) error {
	return server.store.AddRoute(route)
}

// LoadRoutes loads the routes from the store
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

func (server *Server) ClearDatabase() error {
	return server.store.Clear()
}

// Start starts the API server
func (server *Server) Start() error {
	return server.router.Run(fmt.Sprintf(":%d", server.port))
}
