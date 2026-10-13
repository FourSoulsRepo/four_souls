package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBoil(t *testing.T) {
	tb := fight(t, "boil", "fly")
	tb.G.DealDamageTo(engine.Target{Object: topID(tb, 0)}, 1, 0, 0) // Boil at 1 HP
	tb.Settle()
	tb.G.ForceRolls(1, 6, 6) // miss: Boil hits Isaac and heals; then two hits
	ev := tb.Do(attackIntent(tb), "boil")
	if n := attackRolls(ev); n != 3 {
		t.Errorf("attack rolls = %d, want 3: Boil healed after its combat damage", n)
	}
}
