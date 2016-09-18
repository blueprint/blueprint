package router

import (
	"log"
	"net/http"
)

// Middleware represents a function that takes a Handler and returns a Handler.
type Middleware func(http.Handler) http.Handler

// MiddlewareChain acts as a list of Middleware Handler functions.
// MiddlewareChain is effectively immutable, so that once created it will
// always hold the same set of Middleware Handler functions, in the same order.
type MiddlewareChain struct {
	req        *http.Request
	resp       http.ResponseWriter
	middleware []Middleware
}

// NewMiddlewareChain creates a new chain of middleware, appends any middleware
// provided, and then returns the appneded chain.
func NewMiddlewareChain(m ...Middleware) MiddlewareChain {
	return MiddlewareChain{middleware: m}
}

// Then chains the Middleware and returns the final Handler.
//     NewMiddlewareChain(m1, m2, m3).Then(h)
// is equivalent to:
//     m1(m2(m3(h)))
// When the request comes in, it will be passed to m1, then m2, then m3.
func (mc MiddlewareChain) Then(h http.Handler) http.Handler {
	if h == nil {
		log.Fatalln("MiddlewareChain requires a Then(h Handler) or ThenFunc(f HandlerFunc)")
	}
	// iterate backwares through each middleware Handler
	final := h
	for i := len(mc.middleware) - 1; i >= 0; i-- {
		final = mc.middleware[i](final)
	}
	return final
}

// ThenFunc works identically to Then, but takes a HandlerFunc instead of a
// Handler.
//
// The following two statements are equivalent:
//     c.Then(http.HandlerFunc(fn))
//     c.ThenFunc(fn)
//
// ThenFunc provides all the guarantees of Then.
func (mc MiddlewareChain) ThenFunc(f http.HandlerFunc) http.Handler {
	// set nil as the caboose value at the end of the middleware chain/train
	if f == nil {
		return mc.Then(nil)
	}
	return mc.Then(http.HandlerFunc(f))
}

// Append extends a chain, adding the specified Middleware handler as the last
// handler in the request flow.
func (mc MiddlewareChain) Append(m ...Middleware) MiddlewareChain {
	middlewares := make([]Middleware, len(mc.middleware)+len(m))
	copy(middlewares, mc.middleware)
	copy(middlewares[len(mc.middleware):], m)
	return NewMiddlewareChain(middlewares...)
}
