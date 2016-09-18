// Package http manages HTTP requests, responses, errors and statuses.
package http

// Handler is an interface.  Objects implementing the Handler interface can be
// registered to serve a particular path or subtree in the gophersaurus HTTP
// server.
//
// ServeHTTP should write reply headers and data to the ResponseWriter and then
// return. Returning signals that the request is finished and that the HTTP
// server can move on to the next request on the connection.
//
// If ServeHTTP panics, the server (the caller of ServeHTTP) assumes that the
// effect of the panic was isolated to the active request. It recovers the
// panic, logs a stack trace to the server error log, and hangs up the
// connection.
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}

// The HandlerFunc type is an adapter to allow the use of ordinary functions as
// gophersaurus HTTP handlers. If f is a function with the appropriate
// signature, HandlerFunc(f) is a Handler object that calls f.
type HandlerFunc func(ResponseWriter, *Request)

// ServeHTTP satisfies the http.Hander interface.  This maintains flexiblity
// and compatiblity with the standard http package.
//
// ServeHTTP creates a new Response and Request and executes the Action method.
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {

	// execute the action
	f(w, r)

	// The last thing ServeHTTP should do is make room for another request
	// after the current request body bytes have been read.
	//
	// Go has a rule that after the http.Request close method has been executed
	// on the request body, its bytes cannot be read again.  Also because
	// http.Request body is a stream, it can only be read once.  Executing a
	// write on the http.ResponseWriter will also call close on the
	// http.Request body.
	//
	// After close is called the http.Request body cannot be read again.
	//
	// To support gf Requests that might be created for convenience (such as
	// requests in the application middleware layer), this ServeHTTP method must
	// close the http.Request body.
	//
	// Having ServeHTTP manage closing the http.Request body allows us to be
	// flexible, yet responible when building applications.
	defer r.Body.Close()
}
