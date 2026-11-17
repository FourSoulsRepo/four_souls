package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestPlanC(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"plan_c"}}, seat("cain")},
		Monsters: []engine.CardRef{"monstro", "fatty"},
	})
	tb.Activate(0, "plan_c", 0, "monstro")
	if hasItem(tb.G, 0, "plan_c") {
		t.Error("Plan C was not destroyed")
	}
	if top, _ := tb.G.Monsters[0].TopOf(); tb.G.Object(top).Card == "monstro" {
		t.Error("Monstro was not killed")
	}
	if len(tb.G.Players[0].InPlay) == 0 {
		t.Error("Isaac should gain Monstro's soul")
	}
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentEndTurn}); err != nil {
		t.Fatal(err)
	}
	settleAny(tb)
	for range 20 {
		if tb.G.Players[0].Dead && tb.G.Players[1].Dead {
			break
		}
		switch w := tb.G.Prompt(); w.Kind {
		case engine.PromptChoose:
			tb.Choose(tb.Options()[0])
		case engine.PromptPriority:
			tb.Pass(int(w.Player))
		default:
		}
	}
	if !tb.G.Players[0].Dead || !tb.G.Players[1].Dead {
		t.Errorf("dead = %v, %v: at the end of the turn all players die", tb.G.Players[0].Dead, tb.G.Players[1].Dead)
	}
}

// TestPlanCSaveAtKill saves while Plan C's "kill a monster" question
// waits inside DestroyThis and loads: the load used to panic (19.177).
func TestPlanCSaveAtKill(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"plan_c"}}, seat("cain")},
		Monsters: []engine.CardRef{"monstro", "fatty"},
	})
	tb.Start(0, "plan_c", 0)
	tb.UntilChoose()
	if !tb.G.NestedAsk() {
		t.Fatal("the kill question is not nested in DestroyThis")
	}
	data, err := tb.G.Save()
	if err != nil {
		t.Fatal(err)
	}
	l, err := engine.Load(data, sets...)
	if err != nil {
		t.Fatal(err)
	}
	played := tb.G
	tb.G = l
	tb.Choose("monstro")
	tb.Settle()
	if hasItem(l, 0, "plan_c") {
		t.Error("Plan C was not destroyed")
	}
	if top(tb, 0) == "monstro" {
		t.Error("Monstro was not killed")
	}
	tb.G = played
	tb.Choose("monstro")
	tb.Settle()
	a, _ := played.Checksum() //nolint:errcheck // compared below
	b, _ := l.Checksum()      //nolint:errcheck // compared below
	if a != b {
		t.Error("the loaded game differs from the played one after the answer")
	}
}
