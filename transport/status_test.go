package http

import "testing"

func TestStatusText(t *testing.T) {
	str1 := StatusText(StatusOK)
	str2 := StatusText(9999)

	if str1 != "OK" {
		t.Errorf("expected '%s' got '%s'", "OK", str1)
	}

	if str2 != "" {
		t.Errorf("expected '%s' got '%s'", "", str2)
	}
}

func TestShameOnYou(t *testing.T) {
	if !ShameOnYou(400) {
		t.Error("expected true got false")
	}
	if ShameOnYou(500) {
		t.Error("expected false got true")
	}
}

func TestStatus400(t *testing.T) {
	if !Status400(400) {
		t.Error("expected true got false")
	}
	if Status400(500) {
		t.Error("expected false got true")
	}
}

func TestShameOnMe(t *testing.T) {
	if ShameOnMe(400) {
		t.Error("expected true got false")
	}
	if !ShameOnMe(500) {
		t.Error("expected false got true")
	}
}

func TestStatus500(t *testing.T) {
	if Status500(400) {
		t.Error("expected true got false")
	}
	if !Status500(500) {
		t.Error("expected false got true")
	}
}
