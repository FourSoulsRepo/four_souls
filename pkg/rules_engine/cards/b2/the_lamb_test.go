package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTheLamb(t *testing.T) {
	tb := slayTable(t, "the_lamb", seat("isaac"), engine.SituationPlayer{Character: "cain", Souls: items("monstro")})
	killWithSixes(t, tb, foe, "monstro")
	if s := tb.G.SoulValue(0); s != 3 { // the lamb's 2 and Cain's monstro
		t.Errorf("soul value %d, want 3", s)
	}
}
