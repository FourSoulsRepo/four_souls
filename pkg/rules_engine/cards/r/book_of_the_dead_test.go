package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBookOfTheDead(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", "book_of_the_dead"), seat("cain")},
		Monsters: []engine.CardRef{"fly", "leech"},
	})
	tb.G.PutOnDeck(engine.MonsterDeck, "pooter")
	tb.G.MillTop(engine.MonsterDeck)
	tb.Activate(0, "book_of_the_dead", 0, "pooter", "ask.slot_card 2 leech")
	if got := top(tb, 1); got != "pooter" {
		t.Errorf("slot 2 = %s, want pooter", got)
	}
}
