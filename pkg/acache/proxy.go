package acache

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

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

			replacedURL := fmt.Sprintf("%s%s", proxyURL, c.Request.URL.Path)
			route, err := NewRouteFromURL(replacedURL, c.Request.URL.Path)
			if err != nil {
				c.AbortWithError(500, err)
				return
			}

			server.Store.AddRoute(route)
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}

	return proxyHandler
}

func (server *Server) RegisterProxyHandler(proxyURL string) {
	server.router.NoRoute(server.proxyRoute(proxyURL))
}
