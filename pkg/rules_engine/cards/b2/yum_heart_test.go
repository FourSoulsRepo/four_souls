package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestYumHeart(t *testing.T) {
	tb := monsterTable(t, []engine.CardRef{"fly"}, seat("maggy", "yum_heart"), seat("cain"))
	tb.Activate(0, "yum_heart", 0, "player 1 (maggy)")
	// Miss (fly's 1 damage is prevented), miss (1 damage), then kill it.
	ev := tb.Attack("fly", 1, 1, 6)
	if !enginetest.Has(ev, engine.EvPrevented, 0) {
		t.Error("no damage was prevented")
	}
	if got := tb.G.Players[0].Damage; got != 1 {
		t.Errorf("damage = %d, want 1: only the first instance is prevented", got)
	}
	tb.EndTurn()
	if !tb.G.Object(tb.Find(0, "yum_heart")).Charged {
		t.Error("Yum Heart did not recharge at the end of the turn")
	}
	if len(tb.G.Shields) != 0 {
		t.Error("shields last past the end of the turn")
	}
}
