package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBloodLust(t *testing.T) {
	tb := table(t, seat("samson", "blood_lust"), seat("cain"))
	tb.Activate(0, "blood_lust", 0, "player 2 (cain)")
	if got := tb.G.PlayerATK(1); got != 2 {
		t.Errorf("target's ATK = %d, want 2", got)
	}
	tb.EndTurn()
	if got := tb.G.PlayerATK(1); got != 1 {
		t.Errorf("ATK after the turn = %d, want 1 (till end of turn)", got)
	}
	if !tb.G.Object(tb.Find(0, "blood_lust")).Charged {
		t.Error("Blood Lust did not recharge at the end of the turn")
	}
}

func TestBloodLustOnAMonster(t *testing.T) {
	tb := monsterTable(t, []engine.CardRef{"fly"}, seat("samson", "blood_lust"), seat("cain"))
	fly, _ := tb.G.Monsters[0].TopOf()
	before := tb.G.MonsterATK(fly)
	tb.Activate(0, "blood_lust", 0, "fly")
	if got := tb.G.MonsterATK(fly); got != before+1 {
		t.Errorf("monster ATK = %d, want %d", got, before+1)
	}
}
