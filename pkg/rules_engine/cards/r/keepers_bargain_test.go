package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestKeepersBargain(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("the_miser", "keepers_bargain"), seat("cain"), seat("isaac")}})
	tb.Activate(0, "keepers_bargain", 0, "yes", "no")
	if tb.G.Players[0].Cents != 2 || tb.G.Players[1].Cents != 1 || tb.G.Players[2].Cents != 0 {
		t.Error("Cain took 1¢: the Miser gains 1 + 1")
	}
}
