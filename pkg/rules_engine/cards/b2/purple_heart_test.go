package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestPurpleHeart(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{seat("isaac", "purple_heart"), seat("cain")},
		Active:  1,
	}, Set)
	top := tb.G.DeckTop(engine.MonsterDeck, 1)[0]
	tb.EndTurn("put " + string(tb.G.Object(top).Card) + " on the bottom")
	if tb.G.Decks[engine.MonsterDeck][0] != top {
		t.Error("the top card did not go to the bottom")
	}
}
