package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/b2"
	"github.com/FourSoulsRepo/rules_engine/cards/g2"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// The copy of Lord of the Pit waits on the stack while Cain's own Lord
// of the Pit cancels the attack first (19.238). Go, not a file: Cain
// responds in the middle of the stack, which steps don't do (24 owner
// answer 7), and ? Card can't start in setup.stack.
func TestQuestionmarkCardCopyFizzlesWhenTheAttackEnds(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{
			{Character: "isaac", Hand: []engine.CardRef{"questionmark_card"}},
			{Character: "cain", Items: []engine.CardRef{"lord_of_the_pit"}},
		},
		Monsters: []engine.CardRef{"fly", "leech"},
		Decks:    map[string][]engine.CardRef{"monster": {"dip"}},
		Stack:    []engine.SituationStack{{Kind: "attack_roll", DiceRoll: 6, Card: "fly", Owner: 0}},
	}, b2.Set, Set, g2.Set)
	card := tb.Find(0, "questionmark_card")
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentPlayLoot, Objects: []engine.ObjectID{card}}); err != nil {
		t.Fatal(err)
	}
	tb.Choose("lord_of_the_pit: ↷: Cancel an attack. That player may attack an additional time this turn.")
	tb.Pass(0)
	if w := tb.G.Prompt(); w.Kind != engine.PromptPriority || w.Player != 1 {
		t.Fatalf("Cain should hold priority over the ? Card, prompt %+v", w)
	}
	tb.Start(1, "lord_of_the_pit", 0)
	tb.Settle()

	// Cain's Lord of the Pit resolves first and cancels the attack: one
	// extra attack for Isaac, "that player" (card lord_of_the_pit). The
	// copy then has no attack to cancel and fizzles: no second extra
	// attack (R-ABIL-06, R-MECH-32).
	if got := tb.G.Turn.Attacks; got != 1 {
		t.Errorf("attacks left = %d, want 1: only Cain's Lord of the Pit gives one", got)
	}
	// The attack roll of 6 left the stack with the attack (R-ATK-15).
	if tb.G.Attack.On {
		t.Error("the attack ended")
	}
	if fly, ok := tb.G.Monsters[0].TopOf(); !ok || tb.G.Object(fly).Card != "fly" || tb.G.Object(fly).Damage != 0 {
		t.Error("the Fly is still in slot 1 with no damage (R-ATK-15)")
	}
	// A fizzle cancels the loot; it goes to the loot discard (R-MECH-32,
	// R-STACK-06).
	discard := tb.G.Discards[engine.LootDeck]
	if len(discard) != 1 || tb.G.Object(discard[0]).Card != "questionmark_card" {
		t.Errorf("the loot discard holds the ? Card, got %d cards", len(discard))
	}
}
