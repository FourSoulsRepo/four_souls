package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

const lost = "player 1 (the_lost)"

const lost2 = "player 2 (the_lost)"

func TestHolyMantleOnYourTurn(t *testing.T) {
	tb := table(t, seat("the_lost", "holy_mantle"), seat("cain"))
	tb.Activate(0, "holy_mantle", 0, lost)
	tb.G.DealDamageTo(engine.Target{Player: 0, IsPlayer: true}, 1, 1, 0)
	ev := tb.Settle()
	if tb.G.Players[0].Dead {
		t.Fatal("the Lost died; Holy Mantle should prevent it")
	}
	if !enginetest.Has(ev, engine.EvTurnEndedEarly, 0) {
		t.Error("on the Lost's own turn the turn should end (R-MECH-52)")
	}
	// The ward is used up: a second death this turn is not prevented.
	tb.G.DealDamageTo(engine.Target{Player: 0, IsPlayer: true}, 1, 1, 0)
	tb.Settle()
	if !tb.G.Players[0].Dead {
		t.Error("the ward prevents only the next death")
	}
}

func TestHolyMantleOnAnotherPlayer(t *testing.T) {
	tb := table(t, seat("isaac"), seat("the_lost", "holy_mantle"))
	tb.Start(0, "isaac", 0)
	tb.Pass(0) // player 2 responds on Isaac's turn
	tb.Activate(1, "holy_mantle", 0, "player 1 (isaac)")
	tb.G.DealDamageTo(engine.Target{Player: 0, IsPlayer: true}, 2, 1, 0)
	ev := tb.Settle()
	if tb.G.Players[0].Dead {
		t.Error("Isaac died; Holy Mantle should prevent it")
	}
	if !enginetest.Has(ev, engine.EvTurnEndedEarly, 0) {
		t.Error("it is Isaac's turn, so the turn ends")
	}
	if !tb.G.Object(tb.Find(1, "holy_mantle")).Charged {
		t.Log("Holy Mantle recharges at the end of its controller's turn, not now")
	}
}

func TestHolyMantleWardEndsWithTheTurn(t *testing.T) {
	tb := table(t, seat("isaac"), seat("the_lost", "holy_mantle"))
	tb.Start(0, "isaac", 0)
	tb.Pass(0)
	tb.Activate(1, "holy_mantle", 0, lost2)
	tb.EndTurn()
	tb.G.DealDamageTo(engine.Target{Player: 1, IsPlayer: true}, 1, 0, 0)
	settleAny(tb)
	if !tb.G.Players[1].Dead {
		t.Error("the ward lasts only this turn")
	}
}
