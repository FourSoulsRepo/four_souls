package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// Eve's line (owner, 2026-10-29): gold on Tech X, Eve dies and loses the
// trinket, the trinket comes back to her hand, then gold on Bum-Bo!.
func TestEveGoldenTrinketLine(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "eve", Items: []engine.CardRef{"tech_x", "bum_bo"}, Hand: []engine.CardRef{"golden_trinket"}},
		{Character: "isaac", Cents: 5},
	}})
	tech, bum := tb.Find(0, "tech_x"), tb.Find(0, "bum_bo")

	// 1-2: Golden Trinket puts a gold counter on Tech X.
	tb.Play(0, "golden_trinket", "tech_x")
	if got := tb.G.Object(tech).CountersOf(gold); got != 1 {
		t.Fatalf("Tech X gold = %d, want 1", got)
	}

	// 3: Eve dies and loses the trinket; Tech X and Bum-Bo! stay in play.
	tb.G.KillPlayer(0, 1)
	tb.Settle("golden_trinket")
	if !tb.G.Players[0].Dead {
		t.Fatal("Eve should be dead")
	}
	if !inDiscard(tb.G, engine.LootDeck, "golden_trinket") {
		t.Fatal("the penalty item is not the trinket in the loot discard")
	}
	if got := tb.G.Object(tech).CountersOf(gold); got != 1 {
		t.Errorf("Tech X gold = %d after the death, want 1: a gold counter stays on a card that stays in play", got)
	}

	// 4: Eve gets the trinket back from the loot discard, as The Curse
	// would: it is a new card with no counters.
	discarded := tb.G.Discards[engine.LootDeck]
	tb.G.DiscardToHand(engine.LootDeck, discarded[len(discarded)-1], 0)
	if got := tb.G.Object(tb.Find(0, "golden_trinket")).CountersOf(""); got != 0 {
		t.Fatalf("the trinket came back with %d counters, want 0", got)
	}

	// Eve is back on her next turn, with her loot play.
	// Her death ends her turn once everyone has passed (R-DEATH-16).
	for range 100 {
		if tb.G.Turn.Active == 1 && tb.G.Turn.Step == engine.StepAction {
			break
		}
		tb.Pass(int(tb.G.Prompt().Player))
	}
	tb.EndTurn() // Isaac's turn; then Eve is back
	if tb.G.Players[0].Dead {
		t.Fatal("Eve should be alive on her next turn")
	}

	// 5: Golden Trinket puts a gold counter on Bum-Bo!.
	tb.Play(0, "golden_trinket", "bum_bo")
	if got := tb.G.Object(bum).CountersOf(gold); got != 1 {
		t.Errorf("Bum-Bo! gold = %d, want 1", got)
	}
	if got := tb.G.Object(bum).CountersOf(""); got != 0 {
		t.Errorf("Bum-Bo! level = %d, want 0: gold is not a level", got)
	}
	if got := tb.G.Object(tech).CountersOf(gold); got != 1 {
		t.Errorf("Tech X gold = %d, want 1: it does not move", got)
	}
}
