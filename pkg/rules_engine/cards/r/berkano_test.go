package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBerkano(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: []engine.CardRef{"berkano"}}, seat("cain", "moms_purse"),
	}, Shop: []engine.CardRef{"spoon_bender"}})
	tb.Play(0, "berkano", "moms_purse", "spoon_bender")
	if hasItem(tb.G, 1, "moms_purse") || !hasItem(tb.G, 1, "spoon_bender") {
		t.Error("Cain loses the purse and steals the shop's Spoon Bender")
	}
}

func TestBerkanoShopItem(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: []engine.CardRef{"berkano"}}, seat("cain"),
	}, Shop: []engine.CardRef{"spoon_bender"}})
	tb.Play(0, "berkano", "spoon_bender")
	if !inDiscard(tb.G, engine.TreasureDeck, "spoon_bender") {
		t.Error("the shop item goes to discard")
	}
	if len(tb.G.Players[0].InPlay)+len(tb.G.Players[1].InPlay) != 0 {
		t.Error("nobody controlled it: nobody steals")
	}
}
