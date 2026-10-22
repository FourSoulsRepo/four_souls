package b2

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestXXJudgement(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: []engine.CardRef{"xx_judgement"}},
		{Character: "cain", Souls: []engine.CardRef{"monstro", "gurdy"}},
	}}, Set)
	tb.Play(0, "xx_judgement", foe, "gurdy")
	if s := tb.G.SoulValue(1); s != 1 {
		t.Errorf("soul value = %d, want 1", s)
	}
}

// TestXXJudgementSkipsEternalSouls: a soul that can't be destroyed is
// not offered (19.36, Q-23).
func TestXXJudgementSkipsEternalSouls(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: []engine.CardRef{"xx_judgement"}},
		{Character: "cain", Souls: []engine.CardRef{"monstro", "gurdy"}},
	}}, Set)
	tb.G.Object(tb.Find(1, "monstro")).Eternal = true
	in := engine.Intent{Player: 0, Kind: engine.IntentPlayLoot, Objects: []engine.ObjectID{tb.Find(0, "xx_judgement")}}
	if _, err := tb.G.Apply(in); err != nil {
		t.Fatal(err)
	}
	if _, err := tb.G.Settle(foe); err == nil {
		t.Fatal("no soul to choose")
	}
	if got := tb.Options(); !slices.Equal(got, []string{"gurdy"}) {
		t.Fatalf("options = %v, want [gurdy]", got)
	}
	tb.Choose("gurdy")
	tb.Settle()
	if s := tb.G.SoulValue(1); s != 1 {
		t.Errorf("soul value = %d, want 1", s)
	}
}
