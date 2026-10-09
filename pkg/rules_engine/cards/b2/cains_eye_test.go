package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestCainsEye(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{seat("isaac", "cains_eye"), seat("cain")},
		Active:  1,
	}, Set)
	top := tb.G.DeckTop(engine.LootDeck, 1)[0]
	tb.EndTurn("put " + string(tb.G.Object(top).Card) + " on the bottom")
	if tb.G.Decks[engine.LootDeck][0] != top {
		t.Error("the top card did not go to the bottom")
	}
}
