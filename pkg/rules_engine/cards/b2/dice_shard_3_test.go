package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestDiceShard(t *testing.T) {
	tb := rollTable(t, 2, engine.SituationPlayer{Character: "isaac", Hand: []engine.CardRef{"dice_shard_3"}}, seat("cain"))
	tb.Pass(1)
	tb.G.ForceRolls(6)
	if got := resolvedRoll(t, tb.Play(0, "dice_shard_3", "roll of 2")); got != 6 {
		t.Errorf("resolved as %d, want 6", got)
	}
}
