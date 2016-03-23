package utils

import "fmt"

var (
	toHuman   func(int64) float64
	toMachine func(float64) int64
)

// InitializePriceConverter initialize price converter only once
func InitializePriceConverter(typ string) error {
	switch typ {
	case "usd":
		toHuman = humanReadableUSD
		toMachine = machineReadableUSD
	case "btc":
		toHuman = humanReadableBTC
		toMachine = machineReadableBTC
	default:
		return fmt.Errorf("converter %s not exist", typ)
	}

	return nil
}

// ToHuman converts to human readable price
func ToHuman(i int64) float64 {
	return toHuman(i)
}

// ToMachine converts to what machine want to save in db
func ToMachine(f float64) int64 {
	return toMachine(f)
}

const k10 = 1e4

func humanReadableUSD(v int64) float64 {
	return float64(v) / k10
}

func machineReadableUSD(v float64) int64 {
	return int64(v * k10)
}

const satonish = 1e8

func humanReadableBTC(v int64) float64 {
	return float64(v) / satonish
}

func machineReadableBTC(v float64) int64 {
	return int64(v * satonish)
}
