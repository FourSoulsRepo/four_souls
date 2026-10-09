package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestLazarusRags(t *testing.T) {
	tb := monsterTable(t, []engine.CardRef{"fly"}, seat("lazarus", "lazarus_rags"), seat("cain"))
	tb.G.Players[0].Cents = 3
	tb.Attack("fly", 1, 1) // two misses: 2 damage kills Lazarus
	pl := tb.G.Players[0]
	if !pl.Dead {
		t.Fatal("Lazarus is not dead")
	}
	if pl.Cents != 2 {
		t.Errorf("cents = %d, want 2: the penalty is paid first", pl.Cents)
	}
	if len(pl.InPlay) != 2 || !hasItem(tb.G, 0, "lazarus_rags") {
		t.Errorf("in play %v, want the rags (eternal, kept) and one new treasure", pl.InPlay)
	}
}
