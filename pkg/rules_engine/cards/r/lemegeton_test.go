package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestLemegeton(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("the_zealot", "lemegeton", "spoon_bender"), seat("cain")}})
	tb.Activate(0, "lemegeton", 0)
	tb.G.Object(tb.Find(0, "lemegeton")).Charged = true
	tb.Activate(0, "lemegeton", 0)
	if len(tb.G.Players[0].InPlay) != 4 {
		t.Fatal("two treasures gained")
	}
	// Both destructions may pick the Spoon Bender: the second does nothing (Q-15).
	tb.EndTurn("spoon_bender", "spoon_bender")
	if len(tb.G.Players[0].InPlay) != 3 || hasItem(tb.G, 0, "spoon_bender") {
		t.Errorf("one item destroyed; items %d", len(tb.G.Players[0].InPlay))
	}
}
