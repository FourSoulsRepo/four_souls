package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestHolyCard(t *testing.T) {
	tb := lootTable(t, "holy_card")
	tb.Play(0, "holy_card", me)
	tb.G.DealDamageTo(engine.Target{Player: 0, IsPlayer: true}, 2, 1, 0)
	ev := tb.Settle()
	if tb.G.Players[0].Dead {
		t.Error("Isaac died; Holy Card should prevent it")
	}
	if !enginetest.Has(ev, engine.EvTurnEndedEarly, 0) {
		t.Error("on Isaac's own turn the turn ends")
	}
}
