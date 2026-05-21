package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewReverseProxy(target string) (*httputil.ReverseProxy, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	origDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		origDirector(req)
		// Preserve original path and query
		req.Host = u.Host
	}
	return proxy, nil
}

func ProxyHandler(target string) gin.HandlerFunc {
	p, err := NewReverseProxy(target)
	if err != nil {
		panic(err)
	}
	return func(c *gin.Context) {
		// Adjust path: keep everything as-is relative to the target
		c.Request.URL.Path = c.Param("path")
		if !strings.HasPrefix(c.Request.URL.Path, "/") {
			c.Request.URL.Path = "/" + c.Request.URL.Path
		}
		p.ServeHTTP(c.Writer, c.Request)
	}
}
