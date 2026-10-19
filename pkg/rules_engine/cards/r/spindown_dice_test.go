package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSpindownDice(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{seat("the_broken", "spindown_dice"), seat("cain", "spoon_bender")},
		Shop:    []engine.CardRef{"jawbone"},
	})
	tb.Activate(0, "spindown_dice", 0, "spoon_bender", "jawbone")
	if hasItem(tb.G, 1, "spoon_bender") || !hasItem(tb.G, 1, "jawbone") {
		t.Error("Cain's item is destroyed; he steals the Jawbone")
	}
}

func TestSpindownDiceBagLunch(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{seat("the_broken", "spindown_dice"), seat("cain", "bag_lunch")},
		Shop:    []engine.CardRef{"jawbone"},
	})
	tb.Activate(0, "spindown_dice", 0, "bag_lunch")
	if !hasItem(tb.G, 1, "bag_lunch") || hasItem(tb.G, 1, "jawbone") {
		t.Error("Bag Lunch gets a counter instead; no destruction, no steal (Q-16)")
	}
}
