package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestStevenTwo(t *testing.T) {
	tb := fight(t, "steven", "fly")
	tb.G.ForceRolls(2) // Isaac (HP 2) takes 2 from the trigger before the hit
	tb.Do(attackIntent(tb), "steven", "isaac")
	if !tb.G.Players[0].Dead {
		t.Error("Isaac should die from Steven's 2 damage")
	}
	if tb.G.Object(topID(tb, 0)).Damage != 0 {
		t.Error("the attack was cancelled by the death: Steven takes no damage (S-YUG)")
	}
}

func TestStevenSix(t *testing.T) {
	tb := fight(t, "steven", "fly")
	tb.G.ForceRolls(6)
	ev := tb.Do(attackIntent(tb), "steven")
	if !enginetest.Has(ev, engine.EvTurnEndedEarly, 0) {
		t.Error("a 6 should end the turn")
	}
}
