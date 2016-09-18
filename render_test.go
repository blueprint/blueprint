package render

import (
	"net/http/httptest"
	"testing"
)

func TestRenderBinary(t *testing.T) {

	// build a list of tests
	tests := []struct {
		code    int
		content string
	}{
		{200, "binary test"},
	}

	// range through the slice of tests
	for _, test := range tests {

		// make a new http test recorder
		rec := httptest.NewRecorder()

		// execute the Binary method
		n, err := Binary(rec, test.code, []byte(test.content))

		// check errors
		if err != nil {
			t.Error(err)
		}

		// check http status code
		if rec.Code != test.code {
			t.Errorf("render.Binary failed: expected http status code %d got %d.", test.code, rec.Code)
		}

		// check byte length
		if len(rec.Body.Bytes()) != n {
			t.Errorf("render.Binary failed: expected a written byte length of %d got %d.", n, len(rec.Body.Bytes()))
		}

		// check if bytes match
		if rec.Body.String() != test.content {
			t.Errorf("render.Binary failed: expected bytes '%s' got '%s'.", test.content, rec.Body)
		}
	}
}

func TestRenderJSON(t *testing.T) {

	// build a list of tests
	tests := []struct {
		code        int
		prettyprint bool
		content     string
	}{
		{200, true, "json test"},
		{200, false, "json test"},
	}

	// range through the slice of tests
	for _, test := range tests {

		// make a new http test recorder
		rec := httptest.NewRecorder()

		// execute the JSON method
		n, err := JSON(rec, test.code, test.prettyprint, test.content)

		// check errors
		if err != nil {
			t.Error(err)
		}

		want := []byte("\"" + test.content + "\"")

		if test.prettyprint {
			endline := []byte("\n")
			want = append(want, endline...)
		}

		// check http status code
		if rec.Code != test.code {
			t.Errorf("render.JSON failed: expected http status code %d got %d.", test.code, rec.Code)
		}

		// check byte length
		if len(rec.Body.Bytes()) != n {
			t.Errorf("render.JSON failed: expected a written byte length of %d got %d.", n, len(rec.Body.Bytes()))
		}

		// check if bytes match
		if rec.Body.String() != string(want) {
			t.Errorf("render.JSON failed: expected bytes '%s' got '%s'.", string(want), rec.Body)
		}
	}
}

func TestRenderJSONP(t *testing.T) {

	// build a list of tests
	tests := []struct {
		code        int
		prettyprint bool
		content     string
		callback    string
	}{
		{200, true, "json-p test", "yolo"},
		{200, false, "json-p test", "yolo"},
		{200, true, "json-p test", ""},
		{200, false, "json-p test", ""},
	}

	// range through the slice of tests
	for _, test := range tests {

		// make a new http test recorder
		rec := httptest.NewRecorder()

		// execute the JSONP method
		n, err := JSONP(rec, test.code, test.prettyprint, test.callback, test.content)

		// check errors
		if err != nil {
			t.Error(err)
		}

		want := []byte("/**/" + test.callback + "(\"" + test.content + "\");")

		if test.prettyprint {
			endline := []byte("\n")
			want = append(want, endline...)
		}

		// check http status code
		if rec.Code != test.code {
			t.Errorf("render.JSONP failed: expected http status code %d got %d.", test.code, rec.Code)
		}

		// check byte length
		if len(rec.Body.Bytes()) != n {
			t.Errorf("render.JSONP failed: expected a written byte length of %d got %d.", n, len(rec.Body.Bytes()))
		}

		// check if bytes match
		if rec.Body.String() != string(want) {
			t.Errorf("render.JSONP failed: expected bytes '%s' got '%s'.", string(want), rec.Body)
		}
	}
}

func TestRenderXML(t *testing.T) {

	var prettyresult = `<Response>
<string>xml test</string>
</Response>
`
	var result = "<Response><string>xml test</string></Response>"

	// build a list of tests
	tests := []struct {
		code        int
		prettyprint bool
		have        string
		want        string
	}{
		{200, true, "xml test", prettyresult},
		{200, false, "xml test", result},
		{200, true, "xml test", prettyresult},
		{200, false, "xml test", result},
	}

	// range through the slice of tests
	for _, test := range tests {

		// make a new http test recorder
		rec := httptest.NewRecorder()

		// execute the XML method
		n, err := XML(rec, test.code, test.prettyprint, test.have)

		// check errors
		if err != nil {
			t.Error(err)
		}

		// check http status code
		if rec.Code != test.code {
			t.Errorf("render.XML failed: expected http status code %d got %d.", test.code, rec.Code)
		}

		// check byte length
		if len(rec.Body.Bytes()) != n {
			t.Errorf("render.XML failed: expected a written byte length of %d got %d.", n, len(rec.Body.Bytes()))
		}

		// check if bytes match
		if rec.Body.String() != test.want {
			t.Errorf("render.XML failed: expected bytes '%s' got '%s'.", test.want, rec.Body)
		}
	}
}

func TestRenderYML(t *testing.T) {

	var result = `yml test
`

	// build a list of tests
	tests := []struct {
		code int
		have string
		want string
	}{
		{200, "yml test", result},
	}

	// range through the slice of tests
	for _, test := range tests {

		// make a new http test recorder
		rec := httptest.NewRecorder()

		// execute the XML method
		n, err := YML(rec, test.code, test.have)

		// check errors
		if err != nil {
			t.Error(err)
		}

		// check http status code
		if rec.Code != test.code {
			t.Errorf("render.YML failed: expected http status code %d got %d.", test.code, rec.Code)
		}

		// check byte length
		if len(rec.Body.Bytes()) != n {
			t.Errorf("render.YML failed: expected a written byte length of %d got %d.", n, len(rec.Body.Bytes()))
		}

		// check if bytes match
		if rec.Body.String() != test.want {
			t.Errorf("render.YML failed: expected bytes '%s' got '%s'.", test.want, rec.Body)
		}
	}
}
