package router

// Endpoints contains a list of Endpoints.
var endpoints []Endpoint

// Endpoint represents an API URI endpoint for documentation purposes.
type Endpoint struct {
	Type string
	Path string
}

// Endpoints returns the list of registered HTTP endpoints.
func Endpoints() []Endpoint {
	return endpoints
}
