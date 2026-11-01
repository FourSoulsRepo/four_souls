package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestHeadlessHorseman(t *testing.T) {
	tb := fight(t, "headless_horseman", "fly")
	id := topID(tb, 0)
	tb.G.DealDamageTo(engine.Target{Object: id}, 3, 0, 0)
	tb.Settle()
	if top(tb, 0) != "headless_horseman" {
		t.Fatal("the first death this turn is prevented")
	}
	// Prevention is not healing: 0 HP, then the heal of 2.
	if hp, dc, atk := tb.G.HP(id), tb.G.Evasion(id), tb.G.MonsterATK(id); hp != 2 || dc != 4 || atk != 1 {
		t.Errorf("HP %d DC %d ATK %d, want 2 (0 + 2), 4, 1", hp, dc, atk)
	}
	tb.G.DealDamageTo(engine.Target{Object: id}, 3, 0, 0)
	settleAny(tb)
	if top(tb, 0) == "headless_horseman" {
		t.Error("the second death this turn goes through")
	}
}
