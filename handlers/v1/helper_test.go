package v1

import "testing"

func TestMax(t *testing.T) {
	if v := min(1, 2); v != 1 {
		t.Errorf("Min(1, 2) should be 1 but get %v", v)
	}

	if v := min(2, 1); v != 1 {
		t.Errorf("Min(1, 2) should be 1 but get %v", v)
	}
}
