package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// Golden Trinket with gold counters on Bum Friend and Tech X: both ↷
// abilities are its, and they share its one tap.
func TestGoldOnBumFriendAndTechX(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"golden_trinket", "bum_friend", "tech_x"}}, {Character: "cain", Cents: 5}},
	})
	trinket := tb.Find(0, "golden_trinket")
	tb.G.AddNamedCounters(tb.Find(0, "bum_friend"), gold, 1)
	tb.G.AddNamedCounters(tb.Find(0, "tech_x"), gold, 1)

	// Own trigger, Bum Friend's ↷, then Tech X's two.
	if got := len(tb.G.AbilitiesOf(trinket)); got != 4 {
		t.Fatalf("Golden Trinket has %d abilities, want 4", got)
	}
	tb.Activate(0, "golden_trinket", 1, "three_cents") // loots, then puts a card back on the deck
	if got := tb.G.Object(trinket).Charged; got {
		t.Error("the trinket is still charged after Bum Friend's ↷")
	}
	if rule := tb.Refused(0, "golden_trinket", 2); rule == "" {
		t.Error("Tech X's ↷ was allowed after Bum Friend's ↷ on the same turn: both use the trinket's one tap")
	}
}

// Bum-Bo!'s statics reach Golden Trinket with gold on Bum-Bo!, and read
// the trinket's own counters, not Bum-Bo!'s.
func TestBumBoStaticsOnGoldenTrinket(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"golden_trinket", "bum_bo"}}, {Character: "cain", Cents: 5}},
	})
	trinket := tb.Find(0, "golden_trinket")
	tb.G.AddNamedCounters(tb.Find(0, "bum_bo"), gold, 1)
	base := tb.G.PlayerATK(0)

	tb.G.AddNamedCounters(trinket, "", 10)
	if got := tb.G.PlayerATK(0); got != base+1 {
		t.Errorf("ATK with 10 levels on the trinket = %d, want %d: the static reads the trinket's own level", got, base+1)
	}
	tb.G.AddNamedCounters(trinket, "", -10)
	if got := tb.G.PlayerATK(0); got != base {
		t.Errorf("ATK with the trinket at 0 levels = %d, want %d: gold on Bum-Bo! alone gives no static", got, base)
	}
}

// Gaining ¢ levels Bum-Bo! itself, even with gold on it; the trinket,
// which only has the gold counter's abilities, does not level.
func TestGainCentsWithGoldOnBumBo(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"golden_trinket", "bum_bo"}, Hand: []engine.CardRef{"a_dime"}}, {Character: "cain", Cents: 5}},
	})
	tb.G.AddNamedCounters(tb.Find(0, "bum_bo"), gold, 1)
	tb.Play(0, "a_dime")
	if got := tb.G.Object(tb.Find(0, "bum_bo")).CountersOf(""); got != 10 {
		t.Errorf("Bum-Bo! level = %d, want 10", got)
	}
	if got := tb.G.Object(tb.Find(0, "golden_trinket")).CountersOf(""); got != 0 {
		t.Errorf("trinket level = %d, want 0", got)
	}
	if got := tb.G.Players[0].Cents; got != 0 {
		t.Errorf("Isaac's cents = %d, want 0: the dime became levels", got)
	}
}

// Open question (owner): Golden Trinket has Bum-Bo!'s gain-¢ replacement
// as an ability. Engine: replacements come only from items' own cards,
// so the trinket does not take the gain.
func TestGainCentsReplacementGained(t *testing.T) {
	t.Skip("open question Q-27: does the trinket gain Bum-Bo!'s gain-¢ replacement? Engine: no")
}
