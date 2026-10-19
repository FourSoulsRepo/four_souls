package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestHagalaz(t *testing.T) {
	tb := lootTable(t, "hagalaz")
	tb.G.PutOnDeck(engine.MonsterDeck, "pooter")
	tb.Play(0, "hagalaz", "monster deck", "clotty", "pooter", "clotty", "clotty", "clotty", "fly")
	// The fly's slot takes the top clotty, so the pooter is now on top.
	if top := tb.G.DeckTop(engine.MonsterDeck, 1); tb.G.Object(top[0]).Card != "pooter" {
		t.Error("the pooter went back second")
	}
	if !inDiscard(tb.G, engine.MonsterDeck, "fly") {
		t.Error("then the fly goes to discard")
	}
}
