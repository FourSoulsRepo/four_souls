package r

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCharmedDip(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), {Character: "cain", Items: []engine.CardRef{"jawbone"}, Deactivated: []engine.CardRef{"jawbone"}}},
		Monsters: []engine.CardRef{"charmed_dip"},
	})
	tb.G.ForceRolls(6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "charmed_dip", foe, "jawbone")
	if !tb.G.Object(tb.Find(1, "jawbone")).Charged {
		t.Error("Cain recharges his Jawbone")
	}
}

// Cain is offered only his spent active item: not a charged one, and
// not a passive one even when it is spent (19.146).
func TestCharmedDipOffersOnlyRechargeableItems(t *testing.T) {
	cain := engine.SituationPlayer{Character: "cain", Items: []engine.CardRef{"jawbone", "the_d6", "breakfast"}, Deactivated: []engine.CardRef{"jawbone", "breakfast"}}
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "mystery_sack"), cain}, Monsters: []engine.CardRef{"charmed_dip"}})
	tb.G.ForceRolls(6)
	if _, err := tb.G.Settle(); err != nil {
		t.Fatal(err)
	}
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentAttack}); err != nil {
		t.Fatal(err)
	}
	if _, err := tb.G.Settle("charmed_dip", foe); err == nil {
		t.Fatal("no question: Cain should recharge an item")
	}
	if got, want := tb.Options(), []string{"jawbone"}; !slices.Equal(got, want) {
		t.Errorf("options = %v, want %v", got, want)
	}
}
