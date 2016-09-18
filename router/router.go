// Package router provides multiplex routers.
package router

import (
	"net/http"

	"github.com/blueprint/blueprint/resource"
)

// Router represents an HTTP router and impliments the http.Handler interface.
type Router interface {

	// Resource controller action methods.
	GET(uri string, f http.HandlerFunc, m ...Middleware)
	POST(uri string, f http.HandlerFunc, m ...Middleware)
	PUT(uri string, f http.HandlerFunc, m ...Middleware)
	PATCH(uri string, f http.HandlerFunc, m ...Middleware)
	DELETE(uri string, f http.HandlerFunc, m ...Middleware)

	// Middleware adds HTTP middleware to the router.
	Middleware(m ...Middleware) Router

	// Flush the current MiddlewareChain.
	FlushMiddleware(path string) Router

	// Subrouting for API versioning support.
	Subrouter(path string) Router

	// Serve static files.
	Static(uri, dir string)

	// Impliment the http.Handler interface.
	ServeHTTP(w http.ResponseWriter, req *http.Request)

	// Generate routes from a resource controller.
	Resource(path, id string, rs resource.Resourcer, m ...Middleware)
}
