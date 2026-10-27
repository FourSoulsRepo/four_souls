package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTheMirror(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), seat("cain")},
		Monsters: []engine.CardRef{"fly", "leech"},
		Rooms:    []engine.CardRef{"the_mirror"},
	})
	tb.EndTurn("ask.yes")
	if tb.G.PlayerHP(1) != 1 {
		t.Fatalf("Cain pays HP down to 1: HP %d", tb.G.PlayerHP(1))
	}
	tb.Attack("fly", 6)
	if len(itemsOf(tb.G, 1)) != 1 {
		t.Error("a monster died: Cain gains +1 treasure")
	}
}
