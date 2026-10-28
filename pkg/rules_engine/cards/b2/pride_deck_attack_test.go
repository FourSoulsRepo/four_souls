package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// TestPrideDeckAttackTriggersBeforeTheRoll: Pride revealed from the deck
// has "when an attack is declared on this" on the stack before the attack
// roll is thrown (19.105).
func TestPrideDeckAttackTriggersBeforeTheRoll(t *testing.T) {
	tb := monsterTable(t, []engine.CardRef{"fly"}, seat("isaac"), seat("cain"))
	tb.G.PutOnDeck(engine.MonsterDeck, "pride")
	tb.G.ForceRolls(6)
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentAttack}); err != nil {
		t.Fatal(err)
	}
	for tb.G.Prompt().Kind == engine.PromptPriority {
		tb.Pass(int(tb.G.Prompt().Player))
	}
	tb.Choose("ask.deck.monster")
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentChoose, Choice: 0}); err != nil {
		t.Fatal(err)
	}
	if tb.G.Prompt().Kind == engine.PromptThrow && len(tb.G.PendingTriggers) > 0 {
		t.Fatal("the attack roll waits for its throw while Pride's ability is still pending; the ability goes on the stack first (19.105)")
	}
}
