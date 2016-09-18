package resource

import (
	"bytes"
	"encoding/json"
	stdhttp "net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/target/gophersaurus/http"
	httpmock "github.com/target/gophersaurus/http/mock"
	modelmock "github.com/target/gophersaurus/model/mock"
)

func TestIndex(t *testing.T) {

	// build a slice of tests to check an gf.Extended Index
	// method executes as expected for multiple formatted responses.
	tests := []struct {
		extension   string
		contentType string
		status      int
	}{
		{"", "application/json; charset=UTF-8", http.StatusOK},
		{".json", "application/json; charset=UTF-8", http.StatusOK},
		{".xml", "text/xml; charset=UTF-8", http.StatusOK},
		{".yml", "text/x-yaml", http.StatusOK},
		{".badformat", "application/json; charset=UTF-8", http.StatusBadRequest}}

	// range through the slice of tests
	for _, test := range tests {

		// mock some models
		model := modelmock.NewModel()

		// create a resource
		resource := New(model, func(req *http.Request) (string, error) {
			return strconv.Itoa(model.ID), nil
		})

		// make an empty io.Reader
		body := strings.NewReader("")

		// make a new []httprouter.Param
		params := []httprouter.Param{}

		// execute a request to the .Index method
		url := httpmock.Ext("/model", test.extension)
		rec, err := httpmock.Request(url, "GET", params, body, resource.Index)
		if err != nil {
			t.Errorf("resource.Index failed: %s", err.Error())
		}

		// check the http status code
		if rec.Code != test.status {
			t.Errorf("resource.Index failed: expected http code %d got %d.", test.status, rec.Code)
		}

		// check a http header Content-Type exists
		if len(rec.HeaderMap["Content-Type"]) < 1 {
			t.Errorf("resource.Index failed: expected Content-Type %s but got nothing.", "application/json; charset=UTF-8")
		} else {
			// check the http header Content-Type
			if rec.HeaderMap["Content-Type"][0] != test.contentType {
				t.Errorf("resource.Index failed: expected Content-Type %s but got %s.", test.contentType, rec.HeaderMap["Content-Type"][0])
			}
		}
	}
}

func TestShow(t *testing.T) {

	// build a slice of tests to check an gf.Extended Index
	// method executes as expected for multiple formatted responses.
	tests := []struct {
		extension   string
		contentType string
		status      int
	}{{"", "application/json; charset=UTF-8", http.StatusOK},
		{".json", "application/json; charset=UTF-8", http.StatusOK},
		{".xml", "text/xml; charset=UTF-8", http.StatusOK},
		{".yml", "text/x-yaml", http.StatusOK},
		{".badformat", "application/json; charset=UTF-8", http.StatusBadRequest}}

	// range through the slice of tests
	for _, test := range tests {

		// mock a model
		model := modelmock.NewModel()

		// create a resource
		resource := New(model, func(req *http.Request) (string, error) {
			return strconv.Itoa(model.ID), nil
		})

		// make an empty io.Reader
		body := strings.NewReader("")

		modelID := strconv.Itoa(model.ID)

		// make a new []httprouter.Param
		params := []httprouter.Param{{
			Key:   "modelID",
			Value: httpmock.Ext(modelID, test.extension),
		}}

		// execute a request to the .Show method
		url := "/model/:modelID"
		rec, err := httpmock.Request(url, "GET", params, body, resource.Show)

		// check error
		if err != nil {
			t.Errorf("resource.Show failed: %s", err.Error())
		}

		// check the http status code
		if rec.Code != test.status {

			// error since http status codes must match
			t.Errorf("resource.Show failed: expected http status code %d got %d.", test.status, rec.Code)
		}

		// check a http header Content-Type exists
		if len(rec.HeaderMap["Content-Type"]) < 1 {
			t.Errorf("resource.Show failed: expected Content-Type %s but got nothing.", "application/json; charset=UTF-8")
		} else {
			// check the http header Content-Type
			if rec.HeaderMap["Content-Type"][0] != test.contentType {
				t.Errorf("resource.Show failed: expected Content-Type %s but got %s.", test.contentType, rec.HeaderMap["Content-Type"][0])
			}
		}
	}
}

func TestStore(t *testing.T) {

	// build a slice of tests to check an gf. Index
	// method executes as expected for multiple formatted responses.
	tests := []struct {
		extension   string
		contentType string
		status      int
	}{{"", "application/json; charset=UTF-8", http.StatusCreated}, // json indented with list
		{".json", "application/json; charset=UTF-8", http.StatusCreated},         // json indented without list
		{".xml", "text/xml; charset=UTF-8", http.StatusCreated},                  // xml not indented without list
		{".yml", "text/x-yaml", http.StatusCreated},                              // yml not indented without list
		{".badformat", "application/json; charset=UTF-8", http.StatusBadRequest}} // bad request

	// range through the slice of tests
	for _, test := range tests {

		// mock some models
		model := modelmock.NewModel()

		// create a resource
		resource := New(model, func(req *http.Request) (string, error) {
			return strconv.Itoa(model.ID), nil
		})

		// create a new model
		m := modelmock.NewModel()

		// marshal model
		b, err := json.Marshal(m)
		if err != nil {
			t.Error("resource.Store failed: unable to marshal model.")
		}

		// make a body byte.buffer
		body := bytes.NewBuffer(b)

		// make a new []httprouter.Param
		params := []httprouter.Param{}

		// execute a request to the Store.Show method
		url := httpmock.Ext("/model", test.extension)
		rec, err := httpmock.Request(url, "POST", params, body, resource.Store)

		// check error
		if err != nil {
			t.Errorf("Store.Store failed: %s", err.Error())
		}

		// check the http status code
		if rec.Code != test.status {

			// error since http status codes must match
			t.Errorf("resource.Store failed: expected http status code %d but got %d.", test.status, rec.Code)
		}

		// check a http header Content-Type exists
		if len(rec.HeaderMap["Content-Type"]) < 1 {
			t.Errorf("resource.Store failed: expected Content-Type %s but got nothing.", "application/json; charset=UTF-8")
		} else {
			// check the http header Content-Type
			if rec.HeaderMap["Content-Type"][0] != test.contentType {
				t.Errorf("resource.Store failed: expected Content-Type %s but got %s.", test.contentType, rec.HeaderMap["Content-Type"][0])
			}
		}
	}
}

func TestApply(t *testing.T) {

	// build a slice of tests to check an gf.Extended Index
	// method executes as expected for multiple formatted responses.
	tests := []struct {
		extension   string
		contentType string
		status      int
	}{{"", "application/json; charset=UTF-8", http.StatusOK},
		{".json", "application/json; charset=UTF-8", http.StatusOK},
		{".xml", "text/xml; charset=UTF-8", http.StatusOK},
		{".yml", "text/x-yaml", http.StatusOK},
		{".badformat", "application/json; charset=UTF-8", http.StatusBadRequest}}

	// range through the slice of tests
	for _, test := range tests {

		// mock some models
		model := modelmock.NewModel()

		// create a resource
		resource := New(model, func(req *http.Request) (string, error) {
			return strconv.Itoa(model.ID), nil
		})

		// create a new model
		m := modelmock.NewModel()

		// marshal model
		b, err := json.Marshal(m)
		if err != nil {
			t.Error("Extended.Apply failed: unable to marshal model.")
		}

		// make a body byte.buffer
		body := bytes.NewBuffer(b)

		modelID := strconv.Itoa(model.ID)

		// make a new []httprouter.Param
		params := []httprouter.Param{{
			Key:   "modelID",
			Value: httpmock.Ext(modelID, test.extension),
		}}

		// execute a request to the .Show method
		url := "/model/:modelID"
		rec, err := httpmock.Request(url, "PATCH", params, body, resource.Apply)

		// check error
		if err != nil {
			t.Errorf("resource.Apply failed: %s", err.Error())
		}

		// check the http status code
		if rec.Code != test.status {

			// error since http status codes must match
			t.Errorf("resource.Apply failed: expected http status code %d got %d.", test.status, rec.Code)
		}

		// check a http header Content-Type exists
		if len(rec.HeaderMap["Content-Type"]) < 1 {
			t.Errorf("resource.Apply failed: expected Content-Type %s but got nothing.", "application/json; charset=UTF-8")
		} else {
			// check the http header Content-Type
			if rec.HeaderMap["Content-Type"][0] != test.contentType {
				t.Errorf("resource.Apply failed: expected Content-Type %s but got %s.", test.contentType, rec.HeaderMap["Content-Type"][0])
			}
		}
	}
}

func TestUpdate(t *testing.T) {

	// build a slice of tests to check an gf.Extended Index
	// method executes as expected for multiple formatted responses.
	tests := []struct {
		extension   string
		contentType string
		status      int
	}{{"", "application/json; charset=UTF-8", http.StatusOK},
		{".json", "application/json; charset=UTF-8", http.StatusOK},
		{".xml", "text/xml; charset=UTF-8", http.StatusOK},
		{".yml", "text/x-yaml", http.StatusOK},
		{".badformat", "application/json; charset=UTF-8", http.StatusBadRequest}}

	// range through the slice of tests
	for _, test := range tests {

		// mock some models
		model := modelmock.NewModel()

		// create a resource
		resource := New(model, func(req *http.Request) (string, error) {
			return strconv.Itoa(model.ID), nil
		})

		// create a new model
		m := modelmock.NewModel()

		// marshal model
		b, err := json.Marshal(m)
		if err != nil {
			t.Error("resource.Update failed: unable to marshal model.")
		}

		// make a body byte.buffer
		body := bytes.NewBuffer(b)

		// make a new []httprouter.Param
		params := []httprouter.Param{}

		// execute a request to the .Show method
		url := httpmock.Ext("/model", test.extension)
		rec, err := httpmock.Request(url, "PUT", params, body, resource.Update)

		// check error
		if err != nil {
			t.Errorf("resource.Update failed: %s", err.Error())
		}

		// check the http status code
		if rec.Code != test.status {

			// error since http status codes must match
			t.Errorf("resource.Update failed: expected http status code %d got %d.", test.status, rec.Code)
		}

		// check a http header Content-Type exists
		if len(rec.HeaderMap["Content-Type"]) < 1 {
			t.Errorf("resource.Update failed: expected Content-Type %s but got nothing.", "application/json; charset=UTF-8")
		} else {
			// check the http header Content-Type
			if rec.HeaderMap["Content-Type"][0] != test.contentType {
				t.Errorf("resource.Update failed: expected Content-Type %s but got %s.", test.contentType, rec.HeaderMap["Content-Type"][0])
			}
		}
	}
}

func TestDestroy(t *testing.T) {

	// build a slice of tests to check an gf. Index
	// method executes as expected for multiple formatted responses.
	tests := []struct {
		extension   string
		contentType string
		status      int
	}{{"", "application/json; charset=UTF-8", http.StatusNoContent},
		{".json", "application/json; charset=UTF-8", http.StatusNoContent},
		{".xml", "text/xml; charset=UTF-8", http.StatusNoContent},
		{".yml", "text/x-yaml", http.StatusNoContent},
		{".badformat", "application/json; charset=UTF-8", http.StatusBadRequest}}

	// range through the slice of tests
	for _, test := range tests {

		// mock some models
		model := modelmock.NewModel()

		// create a resource
		resource := New(model, func(req *http.Request) (string, error) {
			return strconv.Itoa(model.ID), nil
		})

		modelID := strconv.Itoa(model.ID)

		// make a new []httprouter.Param
		params := []httprouter.Param{{
			Key:   "modelID",
			Value: httpmock.Ext(modelID, test.extension),
		}}

		// make an empty io.Reader
		body := strings.NewReader("")

		// execute a request to the .Destroy method
		url := "/model/:modelID"
		rec, err := httpmock.Request(url, "GET", params, body, resource.Destroy)

		// check error
		if err != nil {
			t.Errorf("resource.Destroy failed: %s", err.Error())
		}

		// check the http status code
		if rec.Code != test.status {

			// error since http status codes must match
			t.Errorf("resource.Destroy failed: expected http status code %d got %d.", test.status, rec.Code)
		}

		// check a http header Content-Type exists
		if len(rec.HeaderMap["Content-Type"]) < 1 {
			t.Errorf("resource.Destroy failed: expected Content-Type %s but got nothing.", "application/json; charset=UTF-8")
		} else {
			// check the http header Content-Type
			if rec.HeaderMap["Content-Type"][0] != test.contentType {
				t.Errorf("resource.Destroy failed: expected Content-Type %s but got %s.", test.contentType, rec.HeaderMap["Content-Type"][0])
			}
		}
	}
}

func TestPathID(t *testing.T) {

	body := strings.NewReader("")

	r, err := stdhttp.NewRequest("get", "foo.com/some/endpoint/:id", body)
	if err != nil {
		t.Error(err)
	}

	ps := []httprouter.Param{{Key: "id", Value: "5225"}}
	req := http.NewRequest(r, ps)

	val, err := PathID(req, "id")
	if err != nil {
		t.Error(err)
	}

	if val != "5225" {
		t.Errorf("expected '5225', got '%s'", val)
	}

	val2, err2 := PathID(req, "fooID")
	if err != nil {
		t.Error(err2)
	}

	if val2 != "" {
		t.Errorf("expected '', got '%s'", val2)
	}

}
