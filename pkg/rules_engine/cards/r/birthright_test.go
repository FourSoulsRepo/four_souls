package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBirthright(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "the_d6", "birthright"), seat("cain")}})
	tb.Activate(0, "birthright", 0, "the_d6")
	id := tb.Find(0, "birthright")
	if tb.G.CardOf(id) != "the_d6" {
		t.Fatalf("Birthright acts as %s, want the_d6", tb.G.CardOf(id))
	}
	if tb.G.Eternal(id) {
		t.Error("the copy loses eternal")
	}
	if !tb.G.Eternal(tb.Find(0, "the_d6")) {
		t.Error("the copied item stays eternal")
	}
}
