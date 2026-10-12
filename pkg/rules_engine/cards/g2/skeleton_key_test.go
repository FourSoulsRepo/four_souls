package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSkeletonKey(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"skeleton_key"}}, seat("cain")},
		Monsters: []engine.CardRef{"fly", "trite"},
	})
	tb.Attack("fly", 6)
	tb.Attack("trite", 6)
	// A third attack is still allowed.
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentAttack}); err != nil {
		t.Errorf("third attack refused: %v", err)
	}
}
