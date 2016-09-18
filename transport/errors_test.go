package http

import (
	"errors"
	"testing"
)

func TestApplyErrorCode(t *testing.T) {
	validError := errors.New("valid new error")
	invalidError := ErrExpiredSession

	if !ApplyErrorCode(validError, StatusBadRequest) {
		t.Errorf("expected validError: %s as %d to pass", validError, StatusBadRequest)
	}

	if ApplyErrorCode(invalidError, StatusBadRequest) {
		t.Errorf("expected invalidError: %s as %d to fail", validError, StatusBadRequest)
	}
}
