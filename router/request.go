package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// Request describes a HTTP Request.
type Request struct {
	http.Request
	params  httprouter.Params
	queries url.Values
	buf     []byte
}

// NewRequest takes a http.Request returns a new Request.
func NewRequest(r *http.Request, ps httprouter.Params) *Request {
	return &Request{Request: *r, params: ps, queries: r.URL.Query()}
}

// Param searches for a variable identifier in the URL path of the http.Request.
//
// If a match is found then the coresponding value is returned with true.
// If a match is not found then an empty string is returned with false.
func (r *Request) Param(name string) string {
	p := r.params.ByName(name)
	l := len(p)
	if l > 4 {
		if p[l-4:] == ".xml" || p[l-4:] == ".yml" {
			p = p[:l-4]
		}
		if l > 5 && p[l-5:] == ".json" {
			p = p[:l-5]
		}
	}
	if l > 1 && p[0] == '/' {
		return p[1:]
	}
	return p
}

// ParamKey searches for a variable identifier in the URL path of the http.Request.
//
// If a match is found then the coresponding value is returned with true.
// If a match is not found then an empty string is returned with false.
func (r *Request) ParamKey(index int) string {
	return r.params[index].Key
}

// ParamValue searches for a variable identifier in the URL path of the http.Request.
//
// If a match is found then the coresponding value is returned with true.
// If a match is not found then an empty string is returned with false.
func (r *Request) ParamValue(index int) string {
	p := r.params[index].Value
	l := len(p)
	if l > 4 {
		if p[l-4:] == ".xml" || p[l-4:] == ".yml" {
			p = p[:l-4]
		}
		if l > 5 && p[l-5:] == ".json" {
			p = p[:l-5]
		}
	}
	if p[0] == '/' {
		return p[1:]
	}
	return p
}

// Query searches for a query parameter in the URL path of the http.Request.
//
// If a match is found then the coresponding value is returned with true.
// If a match is not found then an empty string is returned with false.
func (r *Request) Query(name string) ([]string, bool) {
	i, ok := r.queries[name]
	return i, ok
}

// QueryBool takes a query name and returns the corresponding value as a
// boolean.
func (r *Request) QueryBool(name string) bool {
	out := false
	flagList, ok := r.Query(name)
	if ok && len(flagList) > 0 {
		flag := strings.ToLower(flagList[0])
		if flag == "y" || flag == "yes" || flag == "true" || flag == "✓" {
			out = true
		}
	}
	return out
}

// Bytes returns the bytes of http.Request body.
func (r *Request) Bytes() ([]byte, error) {
	// check if the http.Request has already been read
	if r.buf == nil {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		// save body bytes for later use
		r.buf = body
	}
	return r.buf, nil
}

// Unmarshal parses the JSON-encoded data in the Request body and stores the result in the value pointed to by v.
func (r *Request) Unmarshal(v interface{}) error {
	bytes, err := r.Bytes()
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, v); err != nil {
		return fmt.Errorf("%s: %s", InvalidJSON, err)
	}
	return nil
}
