package acache

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

type ProxyRoute struct {
	OriginalURL string
	Alias       string
}

// UseProxyRoute creates a catch all route which intercepts all API calls
func (server *Server) proxyHandleFunc(proxyURL string, newRouteChan chan (ProxyRoute)) func(*gin.Context) {
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

			proxyPath := c.Request.RequestURI
			req.URL.Path = proxyPath

			newRouteChan <- ProxyRoute{req.URL.String(), c.Request.URL.Path}
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}

	return proxyHandler
}

func listenAndHandleProxyRoutes(proxyRouteChan chan (ProxyRoute), store RouteStorer) {
	for {
		proxyRoute := <-proxyRouteChan
		route, err := NewRouteFromURL(proxyRoute.OriginalURL, proxyRoute.Alias)
		if err != nil {
			log.Fatal(err)
		}

		store.AddRoute(route)
	}
}

func (server *Server) RegisterProxyHandler(proxyURL string) {
	proxyRouteChan := make(chan ProxyRoute)
	go listenAndHandleProxyRoutes(proxyRouteChan, server.store)
	server.router.NoRoute(server.proxyHandleFunc(proxyURL, proxyRouteChan))
}
