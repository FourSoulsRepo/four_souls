package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTheFallen(t *testing.T) {
	tb := fight(t, "the_fallen", "fly")
	tb.G.PutOnDeck(engine.MonsterDeck, "dip", "spider", "trite", "leech", "fatty") // fatty on top
	tb.G.KillObject(topID(tb, 0))
	tb.Settle("dip", "spider", "trite", "leech", "fatty")
	// The slot refills from the top: the dip.
	if top(tb, 0) != "dip" {
		t.Errorf("slot 1 = %s, want the dip put on top", top(tb, 0))
	}
}
