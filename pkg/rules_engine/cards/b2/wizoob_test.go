package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestWizoob(t *testing.T) {
	tb := slayTable(t, "wizoob", seat("isaac"), engine.SituationPlayer{Character: "cain", Souls: items("monstro")})
	killWithSixes(t, tb, foe, "monstro")
	if s := tb.G.SoulValue(1); s != 0 {
		t.Errorf("Cain's souls = %d, want 0", s)
	}
}
