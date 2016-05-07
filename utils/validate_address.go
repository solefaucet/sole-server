package utils

import (
	"encoding/json"
	"os/exec"
)

type runner interface {
	run(cmd string, args ...string) ([]byte, error)
}

type realRunner struct{}

func (realRunner) run(cmd string, args ...string) ([]byte, error) {
	return exec.Command(cmd, args...).CombinedOutput()
}

var r runner = realRunner{}

// ValidateAddress validates cryptocurrency address
func ValidateAddress(address string) (bool, error) {
	raw, err := r.run("coin-cli", "validateaddress", address)
	if err != nil {
		return false, err
	}

	dest := struct {
		IsValid bool `json:"isvalid"`
	}{}
	if err := json.Unmarshal(raw, &dest); err != nil {
		return false, err
	}

	return dest.IsValid, nil
}
