package utils

import (
	"errors"
	"testing"
)

type mockRunner struct {
	raw []byte
	err error
}

func (m mockRunner) run(cmd string, args ...string) ([]byte, error) {
	return m.raw, m.err
}

func TestValidateAddress(t *testing.T) {
	expectedError := errors.New("err")
	r = mockRunner{nil, expectedError}
	if _, err := ValidateAddress(""); err != expectedError {
		t.Errorf("expected error %v but get %v", expectedError, err)
	}

	r = mockRunner{nil, nil}
	if _, err := ValidateAddress(""); err == nil {
		t.Error("expected json unmarshal error but error is nil")
	}

	r = mockRunner{[]byte(`{"isvalid":true}`), nil}
	if valid, _ := ValidateAddress(""); !valid {
		t.Error("expected address valid but it is invalid")
	}
}
