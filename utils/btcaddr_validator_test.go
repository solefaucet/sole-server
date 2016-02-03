package utils

import "testing"

func TestBitcoinAddressValidator(t *testing.T) {
	testdata := []struct {
		addr  string
		valid bool
	}{
		{"", false},
		{"てΩ:wめまし", false},
		{"nEFJFaeATfp2442TGcHS5mgadXJjsSSP2T", false},
		{"1EFJFaeATfp2442TGcHS5mgadXJjsSSP2Ttoo_long", false},
		{"1EFJFaeATfp2442TGcHS5mgadXJjsSSP2T", true},
	}

	for _, v := range testdata {
		if valid, _ := ValidateBitcoinAddress(v.addr); valid != v.valid {
			expected := "valid"
			but := "invalid"
			if !v.valid {
				expected = "invalid"
				but = "valid"
			}
			t.Errorf("address %s should be %s but %s", v.addr, expected, but)
		}
	}
}
