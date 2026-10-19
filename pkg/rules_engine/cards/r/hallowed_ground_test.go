package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestHallowedGround(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", "hallowed_ground"), seat("cain")},
		Monsters: []engine.CardRef{"pooter"},
	})
	tb.Attack("pooter", 2, 6, 6) // a 2 misses, then two hits kill it
	ground := tb.Find(0, "hallowed_ground")
	if got := tb.G.Object(ground).CountersOf(""); got != 1 {
		t.Fatalf("counters = %d, want 1", got)
	}
	tb.G.AddCountersTo(ground, 2)
	tb.G.PutOnDeck(engine.LootDeck, "bomb")
	tb.G.MillTop(engine.LootDeck)
	tb.Activate(0, "hallowed_ground", 1)
	if !handHas(tb.G, 0, "bomb") {
		t.Error("the top card of the loot discard comes into the hand")
	}
}
