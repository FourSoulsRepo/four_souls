package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSmartFly(t *testing.T) {
	tb := table(t, seat("isaac", "smart_fly"), seat("cain"))
	tb.G.PutOnDeck(engine.MonsterDeck, "dip")
	tb.Activate(0, "smart_fly", 0, "monster deck", "put dip into discard")
	if top := tb.G.DeckTop(engine.MonsterDeck, 1); tb.G.Object(top[0]).Card == "dip" {
		t.Error("the dip should be discarded")
	}
	tb.G.DealDamageTo(engine.Target{Player: 0, IsPlayer: true}, 1, 1, 0)
	tb.Settle("yes")
	if !tb.G.Object(tb.Find(0, "smart_fly")).Charged {
		t.Error("taking damage may recharge it")
	}
}
