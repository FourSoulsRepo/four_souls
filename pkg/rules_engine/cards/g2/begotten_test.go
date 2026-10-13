package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBegotten(t *testing.T) {
	tb := fight(t, "begotten", "fly")
	id := topID(tb, 0)
	if atk := tb.G.MonsterATK(id); atk != 1 {
		t.Errorf("ATK = %d at full HP, want 1", atk)
	}
	tb.G.DealDamageTo(engine.Target{Object: id}, 2, 0, 0)
	tb.Settle()
	if atk := tb.G.MonsterATK(id); atk != 2 {
		t.Errorf("ATK = %d at 1 HP, want 2", atk)
	}
}
