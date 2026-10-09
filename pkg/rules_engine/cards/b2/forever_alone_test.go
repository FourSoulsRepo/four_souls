package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

const (
	steal    = "Steal 1¢ from another player."
	lookTop  = "Look at the top card of a deck."
	discLoot = "Discard a loot card, then loot 1."
)

func TestForeverAloneSteals(t *testing.T) {
	tb := table(t, seat("blue_baby", "forever_alone"), seat("cain"))
	tb.G.Players[0].Cents, tb.G.Players[1].Cents = 0, 3
	tb.Activate(0, "forever_alone", 0, steal, "player 2 (cain)")
	if a, b := tb.G.Players[0].Cents, tb.G.Players[1].Cents; a != 1 || b != 2 {
		t.Errorf("cents %d and %d, want 1 and 2", a, b)
	}
}

func TestForeverAloneLooksAtADeck(t *testing.T) {
	tb := table(t, seat("blue_baby", "forever_alone"), seat("cain"))
	top := tb.G.Object(tb.G.DeckTop(engine.TreasureDeck, 1)[0]).Card
	ev := tb.Activate(0, "forever_alone", 0, lookTop, "treasure deck")
	saw := false
	for _, e := range ev {
		if e.Kind == engine.EvLookedAt && e.Player == 0 && e.Card == top && e.Private {
			saw = true
		}
	}
	if !saw {
		t.Errorf("no private look at %s", top)
	}
}

func TestForeverAloneDiscardsThenLoots(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "blue_baby", Items: []engine.CardRef{"forever_alone"}, Hand: []engine.CardRef{"a_dime"}},
		seat("cain"),
	}}, Set)
	tb.Activate(0, "forever_alone", 0, discLoot, "a_dime")
	if len(tb.G.Players[0].Hand) != 1 {
		t.Errorf("hand has %d cards, want 1", len(tb.G.Players[0].Hand))
	}
	if d := tb.G.Discards[engine.LootDeck]; len(d) != 1 || tb.G.Object(d[0]).Card != "a_dime" {
		t.Error("A Dime!! is not in the loot discard")
	}
}

func TestForeverAloneRechargesOnDamage(t *testing.T) {
	tb := monsterTable(t, []engine.CardRef{"fly"}, seat("blue_baby", "forever_alone"), seat("cain"))
	tb.G.Players[1].Cents = 3
	tb.Activate(0, "forever_alone", 0, steal, "player 2 (cain)")
	item := tb.Find(0, "forever_alone")
	if tb.G.Object(item).Charged {
		t.Fatal("charged after use")
	}
	tb.Attack("fly", 1, 6) // a miss: the fly deals 1 damage
	if !tb.G.Object(item).Charged {
		t.Error("taking damage did not recharge Forever Alone")
	}
}
