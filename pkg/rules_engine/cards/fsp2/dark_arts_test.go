package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestDarkArts(t *testing.T) {
	cents := func(items ...engine.CardRef) int {
		tb := table(t, seat("dark_judas", items...), seat("cain"))
		tb.G.ForceRolls(6)
		tb.Activate(0, "mystery_sack", 0)
		return tb.G.Players[0].Cents
	}
	if with, without := cents("mystery_sack", "dark_arts"), cents("mystery_sack"); with-without != 2 {
		t.Errorf("a 6 gave %d¢ more, want 2", with-without)
	}
}

func TestDarkArtsOnDeath(t *testing.T) {
	tb := table(t, seat("dark_judas", "dark_arts"), seat("cain"))
	tb.G.DealDamageTo(engine.Target{Player: 1, IsPlayer: true}, 2, 0, 0)
	settleAny(tb)
	if n := len(tb.G.Players[0].Hand); n != 1 {
		t.Errorf("hand = %d, want 1 from Cain's death", n)
	}
}
