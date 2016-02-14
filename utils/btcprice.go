package utils

import (
	"encoding/json"
	"errors"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/parnurzeal/gorequest"
)

const blockchainTickerURL = "https://blockchain.info/ticker"

// BitcoinPrice returns the lastest bitcoin price in 1 / 10,000 USD
func BitcoinPrice() (int64, error) {
	_, body, _ := gorequest.New().Get(blockchainTickerURL).EndBytes()
	return bitcoinPriceWithByteFromBlockchain(body)
}

func bitcoinPriceWithByteFromBlockchain(data []byte) (int64, error) {
	m := map[string]struct {
		Last float64 `json:"last"`
	}{}

	err := json.Unmarshal(data, &m)
	if err != nil {
		return 0, err
	}

	p := MachineReadableUSD(m["USD"].Last)
	if p == 0 {
		return 0, errors.New("bitcoin price cannot be 0")
	}

	return p, nil
}
