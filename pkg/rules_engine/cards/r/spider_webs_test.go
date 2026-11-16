package r

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSpiderWebs(t *testing.T) {
	tb := roomTable(t, "spider_webs")
	tb.EndTurn()
	if tb.G.Object(tb.G.Players[0].Character).Charged {
		t.Error("Isaac's character is deactivated at the end of his turn")
	}
}

// The deactivation has its event (19.176).
func TestSpiderWebsDeactivatedEvent(t *testing.T) {
	tb := roomTable(t, "spider_webs")
	char := tb.G.Players[0].Character
	ev := tb.EndTurn()
	if !slices.ContainsFunc(ev, func(e engine.Event) bool { return e.Kind == engine.EvDeactivated && e.Object == char }) {
		t.Error("no deactivated event names Isaac's character")
	}
}
