package b2

import "testing"

func TestBreakfast(t *testing.T) {
	if hp := itemTable(t, items("breakfast")).G.PlayerHP(0); hp != 3 {
		t.Errorf("HP = %d, want 3", hp)
	}
}
