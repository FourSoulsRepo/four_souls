package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSoulbond(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "cain", Hand: []engine.CardRef{"mega_bomb"}},
		{Character: "the_baleful", Items: []engine.CardRef{"soulbond"}, Hand: []engine.CardRef{"bomb"}, Cents: 3},
		seat("isaac"),
	}})
	tb.Play(0, "mega_bomb", "player 2 (the_baleful)", "player 3 (isaac)")
	if !tb.G.Players[2].Dead {
		t.Error("Isaac dies with him")
	}
	if tb.G.Players[1].Cents != 3 || len(tb.G.Players[1].Hand) != 1 {
		t.Error("his penalty takes no ¢ and no loot")
	}
}
