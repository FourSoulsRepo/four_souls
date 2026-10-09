package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestDankGlobin(t *testing.T) {
	tb := slayTable(t, "dank_globin", seat("isaac"), engine.SituationPlayer{Character: "cain", Hand: items("a_dime", "bomb")})
	killWithSixes(t, tb, foe, "a_dime", "bomb")
	if h := len(tb.G.Players[1].Hand); h != 0 {
		t.Errorf("Cain's hand = %d, want 0", h)
	}
}
