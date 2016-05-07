package v1

import "testing"

func TestMax(t *testing.T) {
	if v := max(1, 2); v != 2 {
		t.Errorf("Max(1, 2) should be 2 but get %v", v)
	}

	if v := max(2, 1); v != 2 {
		t.Errorf("Max(1, 2) should be 2 but get %v", v)
	}
}
