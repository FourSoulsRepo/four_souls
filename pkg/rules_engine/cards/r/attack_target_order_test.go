package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// Triggers of the target pick resolve before the first attack roll is
// put on the stack (R-ATK-25, 19.151).

func TestRagMegaResolvesBeforeTheAttackRoll(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{
			seat("isaac"), {Character: "cain", Hand: []engine.CardRef{"bomb"}}, seat("samson"),
		},
		Monsters: []engine.CardRef{"rag_mega"},
	})
	hurt(tb, 0, 1)
	tb.G.ForceRolls(6)
	ev := tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "rag_mega", "bomb", "player 3 (samson)", "ask.none")
	rag, roll := enginetest.ResolvedAt(ev, "rag_mega"), enginetest.FirstAttackRoll(ev)
	if rag < 0 || roll < 0 {
		t.Fatalf("Rag Mega resolved at %d, the attack roll added at %d; both must happen", rag, roll)
	}
	if roll < rag {
		t.Errorf("the attack roll goes on the stack (event %d) before Rag Mega's discard resolves (event %d)", roll, rag)
	}
}

func TestDogmaResolvesBeforeTheAttackRoll(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), seat("cain")},
		Monsters: []engine.CardRef{"dogma", "fly"},
	})
	hurt(tb, 0, 1)
	tb.G.ForceRolls(6)
	ev := tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "dogma", "ask.slot 1")
	dogma, roll := enginetest.ResolvedAt(ev, "dogma"), enginetest.FirstAttackRoll(ev)
	if dogma < 0 || roll < 0 {
		t.Fatalf("Dogma resolved at %d, the attack roll added at %d; both must happen", dogma, roll)
	}
	if roll < dogma {
		t.Errorf("the attack roll goes on the stack (event %d) before Dogma's ability resolves (event %d)", roll, dogma)
	}
}

// TestAngelicInterventionCancelsBeforeAnyRoll: a deck attack cancelled by
// Angelic Intervention puts no attack roll on the stack.
func TestAngelicInterventionCancelsBeforeAnyRoll(t *testing.T) {
	tb := roomTable(t, "angelic_intervention")
	tb.G.PutOnDeck(engine.MonsterDeck, "leech")
	ev := tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "ask.deck.monster", "ask.slot_card 1 fly", "ask.yes")
	if i := enginetest.FirstAttackRoll(ev); i >= 0 {
		t.Errorf("an attack roll went on the stack (event %d) before the attack was cancelled", i)
	}
}
