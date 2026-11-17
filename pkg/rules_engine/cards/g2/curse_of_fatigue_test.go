package g2

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCurseOfFatigue(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{
			seat("isaac"),
			{Character: "cain", Items: []engine.CardRef{"the_bible", "crooked_penny"}, Deactivated: []engine.CardRef{"the_bible", "crooked_penny", "character"}},
		},
		Monsters: []engine.CardRef{"fatty", "monstro"},
	})
	tb.RevealFromDeck("curse_of_fatigue", foe)
	tb.EndTurn("crooked_penny") // Cain recharges only one item
	if !tb.G.Object(tb.G.Players[1].Character).Charged {
		t.Error("the character still recharges")
	}
	bible, penny := tb.G.Object(tb.Find(1, "the_bible")).Charged, tb.G.Object(tb.Find(1, "crooked_penny")).Charged
	if bible || !penny {
		t.Errorf("bible %v, penny %v: only the chosen item recharges", bible, penny)
	}
}

// The step that asks which item to recharge has the character's
// recharged event already (19.175).
func TestCurseOfFatigueCharacterRechargeEvent(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{
			seat("isaac"),
			{Character: "cain", Items: []engine.CardRef{"the_bible", "crooked_penny"}, Deactivated: []engine.CardRef{"the_bible", "crooked_penny", "character"}},
		},
		Monsters: []engine.CardRef{"fatty", "monstro"},
	})
	tb.RevealFromDeck("curse_of_fatigue", foe)
	in := engine.Intent{Player: tb.G.Turn.Active, Kind: engine.IntentEndTurn}
	for range 100 {
		ev, err := tb.G.Apply(in)
		if err != nil {
			t.Fatal(err)
		}
		w := tb.G.Prompt()
		if w.Kind == engine.PromptChoose {
			if !slices.ContainsFunc(ev, func(e engine.Event) bool { return e.Kind == engine.EvRecharged && e.Player == 1 }) {
				t.Errorf("the recharge choice opened with no recharged event; events %v", ev)
			}
			return
		}
		kind := engine.IntentPass
		if w.Kind == engine.PromptThrow {
			kind = engine.IntentThrow
		}
		in = engine.Intent{Player: w.Player, Kind: kind}
	}
	t.Fatal("the recharge choice never opened")
}
