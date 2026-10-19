package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestRedKey(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", "red_key"), seat("cain")},
		Monsters: []engine.CardRef{"fly", "leech"},
		Rooms:    []engine.CardRef{"blessing_of_greed"},
	})
	tb.G.PutOnDeck(engine.RoomDeck, "spider_webs")
	tb.Activate(0, "red_key", 0, "room deck", "spider_webs", "slot 1")
	if tb.G.Object(room(tb)).Card != "spider_webs" {
		t.Fatal("the spider webs cover the room")
	}
	tb.G.DiscardRoom(room(tb))
	if tb.G.Object(room(tb)).Card != "blessing_of_greed" {
		t.Error("the covered room comes back")
	}
}

func TestRedKeyNotDuringAnAttack(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", "red_key"), seat("cain")},
		Monsters: []engine.CardRef{"fly", "leech"},
	})
	tb.ThrowPending()
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentAttack}); err != nil {
		t.Fatal(err)
	}
	if rule := tb.Refused(0, "red_key", 0); rule != "R-ABIL-07" {
		t.Errorf("refused by %s, want the cost (R-ABIL-07)", rule)
	}
}
