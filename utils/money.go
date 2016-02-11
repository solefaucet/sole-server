package utils

const k10 = 1e4

// HumanReadableUSD converts MachineReadableUSD to HumanReadableUSD
func HumanReadableUSD(v int64) float64 {
	return float64(v) / k10
}

// MachineReadableUSD converts HumanReadableUSD to HumanReadableUSD
func MachineReadableUSD(v float64) int64 {
	return int64(v * k10)
}

const satonish = 1e8

// HumanReadableBTC converts MachineReadableBTC to HumanReadableBTC
func HumanReadableBTC(v int64) float64 {
	return float64(v) / satonish
}

// MachineReadableBTC converts HumanReadableBTC to MachineReadableBTC
func MachineReadableBTC(v float64) int64 {
	return int64(v * satonish)
}
