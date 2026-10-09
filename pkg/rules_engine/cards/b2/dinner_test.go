package b2

import "testing"

func TestDinner(t *testing.T) {
	if hp := itemTable(t, items("dinner", "breakfast")).G.PlayerHP(0); hp != 4 {
		t.Errorf("HP = %d, want 4", hp)
	}
}
