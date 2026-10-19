package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestWishbone(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: []engine.CardRef{"berkano"}}, seat("cain", "wishbone"),
	}, Shop: []engine.CardRef{"spoon_bender"}})
	tb.Play(0, "berkano", "wishbone", "spoon_bender")
	if got := len(tb.G.Players[1].InPlay); got != 2 {
		t.Errorf("Cain's items = %d, want the stolen one and +1 treasure", got)
	}
}
