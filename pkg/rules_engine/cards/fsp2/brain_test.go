package fsp2

import "testing"

func TestBrain(t *testing.T) {
	tb := fight(t, "brain", "fly")
	tb.G.ForceRolls(6, 6)
	tb.Do(attackIntent(tb), "brain")
	// Each 6 hurts: the trigger resolves before the hit's damage.
	if d := tb.G.Players[0].Damage; d != 2 {
		t.Errorf("damage = %d, want 2", d)
	}
}
