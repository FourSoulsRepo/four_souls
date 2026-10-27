package b2

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBumBo(t *testing.T) {
	tb := itemTable(t, items("bum_bo"), "a_dime")
	tb.Play(0, "a_dime")
	bumbo := tb.G.Object(tb.Find(0, "bum_bo"))
	if c, l := tb.G.Players[0].Cents, bumbo.CountersOf(""); c != 0 || l != 10 {
		t.Errorf("cents %d level %d, want 0 and 10", c, l)
	}
	if a := tb.G.PlayerATK(0); a != 2 {
		t.Errorf("ATK at level 10 = %d, want 2", a)
	}
}

// Bum-Bo! and Counterfeit Penny both replace a gain of ¢. The player
// picks which applies first; each option says what it gives (19.17).
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
		tb.Settle(string(tc.first))
		if l := tb.G.Object(tb.Find(0, "bum_bo")).CountersOf(""); l != tc.level {
			t.Errorf("%s first: level %d, want %d", tc.first, l, tc.level)
		}
	}
}
