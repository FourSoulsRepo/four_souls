package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestAnimaSola(t *testing.T) {
	tb2 := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("the_deserter", "anima_sola"), seat("cain")}})
	sola := tb2.Find(0, "anima_sola")
	tb2.Activate(0, "anima_sola", 0)
	if tb2.G.CardOf(sola) != "the_revenant" {
		t.Fatal("it flips to The Revenant")
	}
	tb2.G.Object(sola).Charged = true
	if rule := tb2.Refused(0, "anima_sola", 0); rule != "R-ABIL-08" {
		t.Errorf("The Revenant has no ↷: refused by %s", rule)
	}
}

func TestTheRevenantAttacked(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac"), seat("the_deserter", "anima_sola")}})
	sola := tb.Find(1, "anima_sola")
	tb.G.Object(sola).Flipped = true
	tb.G.ForceRolls(3, 3)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "anima_sola")
	if tb.G.CardOf(sola) != "anima_sola" {
		t.Error("at 0 HP it flips back to Anima Sola")
	}
	// Isaac gains its reward (loot 2); the Deserter loots 1 as it flips.
	if len(tb.G.Players[0].Hand) != 2 || len(tb.G.Players[1].Hand) != 1 {
		t.Errorf("hands %d and %d, want 2 and 1", len(tb.G.Players[0].Hand), len(tb.G.Players[1].Hand))
	}
	if tb.G.Turn.Attacks != 1 {
		t.Error("Isaac may attack an additional time")
	}
}
