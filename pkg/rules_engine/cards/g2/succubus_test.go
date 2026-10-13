package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSuccubus(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"succubus"}, Cents: 3}, seat("cain"),
	}})
	tb.Activate(0, "succubus", 0)
	if got := tb.G.Turn.LootPlays; got != 3 {
		t.Errorf("loot plays = %d, want 3", got)
	}
	if c := tb.G.Players[0].Cents; c != 0 {
		t.Errorf("cents = %d, want 0", c)
	}
}
