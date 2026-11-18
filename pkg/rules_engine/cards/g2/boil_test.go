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

// TestBoilReward: Boil's printed reward box shows 3¢, no loot (19.224).
func TestBoilReward(t *testing.T) {
	tb := fight(t, "boil", "fly")
	tb.G.DealDamageTo(engine.Target{Object: topID(tb, 0)}, 1, 0, 0) // Boil at 1 HP
	tb.Settle()
	cents, hand := tb.G.Players[0].Cents, len(tb.G.Players[0].Hand)
	tb.G.ForceRolls(6) // a hit kills Boil
	tb.Do(attackIntent(tb), "boil")
	tb.Settle()
	if top(tb, 0) == "boil" {
		t.Fatal("Boil is still in its slot")
	}
	if got := tb.G.Players[0].Cents - cents; got != 3 {
		t.Errorf("cents gained = %d, want 3 (printed reward: 3¢)", got)
	}
	if got := len(tb.G.Players[0].Hand) - hand; got != 0 {
		t.Errorf("loot gained = %d, want 0 (printed reward: 3¢ only)", got)
	}
}
