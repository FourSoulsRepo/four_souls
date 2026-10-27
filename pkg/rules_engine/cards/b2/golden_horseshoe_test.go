package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestGoldenHorseshoe(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{seat("isaac", "golden_horseshoe"), seat("cain")},
		Active:  1,
	}, Set)
	top := tb.G.DeckTop(engine.TreasureDeck, 1)[0]
	tb.EndTurn("ask.put_on_bottom " + string(tb.G.Object(top).Card))
	if tb.G.Decks[engine.TreasureDeck][0] != top {
		t.Error("the top card did not go to the bottom")
	}
}
