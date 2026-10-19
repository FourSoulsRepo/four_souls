package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestMamaMega(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", "mama_mega"), seat("cain")},
		Monsters: []engine.CardRef{"pooter", "fly"},
	})
	tb.Start(0, "mama_mega", 0)
	settleAny(tb) // the order of the monsters, then the death penalties
	if hasItem(tb.G, 0, "mama_mega") {
		t.Error("it is destroyed")
	}
	if !tb.G.Players[0].Dead || !tb.G.Players[1].Dead {
		t.Error("3 damage kills both players")
	}
	if !inDiscard(tb.G, engine.MonsterDeck, "pooter") || !inDiscard(tb.G, engine.MonsterDeck, "fly") {
		t.Error("3 damage kills both monsters")
	}
}
