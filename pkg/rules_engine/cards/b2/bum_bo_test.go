package b2

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// Bum-Bo! and Counterfeit Penny both replace a gain of ¢. The player
// picks which applies first; each option says what it gives (19.17).
// The level each order reaches is in testdata/situations/bum_bo/.
func TestBumBoWithCounterfeitPenny(t *testing.T) {
	for _, tc := range []struct {
		first  engine.CardRef
		level  int
		unused engine.CardRef
	}{
		{"counterfeit_penny", 11, ""},
		{"bum_bo", 10, "counterfeit_penny"},
	} {
		tb := itemTable(t, items("bum_bo", "counterfeit_penny"), "a_dime")
		if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentPlayLoot, Objects: []engine.ObjectID{tb.Find(0, "a_dime")}}); err != nil {
			t.Fatal(err)
		}
		if _, err := tb.G.Settle(); err == nil {
			t.Fatal("no order was asked")
		}
		w := tb.G.Prompt()
		i := slices.Index(tb.Options(), string(tc.first))
		if w.Purpose != engine.ChooseReplacement || i < 0 {
			t.Fatalf("prompt %+v; want both cards as options", w)
		}
		o := tb.G.View(engine.Viewer{Kind: engine.ViewPlayer, Player: 0}).Options[i].Outcome
		want := engine.OutcomeAction{Kind: engine.ActAddCounters, Player: 0, Amount: tc.level, Card: "bum_bo"}
		if o == nil || len(o.Actions) != 1 || o.Actions[0] != want {
			t.Errorf("%s first: outcome %+v, want %+v", tc.first, o, want)
		}
		if tc.unused != "" && (o == nil || !slices.Equal(o.Unused, []engine.CardRef{tc.unused})) {
			t.Errorf("%s first: unused %+v, want %s", tc.first, o, tc.unused)
		}
	}
}
