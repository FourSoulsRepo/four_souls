package b2

import "testing"

func TestBookOfSin(t *testing.T) {
	tb := itemTable(t, items("book_of_sin"))
	tb.G.ForceRolls(5)
	tb.Activate(0, "book_of_sin", 0)
	if hp := tb.G.PlayerHP(0); hp != 3 {
		t.Errorf("HP = %d, want 3", hp)
	}
}
