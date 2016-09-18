// Package render writes HTTP responses in different formats.
package render

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"

	"gopkg.in/yaml.v2"
)

// Binary writes a raw slice of bytes to a http.ResponseWriter.
func Binary(w http.ResponseWriter, code int, p []byte) (int, error) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(code)
	return w.Write(p)
}

// JSON writes an object to a http.ResponseWriter as JSON.
func JSON(w http.ResponseWriter, code int, prettyprint bool, v interface{}) (int, error) {
	bytes := []byte{}
	var err error

	if prettyprint {
		bytes, err = json.MarshalIndent(v, "", "  ")
	} else {
		bytes, err = json.Marshal(v)
	}
	if err != nil {
		n, err2 := w.Write([]byte(err.Error()))
		return n, errors.New(err2.Error() + " " + err.Error())
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	num, err := w.Write(bytes)
	if err != nil {
		n, err2 := w.Write([]byte(err.Error()))
		return n, errors.New(err2.Error() + " " + err.Error())
	}

	if prettyprint {
		n, err := w.Write([]byte("\n"))
		num += n
		if err != nil {
			n, err2 := w.Write([]byte(err.Error()))
			return n, errors.New(err2.Error() + " " + err.Error())
		}
	}
	return num, nil
}

// JSONP wraps a callback function around an object and writes the result to a
// http.ResponseWriter as JSONP.
func JSONP(w http.ResponseWriter, code int, prettyprint bool, callback string, v interface{}) (int, error) {
	bytes := []byte{}
	var err error

	if prettyprint {
		bytes, err = json.MarshalIndent(v, "", "  ")
	} else {
		bytes, err = json.Marshal(v)
	}
	if err != nil {
		n, _ := w.Write([]byte(err.Error()))
		return n, err
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)

	// This next line is interesting, not because of the callback, but the empty
	// comment characters ('/**/') that precede it.
	//
	// These comments are added to protect against an attack where a third party
	// site bypasses the content-type of the response.
	//
	// This was originally inspired by github's api design for callbacks.
	// https://developer.github.com/v3/#json-p-callbacks
	//
	// Besides github, facebook also uses this technique. Google does something
	// different, but achieves the same result.
	//
	// For more information check out this stackoverflow link.
	// http://stackoverflow.com/questions/8034515/facebook-graph-api-jsonp-format-what-does-the-in-first-line-signify
	n, err := w.Write([]byte("/**/" + callback + "("))
	if err != nil {
		n, _ := w.Write([]byte(err.Error()))
		return n, err
	}

	// write out the collected bytes
	num, err := w.Write(bytes)
	num += n

	// now finish writing the callback bytes
	n, err = w.Write([]byte(");"))
	num += n
	if err != nil {
		n, _ := w.Write([]byte(err.Error()))
		return n, err
	}

	if prettyprint {
		n, err := w.Write([]byte("\n"))
		num += n
		if err != nil {
			n, _ := w.Write([]byte(err.Error()))
			return n, err
		}
	}
	return num, nil
}

// XML writes an object to a http.ResponseWriter as XML.
func XML(w http.ResponseWriter, code int, prettyprint bool, v interface{}) (int, error) {
	bytes := []byte{}
	var err error

	if prettyprint {
		bytes, err = xml.MarshalIndent(v, "", "  ")
	} else {
		bytes, err = xml.Marshal(v)
	}
	if err != nil {
		n, _ := w.Write([]byte(err.Error()))
		return n, err
	}

	w.Header().Set("Content-Type", "text/xml; charset=UTF-8")
	w.WriteHeader(code)

	// stay XML compliant and wrap the output
	num, err := w.Write([]byte("<Response>"))
	if err != nil {
		n, _ := w.Write([]byte(err.Error()))
		return n, err
	}

	if prettyprint {
		n, err := w.Write([]byte("\n"))
		num += n
		if err != nil {
			n, _ := w.Write([]byte(err.Error()))
			return n, err
		}
	}

	// write out the collected bytes
	n, err := w.Write(bytes)
	num += n

	if prettyprint {
		n, err := w.Write([]byte("\n"))
		num += n
		if err != nil {
			n, _ := w.Write([]byte(err.Error()))
			return n, err
		}
	}

	// finish writing XML wrapper
	n, err = w.Write([]byte("</Response>"))
	num += n
	if err != nil {
		n, _ := w.Write([]byte(err.Error()))
		return n, err
	}

	if prettyprint {
		n, err := w.Write([]byte("\n"))
		num += n
		if err != nil {
			n, _ := w.Write([]byte(err.Error()))
			return n, err
		}
	}
	return num, nil
}

// YML writes an object to a http.ResponseWriter as YML.
func YML(w http.ResponseWriter, code int, v interface{}) (int, error) {
	bytes, err := yaml.Marshal(v)
	if err != nil {
		n, _ := w.Write([]byte(err.Error()))
		return n, err
	}

	// set Content-Type to text/x-yaml since it is kinda-sorta the standard
	// for yaml right now
	w.Header().Set("Content-Type", "text/x-yaml")
	w.WriteHeader(code)
	return w.Write(bytes)
}
