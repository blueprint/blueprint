package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestHandlerServeHTTP(t *testing.T) {

	var fn HandlerFunc

	fn = func(resp ResponseWriter, req *Request) {
		if resp == nil {
			t.Error("expected resp not to be nil")
		}
		if req == nil {
			t.Error("expected req not to be nil")
		}
	}

	rec := httptest.NewRecorder()
	resp := NewResponse(rec, "json")

	body := strings.NewReader("")
	r, err := http.NewRequest("GET", "foo.com/some/endpoint", body)
	if err != nil {
		t.Error(err)
	}

	p := []httprouter.Param{}
	req := NewRequest(r, p)

	fn.ServeHTTP(resp, req)
}
