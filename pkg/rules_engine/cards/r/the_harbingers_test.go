package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func harbingersTable(t *testing.T) *enginetest.Table {
	t.Helper()
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), seat("cain")},
		Monsters: []engine.CardRef{"fly", "leech"},
	})
	tb.G.PlaceFromOutside("the_harbingers", 0)
	return tb
}

func TestTheHarbingersCannotDie(t *testing.T) {
	tb := harbingersTable(t)
	h := tb.G.Monsters[len(tb.G.Monsters)-1].Cards[0]
	if tb.G.CardOf(h) != "the_harbingers" || len(tb.G.Monsters) != 3 {
		t.Fatal("Indomitable: it takes a new slot")
	}
	tb.G.Object(h).Damage = 3
	tb.Attack("the_harbingers", 6)
	if tb.G.Object(h).Zone.Kind != engine.ZoneInPlay || tb.G.Object(h).CountersOf("") != 1 || tb.G.Object(h).Damage != 0 {
		t.Error("instead of dying: a counter, healed to full")
	}
	if len(tb.G.Players[0].InPlay) != 1 || tb.G.Attack.On {
		t.Error("the active player gains +1 treasure; the attack is cancelled")
	}
}

func TestTheHarbingersFlipToTheBeast(t *testing.T) {
	tb := harbingersTable(t)
	h := tb.G.Monsters[len(tb.G.Monsters)-1].Cards[0]
	tb.G.AddCountersTo(h, 3)
	tb.G.Object(h).Damage = 3
	tb.G.ForceRolls(6, 2, 5) // the kill, then Isaac rolls 2, Cain 5
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "the_harbingers", "book_of_sin")
	if tb.G.CardOf(h) != "the_beast" {
		t.Fatal("4 counters: it flips to The Beast!")
	}
	if !tb.G.Players[0].Dead || tb.G.Players[1].Dead {
		t.Error("Isaac rolled lowest: he dies")
	}
}
