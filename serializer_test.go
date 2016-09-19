package blueprint

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestBinarySerializer(t *testing.T) {
	v := "test"
	var b bytes.Buffer
	_, err := Binary(&b, []byte(v))
	if err != nil {
		t.Error(err)
	}
}

// Binary writes a raw slice of bytes to a io.Writer.
func Binary(w io.Writer, v interface{}) (int, error) {
	switch v.(type) {
	case []byte:
		return w.Write(v.([]byte))
	default:
		return 0, errors.New("expected a slice of bytes")
	}
}
