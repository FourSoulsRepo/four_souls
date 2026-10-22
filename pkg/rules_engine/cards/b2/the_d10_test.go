package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTheD10(t *testing.T) {
	tb := itemTable(t, items("the_d10", "mystery_sack"))
	fly, _ := tb.G.Monsters[0].TopOf()
	tb.G.ForceRolls(3)
	tb.Activate(0, "mystery_sack", 0, "slot 1")
	if top, _ := tb.G.Monsters[0].TopOf(); top == fly || tb.G.Object(fly).Zone.Kind != engine.ZoneCovered {
		t.Error("the fly was not covered")
	}
}

// An event put in a slot resolves at once, then the covered card is
// back on top (19.22).
func TestTheD10Event(t *testing.T) {
	tb := itemTable(t, items("the_d10", "mystery_sack"))
	tb.G.PutOnDeck(engine.MonsterDeck, "gold_chest")
	cents := tb.G.Players[0].Cents
	tb.G.ForceRolls(3, 5)
	tb.Activate(0, "mystery_sack", 0, "slot 1")
	if got := tb.G.Players[0].Cents - cents; got != 4+7 {
		t.Errorf("gained %d¢, want 11 (sack 4, chest 7)", got)
	}
	if top, _ := tb.G.Monsters[0].TopOf(); tb.G.Object(top).Card != "fly" {
		t.Errorf("slot 1 shows %s, want the fly", tb.G.Object(top).Card)
	}
	pile := tb.G.Discards[engine.MonsterDeck]
	if len(pile) == 0 || tb.G.Object(pile[len(pile)-1]).Card != "gold_chest" {
		t.Error("the gold chest is not in the monster discard")
	}
}
