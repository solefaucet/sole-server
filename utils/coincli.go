package utils

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
)

type runner interface {
	run(cmd string, args ...string) ([]byte, error)
}

type realRunner struct{}

func (realRunner) run(cmd string, args ...string) ([]byte, error) {
	raw, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"cmd":   strings.Join(append([]string{cmd}, args...), " "),
		}).Debug("failed to execute command")
	}
	return raw, err
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
	if err := unmarshalJSON(raw, &dest); err != nil {
		return false, err
	}

	return dest.IsValid, nil
}

// GetInputAddress returns the address used for sending coin
func GetInputAddress() (string, error) {
	raw, err := r.run("coin-cli", "getaddressesbyaccount", "")
	if err != nil {
		return "", err
	}

	dest := []string{}
	if err := unmarshalJSON(raw, &dest); err != nil {
		return "", err
	}

	return dest[0], nil
}

// GetBalance returns the balance of wallet
func GetBalance() (float64, error) {
	raw, err := r.run("coin-cli", "getbalance")
	if err != nil {
		return 0.0, err
	}

	balance, err := strconv.ParseFloat(string(raw), 64)
	if err != nil {
		return 0.0, err
	}

	return balance, nil
}

// SendTo send amount of coin to the dest address
func SendTo(address string, amount float64) (string, error) {
	raw, err := r.run("coin-cli", "sendtoaddress", address, fmt.Sprint(amount), "comment", "comment-to")
	if err != nil {
		return "", err
	}

	transactionID := string(raw)
	return transactionID, nil
}

func unmarshalJSON(data []byte, dest interface{}) error {
	if err := json.Unmarshal(data, dest); err != nil {
		logrus.WithFields(logrus.Fields{
			"json":  string(data),
			"error": err,
		}).Debug("failed to unmarshal json")
		return err
	}

	return nil
}
