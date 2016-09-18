package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const prettyprint = true

func TestNewResponse(t *testing.T) {
	rec := httptest.NewRecorder()
	resp := NewResponse(rec, "json")
	if resp == nil {
		t.Error("expected resp not to be nil")
	}
}

func TestResponseRaw(t *testing.T) {
	rec := httptest.NewRecorder()
	resp := NewResponse(rec, "json")
	resp.Raw()
}

func TestResponseRawRead(t *testing.T) {
	rec := httptest.NewRecorder()
	resp := NewResponse(rec, "json")
	str := "hello world"
	resp.Read([]byte(str))
	resp.Raw()
}

func TestResponseRawReadList(t *testing.T) {
	rec := httptest.NewRecorder()
	resp := NewResponse(rec, "json")
	str1 := "hello world"
	str2 := "hello world"
	str3 := "hello world"
	resp.Read([]byte(str1))
	resp.Read([]byte(str2))
	resp.Read([]byte(str3))
	resp.Raw()
}

func TestResponseRawReadBytes(t *testing.T) {
	rec := httptest.NewRecorder()
	resp := NewResponse(rec, "json")
	str := "hello world"
	resp.Read([]byte(str))
	resp.Raw()
}

func TestResponseXML(t *testing.T) {
	rec := httptest.NewRecorder()
	resp := NewResponse(rec, "xml")
	str := "hello world"
	resp.WriteXML(prettyprint, str)
}

func TestResponseXMLList(t *testing.T) {
	rec := httptest.NewRecorder()
	resp := NewResponse(rec, "xml")
	str1 := "hello world"
	str2 := "hello world"
	str3 := "hello world"
	resp.WriteXML(prettyprint, []string{str1, str2, str3})
}

func TestResponseEmptyXML(t *testing.T) {
	rec := httptest.NewRecorder()
	resp := NewResponse(rec, "xml")
	resp.WriteXML(prettyprint, nil)
}

func TestResponseYML(t *testing.T) {
	rec := httptest.NewRecorder()
	resp := NewResponse(rec, "yml")
	resp.WriteYML("hello world")
}

func TestResponseYMLList(t *testing.T) {
	rec := httptest.NewRecorder()
	resp := NewResponse(rec, "yml")
	str1 := "hello world"
	str2 := "hello world"
	str3 := "hello world"
	resp.WriteYML([]string{str1, str2, str3})
}

func TestResponseEmptyYML(t *testing.T) {
	rec := httptest.NewRecorder()
	resp := NewResponse(rec, "yml")
	resp.WriteYML(nil)
}

func TestResponseJSON(t *testing.T) {

	// make a new Response
	rec := httptest.NewRecorder()
	resp := NewResponse(rec, "json")
	resp.WriteJSON(prettyprint, "hello world")
}

func TestResponseJSONList(t *testing.T) {

	// make a new Response
	rec := httptest.NewRecorder()
	resp := NewResponse(rec, "json")
	str1 := "hello world"
	str2 := "hello world"
	str3 := "hello world"
	resp.WriteJSON(prettyprint, []string{str1, str2, str3})
}

func TestResponseEmptyJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	resp := NewResponse(rec, "json")
	resp.WriteJSON(prettyprint, nil)
}

func TestResponseJSONP(t *testing.T) {

	// make a new Response
	rec := httptest.NewRecorder()
	resp := NewResponse(rec, "json")
	resp.WriteJSONP(prettyprint, "yolo", "hello world")
}

func TestResponseJSONPList(t *testing.T) {

	// make a new Response
	rec := httptest.NewRecorder()
	resp := NewResponse(rec, "json")
	str1 := "hello world"
	str2 := "hello world"
	str3 := "hello world"
	resp.WriteJSONP(prettyprint, "yolo", []string{str1, str2, str3})
}

func TestResponseEmptyJSONP(t *testing.T) {
	rec := httptest.NewRecorder()
	resp := NewResponse(rec, "json")
	resp.WriteJSONP(prettyprint, "yolo", nil)
}

func TestResponseStatus(t *testing.T) {

	// make a new Response
	rec := httptest.NewRecorder()
	resp := NewResponse(rec, "json")

	resp.Status(http.StatusNonAuthoritativeInfo)

	resp.WriteJSON(prettyprint, nil)

	if rec.Code != http.StatusNonAuthoritativeInfo {
		t.Errorf("Response.Status failed: expected http status code %d but got %d.", http.StatusNonAuthoritativeInfo, rec.Code)
	}
}

func TestResponseHeader(t *testing.T) {

	// make a new Response
	rec := httptest.NewRecorder()
	resp := NewResponse(rec, "json")

	resp.Header().Set("Content-Type", "xyz")

	// Do response
	resp.Raw()

	if rec.HeaderMap["Content-Type"][0] != "xyz" {
		t.Errorf("Response.Header failed: expected header Content-Type %s but got %s.", "xyz", rec.HeaderMap["Content-Type"][0])
	}
}
