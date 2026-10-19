package r

import "testing"

func TestTurdlings(t *testing.T) {
	tb := fight(t, "turdlings")
	tb.Attack("turdlings", 6, 2)
	if top(tb, 0) != "turdlings" || tb.G.Attack.On {
		t.Error("a 6 deals nothing; the 2nd roll (2) hits, then the attack is cancelled")
	}
	id := tb.G.Monsters[0].Cards[len(tb.G.Monsters[0].Cards)-1]
	if tb.G.Object(id).Damage != 1 {
		t.Errorf("damage %d, want 1 from the roll of 2", tb.G.Object(id).Damage)
	}
}
