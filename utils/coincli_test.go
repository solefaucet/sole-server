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

func TestGetInputAddress(t *testing.T) {
	expectedError := errors.New("err")
	r = mockRunner{nil, expectedError}
	if _, err := GetInputAddress(); err != expectedError {
		t.Errorf("expected error %v but get %v", expectedError, err)
	}

	r = mockRunner{nil, nil}
	if _, err := GetInputAddress(); err == nil {
		t.Error("expected json unmarshal error but error is nil")
	}

	r = mockRunner{[]byte(`["address"]`), nil}
	if address, _ := GetInputAddress(); address != "address" {
		t.Errorf("expected to get address but get %v", address)
	}
}

func TestGetBalance(t *testing.T) {
	expectedError := errors.New("err")
	r = mockRunner{nil, expectedError}
	if _, err := GetBalance(); err != expectedError {
		t.Errorf("expected error %v but get %v", expectedError, err)
	}

	r = mockRunner{nil, nil}
	if _, err := GetBalance(); err == nil {
		t.Error("expected error but error is nil")
	}

	r = mockRunner{[]byte(`10.086`), nil}
	if balance, _ := GetBalance(); balance != 10.086 {
		t.Errorf("expected to get balance 10.086 but get %v", balance)
	}
}

func TestSendTo(t *testing.T) {
	expectedError := errors.New("err")
	r = mockRunner{nil, expectedError}
	if _, err := SendTo("", 1); err != expectedError {
		t.Errorf("expected error %v but get %v", expectedError, err)
	}

	r = mockRunner{[]byte("txid"), nil}
	if txid, _ := SendTo("", 1); txid != "txid" {
		t.Errorf("expected transactionID txid but get %v", txid)
	}
}
