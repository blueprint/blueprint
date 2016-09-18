package mock

import (
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/target/gophersaurus/http"
)

type data struct {
	Message string `json:"message" xml:"message"`
}

func TestRequest(t *testing.T) {
	// build a slice of tests to check an gf.Extended Index
	// method executes as expected for multiple formatted responses.
	tests := []struct {
		url         string
		params      []httprouter.Param
		contentType string
		status      int
	}{{"foo.com/some/endpoint", []httprouter.Param{{Key: "callback", Value: "foo"}}, "application/json; charset=UTF-8", http.StatusOK},
		{"foo.com/some/endpoint", []httprouter.Param{}, "application/json; charset=UTF-8", http.StatusOK},
		{"foo.com/some/endpoint.json", []httprouter.Param{}, "application/json; charset=UTF-8", http.StatusOK},
		{"foo.com/some/endpoint.xml", []httprouter.Param{}, "text/xml; charset=UTF-8", http.StatusOK},
		{"foo.com/some/endpoint.yml", []httprouter.Param{}, "text/x-yaml", http.StatusOK},
		{"foo.com/some/endpoint.badformat", []httprouter.Param{}, "application/json; charset=UTF-8", http.StatusBadRequest},
		{"foo.com/some/endpoint/:id", []httprouter.Param{{Key: "id", Value: "5225"}}, "application/json; charset=UTF-8", http.StatusOK},
		{"foo.com/some/endpoint/:id", []httprouter.Param{{Key: "id", Value: "5225.xml"}}, "text/xml; charset=UTF-8", http.StatusOK},
		{"foo.com/some/endpoint/:id", []httprouter.Param{{Key: "id", Value: "5225.yml"}}, "text/x-yaml", http.StatusOK},
		{"foo.com/some/endpoint/:id", []httprouter.Param{{Key: "id", Value: "5225.badformat"}}, "application/json; charset=UTF-8", http.StatusBadRequest},
		{"foo.com/some/endpoint/:id/model", []httprouter.Param{{Key: "id", Value: "5225"}}, "application/json; charset=UTF-8", http.StatusOK},
		{"foo.com/some/endpoint/:id/model.json", []httprouter.Param{{Key: "id", Value: "5225"}}, "application/json; charset=UTF-8", http.StatusOK},
		{"foo.com/some/endpoint/:id/model.json", []httprouter.Param{{Key: "id", Value: "5225"}, {Key: "callback", Value: "foo"}}, "application/json; charset=UTF-8", http.StatusOK},
		{"foo.com/some/endpoint/:id/model.xml", []httprouter.Param{{Key: "id", Value: "5225.xml"}}, "text/xml; charset=UTF-8", http.StatusOK},
		{"foo.com/some/endpoint/:id/model.yml", []httprouter.Param{{Key: "id", Value: "5225.yml"}}, "text/x-yaml", http.StatusOK},
		{"foo.com/some/endpoint/:id/model.badformat", []httprouter.Param{{Key: "id", Value: "5225.badformat"}}, "application/json; charset=UTF-8", http.StatusBadRequest}}

	fn := func(resp http.ResponseWriter, req *http.Request) {
		data := &data{"hello world"}
		resp.WriteFormatList(req, data)
	}

	for _, test := range tests {

		body := strings.NewReader("")
		rec, err := Request(test.url, "GET", test.params, body, fn)
		if err != nil {
			t.Error(err)
		}

		if rec.Code != test.status {
			t.Errorf("expected http status code %d got %d", test.status, rec.Code)
		}

		if len(rec.HeaderMap["Content-Type"]) < 1 {
			t.Errorf("expected header Content-Type %s but Content-Type is empty, url: %s", test.contentType, test.url)
		} else if rec.HeaderMap["Content-Type"][0] != test.contentType {
			t.Errorf("expected header Content-Type %s but got %s, url: %s", test.contentType, rec.HeaderMap["Content-Type"][0], test.url)
		}
	}
}

func TestExt(t *testing.T) {
	str1 := Ext("foo.com/model", ".xml")
	str2 := Ext("foo.com/model", "")
	if str1 != "foo.com/model.xml" {
		t.Errorf("expected 'foo.com/model.xml' got '%s'", str1)
	}
	if str2 != "foo.com/model" {
		t.Errorf("expected 'foo.com/model' got '%s'", str2)
	}
}
