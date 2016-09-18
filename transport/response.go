package http

import (
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/target/gophersaurus/render"
)

// ResponseWriter is an interface representing a HTTP response.
type ResponseWriter interface {
	Status(code int)
	Header() http.Header
	WriteHeader(int)

	Read(p []byte) (n int, err error)
	Bytes() []byte
	HasBody() bool
	FlushBody()

	View(tmpl string, data interface{})
	WriteFormat(req *Request, v interface{})
	WriteFormatList(req *Request, v interface{})

	Raw()
	Write([]byte) (int, error)
	WriteYML(v interface{})
	WriteXML(prettyprint bool, v interface{})
	WriteJSON(prettyprint bool, v interface{})
	WriteJSONP(prettyprint bool, callback string, v interface{})
	WriteErrs(req *Request, errs ...error)
}

// Response describes a HTTP response object.
type Response struct {
	code   int
	bytes  []byte
	format string
	list   bool
	http.ResponseWriter
}

// NewResponse takes a http.ResponseWriter and returns a Response.
func NewResponse(w http.ResponseWriter, format string) *Response {
	return &Response{ResponseWriter: w, code: http.StatusOK, list: false, format: format}
}

// Raw sends a HTTP response.
//
// Raw first looks for bytes its has read, then sends them.
// Raw secondly looks for objects it has been given.  If no bytes exist,
// but objects exist then it casts the objects to bytes and sends them.
func (r *Response) Raw() {
	// body values exist
	if r.HasBody() {
		r.Write(r.bytes)
		return
	}

	// no body content and no empty list needed, write headers only
	r.WriteHeader(r.code)
}

// WriteFormat formats the response by the url type extension and other factors.
// Supported type extensions are .json and .xml, and JSONP if a callback
// value is supplied.
//
// If no format type is supplied Response defaults to JSON, or JSONP if a
// callback method is provided.
func (r *Response) WriteFormat(req *Request, v interface{}) {
	switch r.format {
	case "json", "":
		if callbacks, ok := req.Query("callback"); ok {
			r.WriteJSONP(req.QueryBool("prettyprint"), callbacks[0], v)
			return
		}
		r.WriteJSON(req.QueryBool("prettyprint"), v)
	case "xml":
		r.WriteXML(req.QueryBool("prettyprint"), v)
	case "yml":
		r.WriteYML(v)
	default:
		r.Read([]byte(InvalidInput))
		r.Status(StatusBadRequest)
		r.WriteJSON(true, nil) // default error response is JSON
	}
}

// WriteFormatList takes an interface and writes it as a list.
func (r *Response) WriteFormatList(req *Request, v interface{}) {
	r.list = true
	r.WriteFormat(req, v)
}

// Status sets the status code for the HTTP response.
func (r *Response) Status(code int) { r.code = code }

// ReadBytes takes bytes and saves them in the Response to be written later.
func (r *Response) Read(p []byte) (n int, err error) {
	r.bytes = append(r.bytes, p...)
	return len(p), nil
}

// FlushBody deletes all saved state in the Response body.
func (r *Response) FlushBody() { r.bytes = nil }

// Bytes returns all bytes that have been read by Read.
func (r *Response) Bytes() []byte { return r.bytes }

// HasBody checks if any bytes or object are being stored in the Response.
func (r *Response) HasBody() bool {
	return len(r.bytes) > 0
}

// View takes a template and data and renders a view
func (r *Response) View(tmpl string, data interface{}) {
	// new template
	t := template.New(tmpl)

	// parse the template
	t, err := t.ParseFiles(tmpl)
	if err != nil {
		r.WriteJSON(true, err.Error()) // default error response is JSON
		return
	}

	// get filename
	tmplbase := filepath.Base(tmpl)
	filename := tmplbase[:]

	// execute html
	if err = t.ExecuteTemplate(r, filename, data); err != nil {
		r.WriteJSON(true, err.Error()) // default error response is JSON
		return
	}
}

// WriteXML takes a value and writes it as a HTTP XML response.
func (r *Response) WriteXML(prettyprint bool, v interface{}) {
	render.XML(r, r.code, prettyprint, v)
}

// WriteYML takes a value and writes it as a HTTP YML response.
func (r *Response) WriteYML(v interface{}) {
	render.YML(r, r.code, v)
}

// WriteJSON takes a value and writes it as a HTTP JSON response.
func (r *Response) WriteJSON(prettyprint bool, v interface{}) {
	render.JSON(r, r.code, prettyprint, v)
}

// WriteJSONP takes a value and writes it as a HTTP JSONP response.
func (r *Response) WriteJSONP(prettyprint bool, callback string, v interface{}) {
	render.JSONP(r, r.code, prettyprint, callback, v)
}

// WriteErrs takes many errors.
//
// If the code is nil, the error string will attempt to match with its code
// against the errorMap.
//
// If no match is found and the Response Status is not 200, the current
// Response Status will be used as the default.
//
// If nothing is found and the Response Status is 200, then the HTTP Response
// code will default to 500.
func (r *Response) WriteErrs(req *Request, errs ...error) {
	var body struct {
		Status string   `json:"status,omitempty"`
		Errs   []string `json:"errors,omitempty"`
	}

	// iterate over errors
	for i, err := range errs {
		// first error
		if i == 0 {
			code, ok := ErrorMap[err]
			if !ok {
				// code is not 200, use Status
				if r.code != http.StatusOK {
					code = r.code
				} else {
					code = http.StatusInternalServerError
				}
			}
			r.code = code
			body.Status = http.StatusText(code)
		}
		body.Errs = append(body.Errs, err.Error())
	}

	r.WriteFormat(req, body)
}
