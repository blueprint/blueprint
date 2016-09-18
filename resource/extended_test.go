package resource

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/target/gophersaurus/http"
	httpmock "github.com/target/gophersaurus/http/mock"
	modelmock "github.com/target/gophersaurus/model/mock"
)

func TestExtendedIndex(t *testing.T) {

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
		base := modelmock.NewBaseModel()

		// create a resource
		resource := New(base, func(req *http.Request) (string, error) {
			return strconv.Itoa(base.ID), nil
		})

		// create an extended resource
		extended := resource.Extend(model)

		baseID := strconv.Itoa(base.ID)

		params := []httprouter.Param{{
			Key:   "baseID",
			Value: baseID,
		}}

		url := httpmock.Ext("/base/:baseID/model", test.extension)
		body := strings.NewReader("")
		rec, err := httpmock.Request(url, "GET", params, body, extended.Index)

		// check error
		if err != nil {
			t.Errorf("Extended.Index failed: %s", err.Error())
		}

		// check http status code
		if rec.Code != test.status {
			t.Errorf("Extended.Index failed: expected http status code %d got %d.", test.status, rec.Code)
		}

		// check http header Content-Type
		if len(rec.HeaderMap["Content-Type"]) < 1 {
			t.Errorf("Extended.Index failed: expected header Content-Type %s but got nothing.", "application/json; charset=UTF-8")
		} else if rec.HeaderMap["Content-Type"][0] != test.contentType {
			t.Errorf("Extended.Index failed: expected header Content-Type %s but got %s.", test.contentType, rec.HeaderMap["Content-Type"][0])
		}
	}
}

func TestExtendedShow(t *testing.T) {

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
		base := modelmock.NewBaseModel()

		// create a resource
		resource := New(base, func(req *http.Request) (string, error) {
			return strconv.Itoa(base.ID), nil
		})

		// create an extended resource
		extended := resource.Extend(model, func(req *http.Request) (string, error) {
			return strconv.Itoa(model.ID), nil
		})

		baseID := strconv.Itoa(base.ID)
		modelID := strconv.Itoa(model.ID)

		params := []httprouter.Param{{
			Key:   "baseID",
			Value: baseID,
		}, {
			Key:   "modelID",
			Value: httpmock.Ext(modelID, test.extension),
		}}

		url := "foo.com/base/:baseID/model/:modelID"
		body := strings.NewReader("")
		rec, err := httpmock.Request(url, "GET", params, body, extended.Show)

		// check error
		if err != nil {
			t.Errorf("Extended.Show failed: %s", err.Error())
		}

		// check http status code
		if rec.Code != test.status {
			t.Errorf("Extended.Show failed: expected http status code %d got %d.", test.status, rec.Code)
		}

		// check http header Content-Type
		if len(rec.HeaderMap["Content-Type"]) < 1 {
			t.Errorf("Extended.Index failed: expected header Content-Type %s but got nothing.", "application/json; charset=UTF-8")
		} else if rec.HeaderMap["Content-Type"][0] != test.contentType {
			t.Errorf("Extended.Index failed: expected header Content-Type %s but got %s.", test.contentType, rec.HeaderMap["Content-Type"][0])
		}
	}
}

func TestExtendedStore(t *testing.T) {

	// build a slice of tests to check an gf.Extended Index
	// method executes as expected for multiple formatted responses.
	tests := []struct {
		extension   string
		contentType string
		status      int
	}{{"", "application/json; charset=UTF-8", http.StatusCreated},
		{".json", "application/json; charset=UTF-8", http.StatusCreated},
		{".xml", "text/xml; charset=UTF-8", http.StatusCreated},
		{".yml", "text/x-yaml", http.StatusCreated},
		{".badformat", "application/json; charset=UTF-8", http.StatusBadRequest}}

	// range through the slice of tests
	for _, test := range tests {

		// mock some models
		model := modelmock.NewModel()
		base := modelmock.NewBaseModel()

		// create a resource
		resource := New(base, func(req *http.Request) (string, error) {
			return strconv.Itoa(base.ID), nil
		})

		// create an extended resource
		extended := resource.Extend(model, func(req *http.Request) (string, error) {
			return strconv.Itoa(model.ID), nil
		})

		// create a new model
		m := modelmock.NewModel()

		// marshal model
		b, err := json.Marshal(m)
		if err != nil {
			t.Error("Extended.Store failed: unable to marshal model.")
		}

		// make a body byte.buffer
		body := bytes.NewBuffer(b)

		baseID := strconv.Itoa(base.ID)

		params := []httprouter.Param{{
			Key:   "baseID",
			Value: baseID,
		}}

		url := httpmock.Ext("foo.com/base/:baseID/model", test.extension)
		rec, err := httpmock.Request(url, "POST", params, body, extended.Store)

		// check error
		if err != nil {
			t.Errorf("Extended.Store failed: %s", err.Error())
		}

		// check http status code
		if rec.Code != test.status {
			t.Errorf("Extended.Store failed: expected http status code %d but got %d.", test.status, rec.Code)
		}

		// check http header Content-Type
		if len(rec.HeaderMap["Content-Type"]) < 1 {
			t.Errorf("Extended.Index failed: expected header Content-Type %s but got nothing.", "application/json; charset=UTF-8")
		} else if rec.HeaderMap["Content-Type"][0] != test.contentType {
			t.Errorf("Extended.Index failed: expected header Content-Type %s but got %s.", test.contentType, rec.HeaderMap["Content-Type"][0])
		}
	}
}

func TestExtendedApply(t *testing.T) {

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
		base := modelmock.NewBaseModel()

		// create a resource
		resource := New(base, func(req *http.Request) (string, error) {
			return strconv.Itoa(base.ID), nil
		})

		// create an extended resource
		extended := resource.Extend(model, func(req *http.Request) (string, error) {
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

		baseID := strconv.Itoa(base.ID)

		params := []httprouter.Param{{
			Key:   "baseID",
			Value: baseID,
		}}

		url := httpmock.Ext("foo.com/base/:baseID/model", test.extension)
		rec, err := httpmock.Request(url, "PATCH", params, body, extended.Apply)

		// check error
		if err != nil {
			t.Errorf("Extended.Apply failed: %s", err.Error())
		}

		// check http status code
		if rec.Code != test.status {
			t.Errorf("Extended.Apply failed: expected http status code %d got %d.", test.status, rec.Code)
		}

		// check http header Content-Type
		if len(rec.HeaderMap["Content-Type"]) < 1 {
			t.Errorf("Extended.Index failed: expected header Content-Type %s but got nothing.", "application/json; charset=UTF-8")
		} else if rec.HeaderMap["Content-Type"][0] != test.contentType {
			t.Errorf("Extended.Index failed: expected header Content-Type %s but got %s.", test.contentType, rec.HeaderMap["Content-Type"][0])
		}
	}
}

func TestExtendedUpdate(t *testing.T) {

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
		base := modelmock.NewBaseModel()

		// create a resource
		resource := New(base, func(req *http.Request) (string, error) {
			return strconv.Itoa(base.ID), nil
		})

		// create an extended resource
		extended := resource.Extend(model, func(req *http.Request) (string, error) {
			return strconv.Itoa(model.ID), nil
		})

		// create a new model
		m := modelmock.NewModel()

		// marshal model
		b, err := json.Marshal(m)
		if err != nil {
			t.Error("Extended.Update failed: unable to marshal model.")
		}

		// make a body byte.buffer
		body := bytes.NewBuffer(b)

		baseID := strconv.Itoa(base.ID)

		params := []httprouter.Param{{
			Key:   "baseID",
			Value: baseID,
		}}

		baseURL := httpmock.Ext("foo.com/base/:baseID/model", test.extension)
		rec, err := httpmock.Request(baseURL, "PUT", params, body, extended.Update)

		// check error
		if err != nil {
			t.Errorf("Extended.Update failed: %s", err.Error())
		}

		// check http status code
		if rec.Code != test.status {
			t.Errorf("Extended.Update failed: expected http status code %d got %d.", test.status, rec.Code)
		}

		// check http header Content-Type
		if len(rec.HeaderMap["Content-Type"]) < 1 {
			t.Errorf("Extended.Index failed: expected header Content-Type %s but got nothing.", "application/json; charset=UTF-8")
		} else if rec.HeaderMap["Content-Type"][0] != test.contentType {
			t.Errorf("Extended.Index failed: expected header Content-Type %s but got %s.", test.contentType, rec.HeaderMap["Content-Type"][0])
		}
	}
}

func TestExtendedDestroy(t *testing.T) {

	// build a slice of tests to check an gf.Extended Index
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
		base := modelmock.NewBaseModel()

		// create a resource
		resource := New(base, func(req *http.Request) (string, error) {
			return strconv.Itoa(base.ID), nil
		})

		// create an extended resource
		extended := resource.Extend(model, func(req *http.Request) (string, error) {
			return strconv.Itoa(model.ID), nil
		})

		// make an empty io.Reader
		body := strings.NewReader("")

		baseID := strconv.Itoa(base.ID)
		modelID := strconv.Itoa(model.ID)

		params := []httprouter.Param{{
			Key:   "baseID",
			Value: baseID,
		}, {
			Key:   "modelID",
			Value: httpmock.Ext(modelID, test.extension),
		}}

		url := "foo.com/base/:baseID/model/:modelID"
		rec, err := httpmock.Request(url, "GET", params, body, extended.Destroy)

		// check error
		if err != nil {
			t.Errorf("Extended.Destroy failed: %s", err.Error())
		}

		// check http status code
		if rec.Code != test.status {
			t.Errorf("Extended.Destroy failed: expected http status code %d got %d.", test.status, rec.Code)
		}

		// check http header Content-Type
		if len(rec.HeaderMap["Content-Type"]) < 1 {
			t.Errorf("Extended.Destroy failed: expected header Content-Type %s but got nothing.", "application/json; charset=UTF-8")
		} else if rec.HeaderMap["Content-Type"][0] != test.contentType {
			t.Errorf("Extended.Destroy failed: expected header Content-Type %s but got %s.", test.contentType, rec.HeaderMap["Content-Type"][0])
		}
	}
}
