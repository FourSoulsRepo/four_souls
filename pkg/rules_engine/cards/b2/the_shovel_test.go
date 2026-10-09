package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTheShovel(t *testing.T) {
	tb := itemTable(t, items("the_shovel", "mr_boom"))
	tb.Activate(0, "mr_boom", 0, "fly")
	tb.Activate(0, "the_shovel", 0, "fly")
	if top := tb.G.DeckTop(engine.MonsterDeck, 1); tb.G.Object(top[0]).Card != "fly" {
		t.Error("the fly is not on top of the monster deck")
	}
}
