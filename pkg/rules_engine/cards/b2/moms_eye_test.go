package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestMomsEye(t *testing.T) {
	tb := slayTable(t, "moms_eye", seat("isaac"), engine.SituationPlayer{Character: "cain", Hand: items("bomb")})
	killWithSixes(t, tb, foe)
	_ = tb
}
