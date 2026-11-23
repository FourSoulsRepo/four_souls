package r

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// The deactivation has its event (19.176); the event's Object stays Go
// (24.21): the file checks only its kind and player.
func TestSpiderWebsDeactivatedEvent(t *testing.T) {
	tb := roomTable(t, "spider_webs")
	char := tb.G.Players[0].Character
	ev := tb.EndTurn()
	if !slices.ContainsFunc(ev, func(e engine.Event) bool { return e.Kind == engine.EvDeactivated && e.Object == char }) {
		t.Error("no deactivated event names Isaac's character")
	}
}
