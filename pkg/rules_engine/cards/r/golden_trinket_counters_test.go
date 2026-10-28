package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// Golden Trinket's gold counter is its own kind: general counters (the
// plain "" name) and gold counters never count for each other.

// TestCountersDontMix: a gold counter does not level Bum-Bo! and does not
// pay Tech X's "Remove 3 counters".
func TestCountersDontMix(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"tech_x", "bum_bo"}, Hand: []engine.CardRef{"golden_trinket"}}, {Character: "cain", Cents: 5}},
	})
	tech := tb.Find(0, "tech_x")
	tb.G.AddNamedCounters(tech, "", 2)

	tb.Play(0, "golden_trinket", "tech_x")
	if got := tb.G.Object(tech).CountersOf(gold); got != 1 {
		t.Fatalf("Tech X gold = %d, want 1", got)
	}
	if got := tb.G.Object(tech).CountersOf(""); got != 2 {
		t.Errorf("Tech X general = %d, want 2: the gold counter is not a general one", got)
	}
	if rule := tb.Refused(0, "tech_x", 1); rule == "" {
		t.Error("Remove 3 counters was allowed with 2 general and 1 gold counter")
	}
}

// TestGoldOnBumBoNoLevel: a gold counter on Bum-Bo! is not a level.
func TestGoldOnBumBoNoLevel(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"bum_bo"}, Hand: []engine.CardRef{"golden_trinket"}}, {Character: "cain", Cents: 5}},
	})
	tb.Play(0, "golden_trinket", "bum_bo")
	bum := tb.Find(0, "bum_bo")
	if got := tb.G.Object(bum).CountersOf(gold); got != 1 {
		t.Fatalf("Bum-Bo! gold = %d, want 1", got)
	}
	if got := tb.G.Object(bum).CountersOf(""); got != 0 {
		t.Errorf("Bum-Bo! level counters = %d, want 0: a gold counter is not a level", got)
	}
}

// TestTwoGoldCountersBothGained: with gold on Tech X and on Bum-Bo!, Golden
// Trinket has Tech X's ↷ abilities and nothing to activate from Bum-Bo!.
// The player picks which ability to use; one use does one ability.
func TestTwoGoldCountersBothGained(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"golden_trinket", "tech_x", "bum_bo"}}, {Character: "cain", Cents: 5}},
	})
	tech, bum, trinket := tb.Find(0, "tech_x"), tb.Find(0, "bum_bo"), tb.Find(0, "golden_trinket")
	tb.G.AddNamedCounters(tech, gold, 1)
	tb.G.AddNamedCounters(bum, gold, 1)

	// Golden Trinket's own ability, then Tech X's two. Bum-Bo! has no ↷
	// ability, so it adds none: one activation picks one of these.
	if got := len(tb.G.AbilitiesOf(trinket)); got != 3 {
		t.Fatalf("Golden Trinket has %d abilities, want 3: its own and Tech X's two", got)
	}

	// Ability 1 is Tech X's "Put a counter": it goes on Golden Trinket,
	// not on Tech X, and Bum-Bo! and Tech X keep their counters.
	tb.Activate(0, "golden_trinket", 1)
	if got := tb.G.Object(trinket).CountersOf(""); got != 1 {
		t.Errorf("Golden Trinket counters = %d, want 1: Tech X's ↷ puts it on the trinket", got)
	}
	if got := tb.G.Object(tech).CountersOf(""); got != 0 {
		t.Errorf("Tech X general = %d, want 0", got)
	}
	if got := tb.G.Object(bum).CountersOf(""); got != 0 {
		t.Errorf("Bum-Bo! general = %d, want 0", got)
	}
	// Ability 2 is Tech X's "Remove 3 counters": one counter is not enough.
	if rule := tb.Refused(0, "golden_trinket", 2); rule == "" {
		t.Error("Remove 3 counters was allowed with 1 counter on the trinket")
	}
}
