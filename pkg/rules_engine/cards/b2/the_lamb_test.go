package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTheLamb(t *testing.T) {
	tb := slayTable(t, "the_lamb", seat("isaac"), engine.SituationPlayer{Character: "cain", Souls: items("monstro")})
	killWithSixes(t, tb, foe, "monstro")
	if s := tb.G.SoulValue(0); s != 2 { // the lamb's 1 and Cain's monstro
		t.Errorf("soul value %d, want 2", s)
	}
}
