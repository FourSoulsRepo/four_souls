package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

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
