package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestGurglings(t *testing.T) {
	tb := fight(t, "gurglings", "fly")
	id := topID(tb, 0)
	if dc := tb.G.Evasion(id); dc != 5 {
		t.Errorf("DC = %d at 4 HP, want 5", dc)
	}
	tb.G.DealDamageTo(engine.Target{Object: id}, 2, 0, 0)
	tb.Settle()
	if dc := tb.G.Evasion(id); dc != 4 {
		t.Errorf("DC = %d at 2 HP, want 4", dc)
	}
}
