package g2

import (
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
