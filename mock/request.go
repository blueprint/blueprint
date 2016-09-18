// Package mock provides methods to mock http requests.
package mock

import (
	"io"
	stdhttp "net/http"
	"net/http/httptest"

	"github.com/julienschmidt/httprouter"
	"github.com/target/gophersaurus/http"
)

const domain = "foo.com"

// Request takes a uri string, a HTTP method, the format to respond in with
// an indentation option, an io.Reader for the reqeust body,
// a slice of []httprouter.Param for url parameters, and finally a
// gf.HandlerFunc to execute.
//
// TestRequest will then execute this HandlerFunc using the given values and
// return the resulting httptest.ResponseRecorder.
func Request(uri, method string, ps []httprouter.Param, body io.Reader, action http.HandlerFunc) (*httptest.ResponseRecorder, error) {
	var format string
	if endsInParam(uri) {
		if len(ps) > 0 {
			v := ps[len(ps)-1].Value
			format = ext(v)
		}
	} else {
		format = ext(uri)
	}

	if format != "json" && format != "xml" && format != "yml" && format != "" {
		rec := httptest.NewRecorder()
		rec.Code = http.StatusBadRequest
		rec.HeaderMap["Content-Type"] = []string{"application/json; charset=UTF-8"}
		return rec, nil
	}

	rec := httptest.NewRecorder()
	resp := http.NewResponse(rec, format)

	r, err := stdhttp.NewRequest(method, renderURL(domain, uri, ps), body)
	if err != nil {
		return nil, err
	}

	action.ServeHTTP(resp, http.NewRequest(r, ps))
	return rec, nil
}

// endsInParam checks if a URL path ends in a parameter.
func endsInParam(path string) bool {
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

func ext(path string) string {
	for i := len(path) - 2; i >= 0; i-- {
		if path[i] == '.' {
			return path[i+1:]
		}
		if path[i] == '/' {
			return ""
		}
	}
	return ""
}

func renderURL(base, uri string, ps []httprouter.Param) string {
	var name []rune
	param := false
	url := []rune(base)

	for _, r := range uri {
		if param {
			if r == '/' {
				key := string(name)
				for _, p := range ps {
					if p.Key == key {
						v := []rune(p.Value)
						url = append(url, v...)
						name = []rune{}
						param = false
						url = append(url, '/')
						break
					}
				}
				continue
			}
			name = append(name, r)
			continue
		}

		if r == ':' {
			param = true
			continue
		}

		url = append(url, r)
	}

	if len(name) > 0 {
		key := string(name)
		for _, p := range ps {
			if p.Key == key {
				v := []rune(p.Value)
				url = append(url, v...)
			}
		}
	}
	return string(url)
}

// Ext appends a file like extension to a uri.
func Ext(uri, ext string) string {
	if ext == "" {
		return uri
	}
	return uri + ext
}
