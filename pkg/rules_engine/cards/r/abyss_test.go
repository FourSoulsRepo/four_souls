package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestAbyss(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("the_empty", "abyss", "spoon_bender", "jawbone"), seat("cain")}})
	tb.Activate(0, "abyss", 0, "spoon_bender")
	if rule := tb.Refused(0, "abyss", 0); rule != "R-ABIL-07" {
		t.Errorf("↷: a second use in one turn is refused by %s", rule)
	}
	tb.EndTurn()
	tb.EndTurn()
	tb.Activate(0, "abyss", 0, "jawbone")
	if tb.G.PlayerATK(0) != 2 || len(tb.G.Players[0].InPlay) != 1 {
		t.Errorf("2 items destroyed: 2 counters, +1 ATK; ATK %d", tb.G.PlayerATK(0))
	}
}

// TestAbyssOnBagLunchKeepsHP: a lost {HP} bonus lowers max HP and removes
// that much damage, so the HP left stays (R-MECH-21, 19.46).
func TestAbyssOnBagLunchKeepsHP(t *testing.T) {
	empty := seat("the_empty", "abyss", "bag_lunch")
	empty.Damage = 2
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{empty, seat("cain")}})
	tb.G.AddCountersTo(tb.Find(0, "bag_lunch"), 1)
	if got := tb.G.PlayerHP(0); got != 1 {
		t.Fatalf("HP = %d, want 1 of 3", got)
	}
	tb.Activate(0, "abyss", 0, "bag_lunch")
	if hasItem(tb.G, 0, "bag_lunch") {
		t.Fatal("Bag Lunch with a counter is destroyed")
	}
	if got, dmg := tb.G.PlayerHP(0), tb.G.Players[0].Damage; got != 1 || dmg != 1 {
		t.Errorf("HP %d, damage %d; want 1 of 2, damage 1", got, dmg)
	}
	if len(tb.G.Stack) != 0 || tb.G.Players[0].Dead {
		t.Error("the player keeps 1 HP; no death")
	}
}
