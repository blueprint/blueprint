package middleware

import (
	"errors"
	"net"

	"github.com/target/gophersaurus/http"
)

// Keys describes API keys.
type Keys struct {
	success http.Handler
	keymap  map[string][]string
}

// NewKeys returns takes a map of string to string and returns a Keys object.
func NewKeys(keys map[string][]string) Keys {
	for _, whitelist := range keys {
		for i, uri := range whitelist {
			if uri == "localhost" {
				whitelist[i] = "::1"
			}
		}
	}
	return Keys{keymap: keys}
}

// Do takes a handler and executes key middleware.
func (k Keys) Do(h http.Handler) http.Handler {
	k.success = h
	return k
}

// Check takes a key and checks it.
func (k Keys) Check(key, remoteAddr string) error {
	if whitelist, ok := k.keymap[key]; ok {

		// no whitelist is provided, so any request is ok
		if len(whitelist) == 0 {
			return nil
		}

		for _, uri := range whitelist {
			// if 'all' or '*' is found, allow any remote ip address
			if uri == "all" || uri == "*" || uri == remoteAddr {
				return nil
			}
		}
	}
	return errors.New(http.InvalidPermission)
}

// ServeHTTP fulfills the http package interface for middlewares.
func (k Keys) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// get client IP address
	ip, _, _ := net.SplitHostPort(req.RemoteAddr)

	// sneaky little proxies (or loadbalancers), wicked, tricksy, false!
	if proxy := req.Header.Get("X-FORWARDED-FOR"); len(proxy) > 0 {
		ip = proxy
	}

	// parameter
	if key := req.URL.Query().Get("key"); len(key) > 0 {
		if err := k.Check(key, ip); err != nil {
			resp.WriteErrs(req, err)
			return
		}
		k.success.ServeHTTP(resp, req)
		return
	}

	// header
	if key := req.Header.Get("API-Key"); len(key) > 0 {
		if err := k.Check(key, ip); err != nil {
			resp.WriteErrs(req, err)
			return
		}
		k.success.ServeHTTP(resp, req)
		return
	}
	resp.WriteErrs(req, errors.New(http.InvalidPermission))
}
