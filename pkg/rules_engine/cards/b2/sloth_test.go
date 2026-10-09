package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSloth(t *testing.T) {
	tb := slayTable(t, "sloth", engine.SituationPlayer{Character: "isaac", Hand: items("a_dime", "bomb")}, seat("cain"))
	killWithSixes(t, tb)
	if h := len(tb.G.Players[0].Hand); h != 0 {
		t.Errorf("the killer's hand = %d, want 0", h)
	}
}
