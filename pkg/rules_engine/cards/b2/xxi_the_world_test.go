package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestXXITheWorld(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: []engine.CardRef{"xxi_the_world"}},
		{Character: "cain", Hand: []engine.CardRef{"a_dime"}},
	}}, Set)
	ev := tb.Play(0, "xxi_the_world")
	if !enginetest.Has(ev, engine.EvLookedAt, 0) {
		t.Error("no look at the other hand")
	}
	if n := len(tb.G.Players[0].Hand); n != 2 {
		t.Errorf("hand %d, want 2", n)
	}
}

// Each looked-at card names its hand's owner in By, seat plus 1 (19.71).
func TestXXITheWorldNamesEachHand(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: []engine.CardRef{"xxi_the_world"}},
		{Character: "cain", Hand: []engine.CardRef{"a_dime"}},
		{Character: "judas", Hand: []engine.CardRef{"a_penny", "a_nickel"}},
	}}, Set)
	ev := tb.Play(0, "xxi_the_world")
	got := map[engine.CardRef]int{}
	for _, e := range ev {
		if e.Kind == engine.EvLookedAt {
			if e.Player != 0 {
				t.Errorf("%s looked at by %d, want 0", e.Card, e.Player)
			}
			got[e.Card] = e.By
		}
	}
	want := map[engine.CardRef]int{"a_dime": 2, "a_penny": 3, "a_nickel": 3}
	if len(got) != len(want) {
		t.Errorf("looked at %v, want %v", got, want)
	}
	for card, by := range want {
		if got[card] != by {
			t.Errorf("%s: by %d, want %d", card, got[card], by)
		}
	}
}
