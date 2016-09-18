package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewResponseWriter(t *testing.T) {
	rec := httptest.NewRecorder()
	w := NewResponseWriter(rec)
	if w == nil {
		t.Error("expected ResponseWriter not to be nil")
	}
}

func TestResponseWriterWriteHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	w := NewResponseWriter(rec)
	w.WriteHeader(http.StatusCreated)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected %d got %d", http.StatusCreated, rec.Code)
	}
}

func TestResponseWriterWrite(t *testing.T) {
	rec := httptest.NewRecorder()
	w := NewResponseWriter(rec)

	data := []byte("hello world")
	l := len(data)

	num, err := w.Write(data)
	if err != nil {
		t.Error(err)
	}

	if l != num {
		t.Errorf("expected %d bytes got %d bytes written", l, num)
	}

	bytes := rec.Body.Bytes()

	for i := range bytes {
		if bytes[i] != data[i] {
			t.Errorf("expected '%s' was written got '%s'", string(rec.Body.Bytes()), string(data))
		}
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, rec.Code)
	}
}

func TestResponseWriterStatus(t *testing.T) {
	rec := httptest.NewRecorder()
	w := NewResponseWriter(rec)
	w.WriteHeader(http.StatusCreated)
	if w.Status() != http.StatusCreated {
		t.Errorf("expected %d got %d", http.StatusCreated, w.Status())
	}
}

func TestResponseWriterSize(t *testing.T) {
	rec := httptest.NewRecorder()
	w := NewResponseWriter(rec)

	data := []byte("hello world")
	l := len(data)

	w.Write(data)

	if l != w.Size() {
		t.Errorf("expected %d bytes got %d bytes written", l, w.Size())
	}
}
