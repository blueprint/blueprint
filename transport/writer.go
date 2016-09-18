package router

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

// ResponseWriter is a wrapper around http.ResponseWriter that provides extra
// information about the response. It is recommended that middleware handlers
// use this construct to wrap a ResponseWriter if needed.
type ResponseWriter interface {

	// ResponseWriter should also satisfy http.ResponseWriter.
	http.ResponseWriter

	// ResponseWriter should also satisfy http.Flusher.
	http.Flusher

	// Status returns the status code of the response or 0 if the response has not been written.
	Status() int

	// Written returns whether or not the ResponseWriter has been written.
	Written() bool

	// Size returns the size of the response body.
	Size() int

	// Before allows for a function to be called before the ResponseWriter has been written to. This is
	// useful for setting headers or any other operations that must happen before a response has been written.
	Before(func(ResponseWriter))
}

type beforeFunc func(ResponseWriter)

// NewResponseWriter creates a ResponseWriter that wraps an http.ResponseWriter
func NewResponseWriter(rw http.ResponseWriter) ResponseWriter {
	return &responseWriter{rw, 0, 0, nil}
}

type responseWriter struct {
	http.ResponseWriter
	status      int
	size        int
	beforeFuncs []beforeFunc
}

// WriteHeader writes the HTTP Header along with the HTTP status code provided.
func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.callBefore()
	rw.ResponseWriter.WriteHeader(code)
}

// Write writes the bytes provided to the ResponseWriter. It also satisfies the
// io.Writer interface.
func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.Written() {
		rw.WriteHeader(http.StatusOK) // status if WriteHeader has not been called
	}
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

// Written returns true if the ResponseWriter has called the Write methods.
func (rw *responseWriter) Written() bool {
	return rw.status != 0
}

// Before sets a function to call before the ResponseWriter writes bytes.
func (rw *responseWriter) Before(before func(ResponseWriter)) {
	rw.beforeFuncs = append(rw.beforeFuncs, before)
}

// Hijack allows ResponseWriter to hijack a network request. It also satisfies
// the http.Hijacker interface.
func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the ResponseWriter doesn't support the Hijacker interface")
	}
	return hijacker.Hijack()
}

// CloseNotify allow detection when the underlying connection has gone away. It
// also satsifies the http.CloseNotifier interface.
func (rw *responseWriter) CloseNotify() <-chan bool {
	return rw.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (rw *responseWriter) callBefore() {
	for i := len(rw.beforeFuncs) - 1; i >= 0; i-- {
		rw.beforeFuncs[i](rw)
	}
}

// Flush allows an HTTP handler to flush buffered data to the client. It also
// satisfies the http.Flusher interface.
func (rw *responseWriter) Flush() {
	flusher, ok := rw.ResponseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}
