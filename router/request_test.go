package http

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestNewRequest(t *testing.T) {

	body := strings.NewReader("")
	r, err := http.NewRequest("GET", "foo.com/some/endpoint", body)
	if err != nil {
		t.Error(err)
	}
	p := []httprouter.Param{}
	req := NewRequest(r, p)

	if req == nil {
		t.Error("expected req not to be nil")
	}
}

func TestRequestParam(t *testing.T) {

	tests := []struct {
		key    string
		val    string
		result string
	}{{"fooID", "fooValue", "fooValue"},
		{"fooID", "/fooValue", "fooValue"},
		{"jsonID", "fooValue.json", "fooValue"},
		{"xmlID", "fooValue.xml", "fooValue"},
		{"ymlID", "fooValue.yml", "fooValue"}}

	// range through the slice of tests
	for _, test := range tests {

		body := strings.NewReader("")
		r, err := http.NewRequest("GET", "foo.com/some/endpoint", body)
		if err != nil {
			t.Error(err)
		}
		p := []httprouter.Param{{
			Key:   test.key,
			Value: test.val,
		}}
		req := NewRequest(r, p)
		if req.Param(test.key) != test.result {
			t.Errorf("expected %s got %s", test.result, req.Param(test.key))
		}
	}
}

func TestRequestParamKey(t *testing.T) {

	ts := struct {
		key    string
		val    string
		result string
	}{"fooID", "fooValue", "fooValue"}

	body := strings.NewReader("")
	r, err := http.NewRequest("GET", "foo.com/some/endpoint", body)
	if err != nil {
		t.Error(err)
	}
	p := []httprouter.Param{{
		Key:   ts.key,
		Value: ts.val,
	}}
	req := NewRequest(r, p)
	if req.ParamKey(0) != ts.key {
		t.Errorf("expected %s got %s", ts.key, req.ParamKey(0))
	}
}

func TestRequestParamValue(t *testing.T) {

	tests := []struct {
		key    string
		val    string
		result string
	}{{"fooID", "fooValue", "fooValue"},
		{"fooID", "/fooValue", "fooValue"},
		{"jsonID", "fooValue.json", "fooValue"},
		{"xmlID", "fooValue.xml", "fooValue"},
		{"ymlID", "fooValue.yml", "fooValue"}}

	// range through the slice of tests
	for _, test := range tests {

		body := strings.NewReader("")
		r, err := http.NewRequest("GET", "foo.com/some/endpoint", body)
		if err != nil {
			t.Error(err)
		}
		p := []httprouter.Param{{
			Key:   test.key,
			Value: test.val,
		}}
		req := NewRequest(r, p)
		if req.ParamValue(0) != test.result {
			t.Errorf("expected %s got %s", test.result, req.ParamValue(0))
		}
	}
}

func TestRequestQuery(t *testing.T) {

	tests := []struct {
		query        string
		queryString1 string
		queryString2 string
		result       string
		found        bool
	}{{"fooQuery", "fooQuery=foo", "fooQuery2=foo2", "foo", true},
		{"fooQuery", "fooQuery=fooValue1", "fooQuery=fooValue2", "slice fooValue1 fooValue2", true},
		{"fooQueryfoo", "fooQuery=foo", "fooQuery2=foo2", "", false}}

	// range through the slice of tests
	for _, test := range tests {

		url := fmt.Sprintf("foo.com/some/endpoint?%s&%s", test.queryString1, test.queryString2)
		body := strings.NewReader("")
		r, err := http.NewRequest("GET", url, body)
		if err != nil {
			t.Error(err)
		}
		p := []httprouter.Param{}
		req := NewRequest(r, p)

		q, ok := req.Query(test.query)

		if !ok && len(q) > 0 {
			t.Errorf("ok is false, but still got query values: %s", q)
		}

		if !ok && !test.found {
			continue
		}

		if !ok && test.found {
			if len(q) < 1 {
				t.Errorf("expected query values for query '%s' got none", test.query)
			}
		}

		if ok && len(q) < 0 {
			t.Error("ok is true, but no query values found")
		}

		if len(test.result) > len("slice") && test.result[:len("slice")] == "slice" {
			if len(q) < 1 {
				t.Errorf("expected more than 1 query value for '%s' got %d", test.query, len(q))
				continue
			}
			strs := strings.Split(test.result, " ")

			if q[0] != strs[1] {
				t.Errorf("expected %s got %s", strs[1], q[0])
			}
			if q[1] != strs[2] {
				t.Errorf("expected %s got %s", strs[2], q[1])
			}
		} else {
			if q[0] != test.result {
				t.Errorf("expected %s got %s", test.result, q[0])
			}
		}
	}
}

func TestRequestQueryBool(t *testing.T) {

	tests := []struct {
		query      string
		queryValue string
		result     bool
	}{{"foo", "foo=y", true},
		{"foo", "foo=yes", true},
		{"foo", "foo=true", true},
		{"foo", "foo=✓", true},
		{"foo", "foo=no", false},
		{"foo", "foo=bad", false},
		{"foo", "foo=false", false},
		{"foo", "foo=", false},
		{"foo", "foo=asdf", false}}

	// range through the slice of tests
	for _, test := range tests {

		url := fmt.Sprintf("foo.com/some/endpoint?%s", test.queryValue)
		body := strings.NewReader("")
		r, err := http.NewRequest("GET", url, body)
		if err != nil {
			t.Error(err)
		}
		p := []httprouter.Param{}
		req := NewRequest(r, p)

		if ok := req.QueryBool(test.query); ok != test.result {
			t.Errorf("expected %t got %t", test.result, ok)
		}
	}
}

func TestRequestBytes(t *testing.T) {

	tests := []struct {
		body string
	}{{""},
		{"abcd"}}

	// range through the slice of tests
	for _, test := range tests {

		url := "foo.com/some/endpoint"
		body := strings.NewReader(test.body)
		r, err := http.NewRequest("GET", url, body)
		if err != nil {
			t.Error(err)
		}
		p := []httprouter.Param{}
		req := NewRequest(r, p)

		b, err := req.Bytes()
		if err != nil {
			t.Error(err)
		}
		if string(b) != test.body {
			t.Errorf("expected %s got %s", test.body, string(b))
		}
	}
}

func TestRequestUnmarshalJSON(t *testing.T) {

	url := "foo.com/some/endpoint"
	body := strings.NewReader(`{"title": "hello world"}`)
	r, err := http.NewRequest("GET", url, body)
	if err != nil {
		t.Error(err)
	}
	p := []httprouter.Param{}
	req := NewRequest(r, p)

	var hw helloworld

	if err := req.Unmarshal(&hw); err != nil {
		t.Error(err)
	}

	if hw.Title != "hello world" {
		t.Errorf("expected 'hello world' got '%s'", hw.Title)
	}
}

type helloworld struct {
	Title string `json:"title"`
}
