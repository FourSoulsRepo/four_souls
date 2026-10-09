package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestWeNeedToGoDeeper(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", "mr_boom"), seat("cain")},
		Monsters: items("fly", "leech"),
	}, Set)
	tb.Activate(0, "mr_boom", 0, "fly")
	tb.RevealFromDeck("we_need_to_go_deeper", "fly", "done")
	if top := tb.G.DeckTop(engine.MonsterDeck, 1); tb.G.Object(top[0]).Card != "fly" {
		t.Error("the fly is not on top of the monster deck")
	}
}
