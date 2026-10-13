package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestRubberCement(t *testing.T) {
	tb := armed(t, items("rubber_cement"), "trite", "fly")
	tb.G.PutOnDeck(engine.MonsterDeck, "dip") // the fly's slot refills with a known card
	tb.G.ForceRolls(3, 2, 6)                  // a 3 misses the trite; the cement rolls 2: 1 damage to the fly
	tb.Do(attackIntent(tb), "trite", "fly")
	if top(tb, 1) != "dip" {
		t.Errorf("slot 2 = %s: the fly should die and the dip replace it", top(tb, 1))
	}
}
