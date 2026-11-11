package r

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/b2"
	"github.com/FourSoulsRepo/rules_engine/cards/fsp2"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// An item with a gained ↷ ability is a ↷ item (19.150): Golden Trinket
// with gold on a ↷ item, Glitch with a shop item's abilities.

// trinketOn: Isaac's Golden Trinket put a gold counter on the Mystery
// Sack; he has the other items too.
func trinketOn(t *testing.T, tb *enginetest.Table) engine.ObjectID {
	t.Helper()
	tb.Play(0, "golden_trinket", "mystery_sack")
	return tb.Find(0, "golden_trinket")
}

// The owner's line: the trinket spends its gained ↷, a 2 is rolled, and
// Charged Baby offers the trinket to recharge.
func TestChargedBabyRechargesGoldenTrinket(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"charged_baby", "mystery_sack"}, Hand: []engine.CardRef{"golden_trinket"}},
		seat("cain"),
	}})
	trinket := trinketOn(t, tb)
	tb.G.ForceRolls(2)
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentActivate, Objects: []engine.ObjectID{trinket}, Choice: 1}); err != nil {
		t.Fatal(err)
	}
	tb.UntilChoose()
	want := []string{"golden_trinket", "ask.none"}
	if got := tb.Options(); !slices.Equal(got, want) {
		t.Fatalf("options = %v, want %v", got, want)
	}
	tb.Choose("golden_trinket")
	if !tb.G.Object(trinket).Charged {
		t.Error("the trinket is recharged")
	}
}

// Glitch has a shop item's ↷: it can be deactivated, then recharged.
func TestGlitchGainedTapIsRechargeable(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "cain", Cents: 5}, seat("the_capricious", "glitch")},
		Shop:    []engine.CardRef{"jawbone", "spoon_bender"},
	})
	glitch := tb.Find(1, "glitch")
	if engine.Deactivatable(tb.G, glitch) {
		t.Fatal("before it gains the Jawbone, Glitch has no ↷")
	}
	tb.EndTurn("jawbone")
	if !engine.Deactivatable(tb.G, glitch) {
		t.Error("charged with the Jawbone's ↷: it can be deactivated")
	}
	tb.Activate(1, "glitch", 2, "player 1 (cain)")
	if !engine.Rechargeable(tb.G, glitch) {
		t.Error("its gained ↷ spent: it can be recharged")
	}
}

// The death penalty deactivates the trinket with a gained ↷.
func TestDeathDeactivatesGoldenTrinket(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"mystery_sack", "breakfast"}, Hand: []engine.CardRef{"golden_trinket"}},
		{Character: "cain", Cents: 5},
	}})
	trinket := trinketOn(t, tb)
	tb.G.KillPlayer(0, 1)
	tb.Settle("breakfast")
	if tb.G.Object(trinket).Charged {
		t.Error("the death penalty deactivates the trinket: it has a ↷")
	}
}

// Curse of Impulse deactivates the trinket with a gained ↷.
func TestCurseOfImpulseDeactivatesGoldenTrinket(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"mystery_sack"}, Hand: []engine.CardRef{"golden_trinket"}},
		seat("cain"),
	}, Monsters: []engine.CardRef{"fatty", "monstro"}}, b2.Set, Set, fsp2.Set)
	trinket := trinketOn(t, tb)
	tb.RevealFromDeck("curse_of_impulse", me)
	tb.EndTurn()
	if tb.G.Object(trinket).Charged {
		t.Error("Curse of Impulse deactivates the trinket: it has a ↷")
	}
}
