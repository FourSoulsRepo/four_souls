package b2

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestRingOfFlies(t *testing.T) {
	tb := slayTable(t, "ring_of_flies", seat("isaac"), engine.SituationPlayer{Character: "cain", Hand: items("bomb")})
	tb.G.ForceRolls(3, 6, 6, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "ring_of_flies", foe)
	if h := handCards(tb.G, 0); !slices.Contains(h, "bomb") {
		t.Errorf("hand %v, want the stolen bomb", h)
	}
}
