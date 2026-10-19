package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCeremonialBlade(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("the_deceiver", "ceremonial_blade"), seat("cain")},
		Monsters: []engine.CardRef{"fly"},
	})
	tb.Attack("fly", 6)
	if len(tb.G.Players[0].Hand) != 1 {
		t.Error("a kill: loot 1")
	}
	tb.G.PutOnDeck(engine.LootDeck, "bomb")
	tb.EndTurn("bomb") // loot 1 (the bomb), then discard it
	if len(tb.G.Players[0].Hand) != 1 || handHas(tb.G, 0, "bomb") {
		t.Error("at the end of his turn: loot 1 and discard a loot card")
	}
}
