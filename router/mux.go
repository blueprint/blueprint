package router

import (
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/blueprint/blueprint/resource"
	"github.com/julienschmidt/httprouter"
)

const (
	JSON    = "json"
	YML     = "yml"
	XML     = "xml"
	dotJSON = ".json"
	dotYML  = ".yml"
	dotXML  = ".xml"
)

// Mux describes a HTTP multiplex router.
type Mux struct {
	mux    *httprouter.Router
	mc     MiddlewareChain
	prefix string
}

// NewMux returns a new router.
func NewMux() *Mux {

	// create a new HTTP multiplexer
	mux := httprouter.New()
	mux.HandleMethodNotAllowed = false

	// create a new router
	return &Mux{mux: mux, mc: NewMiddlewareChain()}
}

// Middleware registers HTTP middlware.
func (m *Mux) Middleware(middleware ...Middleware) Router {
	if len(middleware) > 0 {
		m.mc = m.mc.Append(middleware...)
	}
	return m
}

// FlushMiddleware creates a subrouter with a fresh middleware chain.
func (m *Mux) FlushMiddleware(path string) Router {
	return &Mux{mux: m.mux, prefix: m.prefix + path}
}

// ServeHTTP satisfies the http.Hander interface. This provides flexiblity and
// compatiblity with the standard http package.
func (m *Mux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m.mux.ServeHTTP(w, req)
}

// Subrouter creates a new subrouter based on the path prefix and middlweare of its parent.
func (m *Mux) Subrouter(path string) Router {
	return &Mux{mux: m.mux, mc: m.mc, prefix: m.prefix + path}
}

// GET registers a URL path with an Action.
func (m *Mux) GET(uri string, f http.HandlerFunc, middleware ...Middleware) {
	url := path.Join(m.prefix, uri)
	// if the path ends with a param use a dynamic formatted action, otherwise
	// statically define the routes for better performance
	if paramEnd(uri) {
		m.mux.GET(url, m.action(f, middleware...))
	} else {
		m.mux.GET(url, m.actionWithFormat(JSON, f, middleware...))
		m.mux.GET(url+dotJSON, m.actionWithFormat(JSON, f, middleware...))
		m.mux.GET(url+dotXML, m.actionWithFormat(XML, f, middleware...))
		m.mux.GET(url+dotYML, m.actionWithFormat(YML, f, middleware...))
	}
	endpoints = append(endpoints, Endpoint{Type: "GET", Path: url})
}

// POST registers a URL path with an Action.
func (m *Mux) POST(uri string, f http.HandlerFunc, mw ...Middleware) {
	url := path.Join(m.prefix, uri)
	// if the path ends with a param use a dynamic formatted action, otherwise
	// statically define the routes for better performance
	if paramEnd(uri) {
		m.mux.POST(url, m.action(f, mw...))
	} else {
		m.mux.POST(url, m.actionWithFormat(JSON, f, mw...))
		m.mux.POST(url+dotJSON, m.actionWithFormat(JSON, f, mw...))
		m.mux.POST(url+dotXML, m.actionWithFormat(XML, f, mw...))
		m.mux.POST(url+dotYML, m.actionWithFormat(YML, f, mw...))
	}
	endpoints = append(endpoints, Endpoint{Type: "POST", Path: url})
}

// PATCH registers a URL path with an Action.
func (m *Mux) PATCH(uri string, f http.HandlerFunc, mw ...Middleware) {
	url := path.Join(m.prefix, uri)
	// if the path ends with a param use a dynamic formatted action, otherwise
	// statically define the routes for better performance
	if paramEnd(uri) {
		m.mux.PATCH(url, m.action(f, mw...))
	} else {
		m.mux.PATCH(url, m.actionWithFormat(JSON, f, mw...))
		m.mux.PATCH(url+dotJSON, m.actionWithFormat(JSON, f, mw...))
		m.mux.PATCH(url+dotXML, m.actionWithFormat(XML, f, mw...))
		m.mux.PATCH(url+dotYML, m.actionWithFormat(YML, f, mw...))
	}
	endpoints = append(endpoints, Endpoint{Type: "PATCH", Path: url})
}

// PUT registers a URL path with an Action.
func (m *Mux) PUT(uri string, f http.HandlerFunc, mw ...Middleware) {
	url := path.Join(m.prefix, uri)
	// if the path ends with a param use a dynamic formatted action, otherwise
	// statically define the routes for better performance
	if paramEnd(uri) {
		m.mux.PUT(url, m.action(f, mw...))
	} else {
		m.mux.PUT(url, m.actionWithFormat(JSON, f, mw...))
		m.mux.PUT(url+dotJSON, m.actionWithFormat(JSON, f, mw...))
		m.mux.PUT(url+dotXML, m.actionWithFormat(XML, f, mw...))
		m.mux.PUT(url+dotYML, m.actionWithFormat(YML, f, mw...))
	}
	endpoints = append(endpoints, Endpoint{Type: "PUT", Path: url})
}

// DELETE registers a URL path with an Action.
func (m *Mux) DELETE(uri string, f http.HandlerFunc, mw ...Middleware) {
	url := path.Join(m.prefix, uri)
	// if the path ends with a param use a dynamic formatted action, otherwise
	// statically define the routes for better performance
	if paramEnd(uri) {
		m.mux.DELETE(url, m.action(f, mw...))
	} else {
		m.mux.DELETE(url, m.actionWithFormat(JSON, f, mw...))
		m.mux.DELETE(url+dotJSON, m.actionWithFormat(JSON, f, mw...))
		m.mux.DELETE(url+dotXML, m.actionWithFormat(XML, f, mw...))
		m.mux.DELETE(url+dotYML, m.actionWithFormat(YML, f, mw...))
	}
	endpoints = append(endpoints, Endpoint{Type: "DELETE", Path: url})
}

// action is a private HTTP handler that executes a controller method.
//
// action also takes multiple negroni.Handler objects to create the middleware
// chain for a route.
func (m *Mux) action(h http.Handler, mw ...Middleware) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var format string
		if len(ps) > 0 {
			v := ps[len(ps)-1].Value
			l := len(v)
			if l > 4 {
				switch v[l-4:] {
				case dotXML:
					format = XML
				case dotYML:
					format = YML
				default:
					format = JSON
				}
			}
		}

		resp := http.NewResponse(w, format)
		req := http.NewRequest(r, ps)

		if len(mw) > 0 {
			chain := m.mc.Append(mw...).Then(h)
			chain.ServeHTTP(resp, req)
		} else {
			chain := m.mc.Then(h)
			chain.ServeHTTP(resp, req)
		}
	}
}

// paramEnd checks if path ends in a parameter.
func paramEnd(path string) bool {
	for i := len(path) - 2; i >= 0; i-- {
		if path[i] == ':' {
			return true
		}
		if path[i] == '/' {
			return false
		}
	}
	return false
}

// actionWithFormat executes an action method, but also specifies the return format.
func (m *Mux) actionWithFormat(format string, h http.Handler, mw ...Middleware) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		resp := http.NewResponse(w, format)
		req := http.NewRequest(r, ps)

		if len(mw) > 0 {
			chain := m.mc.Append(mw...).Then(h)
			chain.ServeHTTP(resp, req)
		} else {
			chain := m.mc.Then(h)
			chain.ServeHTTP(resp, req)
		}
	}
}

// Resource registers a URL path with a Controller that impliments all
// Index, Store, Show, Update, Apply, and Destory Actions.
func (m *Mux) Resource(path, id string, r resource.Resourcer, mw ...Middleware) {
	pathID := m.prefix + path + "/:" + id
	m.GET(path, r.Index, mw...)
	m.GET(pathID, r.Show, mw...)
	m.POST(path, r.Store, mw...)
	m.PUT(pathID, r.Update, mw...)
	m.PATCH(pathID, r.Apply, mw...)
	m.DELETE(pathID, r.Destroy, mw...)
}

// Static registers a URL path with a public directory to serve its content.
// This directory is meant to serve public static files such as image files,
// CSS files, JavaScript files, and more.
func (m *Mux) Static(uri, dir string) {
	uri = path.Clean(uri)
	if uri == "/" {
		uri = ""
	}

	url := path.Join(m.prefix, uri)
	files := path.Join(url, "*filepath")

	// check directory path is valid
	if _, err := os.Stat(dir); err != nil {

		// since dir is not a valid directory path, assume the path given is
		// relative to the binary executing
		current, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatalln(err)
		}

		dir = filepath.Join(current, dir)

		// serve all files in the directory
		m.mux.ServeFiles(files, http.Dir(dir))
		return
	}

	// serve all files in the directory
	m.mux.ServeFiles(files, http.Dir(dir))
}
