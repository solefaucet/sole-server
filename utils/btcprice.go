package utils

import (
	"encoding/json"

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

	return MachineReadableUSD(m["USD"].Last), nil
}
