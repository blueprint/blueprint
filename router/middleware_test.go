package router

import (
	"net/http"
	"testing"
)

var state = ""

type middleware1 struct {
	success http.Handler
}

func (m1 middleware1) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	state = "m1"
	m1.success.ServeHTTP(resp, req)
}

func (m1 middleware1) Do(h http.Handler) http.Handler {
	m1.success = h
	return m1
}

type middleware2 struct {
	success http.Handler
}

func (m2 middleware2) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	state = "m2"
	m2.success.ServeHTTP(resp, req)
}

func (m2 middleware2) Do(h http.Handler) http.Handler {
	m2.success = h
	return m2
}

type middleware3 struct {
	success http.Handler
}

func (m3 middleware3) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	state = "m3"
	m3.success.ServeHTTP(resp, req)
}

func (m3 middleware3) Do(h http.Handler) http.Handler {
	m3.success = h
	return m3
}

type middleware4 func(resp http.ResponseWriter, req *http.Request)

func (m4 middleware4) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	state = "m4"
}

func TestMiddlewareChain(t *testing.T) {

	m1 := middleware1{}
	m2 := middleware2{}
	mc1 := NewMiddlewareChain(m1.Do, m2.Do)

	if len(mc1.middleware) < 1 {
		t.Error("expected MiddlewareChain not to be empty")
	}

	m3 := middleware3{}

	mc1.Then(m3)

	mc1.ThenFunc(func(resp http.ResponseWriter, req *http.Request) {})

	mc2 := NewMiddlewareChain()
	mc2.Append(mc1.middleware...)

}
