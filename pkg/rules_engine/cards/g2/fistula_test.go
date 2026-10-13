package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestFistula(t *testing.T) {
	tb := fight(t, "fistula", "fly")
	id := topID(tb, 0)
	tb.G.DealDamageTo(engine.Target{Object: id}, 1, 0, 0)
	tb.Settle()
	if dc := tb.G.Evasion(id); dc != 3 {
		t.Errorf("DC = %d after damage, want 2 + 1", dc)
	}
}
