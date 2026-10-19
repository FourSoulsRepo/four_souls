package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestConjoinedTwin(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), seat("cain"), seat("samson")},
		Monsters: []engine.CardRef{"fly", "leech"},
		Rooms:    []engine.CardRef{"conjoined_twin"},
	})
	tb.G.ForceRolls(6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "fly", "player 3 (samson)")
	if !tb.G.Players[2].Dead || tb.G.Players[1].Dead {
		t.Error("Isaac killed the fly, so he kills another player: Samson")
	}
}
