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
	DefaultPort = 4000
)

// Server contains the dependencies and handles the logic
type Server struct {
	Store  Store
	port   int
	router *gin.Engine
}

// NewServer creates a new server
func NewServer(store Store, router *gin.Engine) Server {
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

// UseProxyRoute creates a catch all route which intercepts all API calls
func (server *Server) RegisterProxyRoute(proxyURL string) {
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

			originalURL := c.Request.URL.Path
			replacedURL := fmt.Sprintf("%s%s", proxyURL, originalURL)
			route, err := NewRouteFromURL(replacedURL, originalURL)
			if err != nil {
				c.AbortWithError(500, err)
				return
			}

			server.Store.AddRoute(route)
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}

	server.router.NoRoute(proxyHandler)
}

//Start starts the API server
func (server *Server) Start() error {
	return server.router.Run(fmt.Sprintf(":%d", server.port))
}
