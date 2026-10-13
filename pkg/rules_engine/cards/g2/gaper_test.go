package g2

import "testing"

func TestGaper(t *testing.T) {
	tb := fight(t, "gaper", "fly")
	tb.G.ForceRolls(5, 6) // a 5 hits, the 6 kills
	tb.Do(attackIntent(tb), "gaper")
	if tb.G.Turn.MustAttacks != 1 || tb.G.Turn.MustAttackDeck != 1 {
		t.Errorf("must attacks = %d (deck %d), want 1 on the monster deck", tb.G.Turn.MustAttacks, tb.G.Turn.MustAttackDeck)
	}
	tb2 := fight(t, "gaper", "fly")
	tb2.G.ForceRolls(6, 5)
	tb2.Do(attackIntent(tb2), "gaper")
	if tb2.G.Turn.MustAttacks != 0 {
		t.Error("killed on a 5: no additional attack")
	}
}
