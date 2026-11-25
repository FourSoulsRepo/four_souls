package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// TestTheClickerCharacterIsSamson: the new character card is Samson, and
// Samson is no longer outside the game. No situation check names a
// character card; the rest is a file (the_clicker/).
func TestTheClickerCharacterIsSamson(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "the_d6", "the_clicker"), seat("cain")}})
	tb.Activate(0, "the_clicker", 0, "samson")
	if c := tb.G.Object(tb.G.Players[0].Character); c.Card != "samson" {
		t.Fatalf("character = %s, want Samson", c.Card)
	}
	for _, ch := range tb.G.CharactersOutside() {
		if ch == "samson" {
			t.Error("Samson is in play now, not outside")
		}
	}
}
