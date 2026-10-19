package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestAlgiz(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: []engine.CardRef{"algiz", "berkano"}}, seat("cain", "moms_purse"),
	}, Shop: []engine.CardRef{"spoon_bender"}})
	tb.G.Turn.LootPlays = 2
	tb.Play(0, "algiz", foe)
	tb.Play(0, "berkano", "moms_purse", "spoon_bender")
	if !hasItem(tb.G, 1, "moms_purse") {
		t.Error("the purse is eternal till end of turn: Berkano can't destroy it")
	}
	tb.EndTurn()
	if tb.G.Eternal(tb.Find(1, "moms_purse")) {
		t.Error("eternal should end with the turn")
	}
}
