package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// TestGimpySaveAtDiscard saves while Gimpy's nested discard question
// waits and loads: the answer discards; it does not ask Gimpy's first
// question again (19.143, 19.177).
func TestGimpySaveAtDiscard(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "whore_of_babylon", Items: []engine.CardRef{"gimpy"}, Hand: []engine.CardRef{"a_penny", "a_penny"}}, seat("cain")},
	})
	tb.G.DealDamageTo(engine.Target{Player: 0, IsPlayer: true}, 1, 1, 0)
	tb.UntilChoose()
	tb.Choose("card.gimpy.loot")
	if !tb.G.NestedAsk() {
		t.Fatal("no nested question after Loot 1")
	}
	data, err := tb.G.Save()
	if err != nil {
		t.Fatal(err)
	}
	l, err := engine.Load(data, sets...)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := engine.Resume(data, engine.StateMeaning, sets...); err != nil {
		t.Fatalf("a server restart at the nested question: %v", err)
	}
	hand, discards, cents := len(l.Players[0].Hand), len(l.Discards[engine.LootDeck]), l.Players[0].Cents
	if hand != 3 {
		t.Fatalf("hand = %d after Loot 1, want 3", hand)
	}
	for _, g := range []*engine.Game{tb.G, l} {
		if _, err := g.Apply(engine.Intent{Player: 0, Kind: engine.IntentChoose, Choice: 0}); err != nil {
			t.Fatal(err)
		}
	}
	if got := len(l.Players[0].Hand); got != hand-1 {
		t.Errorf("hand = %d after the discard, want %d (a second loot?)", got, hand-1)
	}
	if got := len(l.Discards[engine.LootDeck]); got != discards+1 {
		t.Errorf("loot discards = %d, want %d", got, discards+1)
	}
	if got := l.Players[0].Cents; got != cents {
		t.Errorf("cents = %d, want %d: Gimpy's first question was asked again", got, cents)
	}
	a, _ := tb.G.Checksum() //nolint:errcheck // compared below
	b, _ := l.Checksum()    //nolint:errcheck // compared below
	if a != b {
		t.Error("the loaded game differs from the played one after the answer")
	}
}
