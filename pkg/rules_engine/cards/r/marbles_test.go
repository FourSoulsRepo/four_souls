package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestMarbles(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "marbles", "jar_of_flies"), seat("cain")}})
	jar := tb.Find(0, "jar_of_flies")
	tb.G.AddCountersTo(jar, 1)
	tb.Activate(0, "marbles", 0, "jar_of_flies", "add a counter")
	if got := tb.G.Object(jar).CountersOf(""); got != 2 {
		t.Errorf("counters = %d, want 2", got)
	}
}
