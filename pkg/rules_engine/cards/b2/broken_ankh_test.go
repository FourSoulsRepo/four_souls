package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestBrokenAnkh(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"broken_ankh"}, Hand: []engine.CardRef{"xiii_death"}},
		seat("cain"),
	}}, Set)
	tb.G.ForceRolls(6)
	tb.Play(0, "xiii_death", me)
	pl := tb.G.Players[0]
	if pl.Dead || tb.G.PlayerHP(0) != 1 {
		t.Errorf("dead %v HP %d, want alive at 1 HP", pl.Dead, tb.G.PlayerHP(0))
	}
	if !tb.G.Turn.EndDeclared && tb.G.Turn.Active == 0 && tb.G.Turn.Step == engine.StepAction {
		t.Error("the turn did not end")
	}

	tb = enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"broken_ankh"}, Hand: []engine.CardRef{"xiii_death"}},
		seat("cain"),
	}}, Set)
	tb.G.ForceRolls(5)
	tb.Play(0, "xiii_death", me, "broken_ankh") // the ankh is lost to the penalty

	if !tb.G.Players[0].Dead {
		t.Error("a 5 prevented the death")
	}
}
