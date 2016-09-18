package router

import (
	"net/http"
	"testing"
)

func TestNewMux(t *testing.T) {
	mux := NewMux()
	if mux == nil {
		t.Error("expected mux not to be nil")
	}
}

func TestMuxSubrouter(t *testing.T) {
	mux := NewMux()
	sub := mux.Subrouter("/foo")
	if sub == nil {
		t.Error("expected sub not to be nil")
	}
}

func TestMuxGet(t *testing.T) {
	mux := NewMux()
	mux.GET("/some/endpoint", func(resp http.ResponseWriter, req *http.Request) {})
	mux.GET("/some/endpoint/:id", func(resp http.ResponseWriter, req *http.Request) {})
}

func TestMuxPost(t *testing.T) {
	mux := NewMux()
	mux.POST("/some/endpoint", func(resp http.ResponseWriter, req *http.Request) {})
	mux.POST("/some/endpoint/:id", func(resp http.ResponseWriter, req *http.Request) {})
}

func TestMuxPATCH(t *testing.T) {
	mux := NewMux()
	mux.PATCH("/some/endpoint", func(resp http.ResponseWriter, req *http.Request) {})
	mux.PATCH("/some/endpoint/:id", func(resp http.ResponseWriter, req *http.Request) {})
}

func TestMuxPUT(t *testing.T) {
	mux := NewMux()
	mux.PUT("/some/endpoint", func(resp http.ResponseWriter, req *http.Request) {})
	mux.PUT("/some/endpoint/:id", func(resp http.ResponseWriter, req *http.Request) {})
}

func TestMuxDELETE(t *testing.T) {
	mux := NewMux()
	mux.DELETE("/some/endpoint", func(resp http.ResponseWriter, req *http.Request) {})
	mux.DELETE("/some/endpoint/:id", func(resp http.ResponseWriter, req *http.Request) {})
}
