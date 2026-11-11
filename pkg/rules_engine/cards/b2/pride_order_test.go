package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// TestPrideResolvesBeforeTheAttackRoll: "when an attack is declared on
// this" resolves before the first attack roll is put on the stack and
// thrown (R-ATK-25, 19.151).
func TestPrideResolvesBeforeTheAttackRoll(t *testing.T) {
	tb := slayTable(t, "pride", seat("isaac"), engine.SituationPlayer{Character: "cain", Hand: items("a_dime", "bomb")})
	tb.G.ForceRolls(6, 6, 6)
	ev := tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "pride", foe, "a_dime", "bomb")
	pride, roll := enginetest.ResolvedAt(ev, "pride"), enginetest.FirstAttackRoll(ev)
	if pride < 0 || roll < 0 {
		t.Fatalf("Pride resolved at %d, the attack roll added at %d; both must happen", pride, roll)
	}
	if roll < pride {
		t.Errorf("the attack roll goes on the stack (event %d) before Pride resolves (event %d)", roll, pride)
	}
	if h := len(tb.G.Players[1].Hand); h != 0 {
		t.Errorf("Cain's hand = %d, want 0", h)
	}
}

// TestAttackWithNoTargetTriggerRollsAtOnce: with no trigger on the target
// pick, the first attack roll waits for its throw right after the pick.
func TestAttackWithNoTargetTriggerRollsAtOnce(t *testing.T) {
	tb := monsterTable(t, items("fly"), seat("isaac"), seat("cain"))
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentAttack}); err != nil {
		t.Fatal(err)
	}
	for tb.G.Prompt().Kind == engine.PromptPriority {
		tb.Pass(int(tb.G.Prompt().Player))
	}
	tb.Choose("fly")
	if w := tb.G.Prompt(); w.Kind != engine.PromptThrow || w.Player != 0 {
		t.Errorf("after the target pick the prompt is %+v, want Isaac's throw", w)
	}
}
