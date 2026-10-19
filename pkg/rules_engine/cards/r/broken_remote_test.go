package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBrokenRemote(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", "broken_remote", "sack_of_pennies"), seat("cain")},
		Monsters: []engine.CardRef{"fly", "leech"},
		Rooms:    []engine.CardRef{"spider_webs"},
	})
	tb.Activate(0, "sack_of_pennies", 0, "spider_webs")
	if !inDiscard(tb.G, engine.RoomDeck, "spider_webs") {
		t.Error("activating an item lets Isaac put the room into discard")
	}
}
