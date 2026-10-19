package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTheClicker(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "the_d6", "the_clicker"), seat("cain")}})
	tb.G.Players[0].Damage = 1
	tb.Activate(0, "the_clicker", 0, "samson")
	if c := tb.G.Object(tb.G.Players[0].Character); c.Card != "samson" || !c.Charged {
		t.Fatalf("character = %s (charged %v), want Samson, charged (Q-18)", c.Card, c.Charged)
	}
	if hasItem(tb.G, 0, "the_d6") || !hasItem(tb.G, 0, "blood_lust") {
		t.Error("the D6 leaves; Samson's Blood Lust comes")
	}
	if tb.G.Players[0].Damage != 1 {
		t.Error("damage stays with the player (R-MECH-19)")
	}
	for _, ch := range tb.G.CharactersOutside() {
		if ch == "samson" {
			t.Error("Samson is in play now, not outside")
		}
	}
}
