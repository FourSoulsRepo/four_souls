package b2

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestSoulOfGuppy(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "guppys_head", "guppys_paw", "mystery_sack"), seat("cain")}}, Set)
	cond := Set.Cards[slices.IndexFunc(Set.Cards, func(c engine.CardDef) bool { return c.Ref == "soul_of_guppy" })].BonusSoul
	if !cond(tb.G, 0) || cond(tb.G, 1) {
		t.Error("the condition is wrong")
	}
}
