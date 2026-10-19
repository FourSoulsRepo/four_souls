package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestRagMega(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{
			seat("isaac"), {Character: "cain", Hand: []engine.CardRef{"bomb"}}, seat("samson"),
		},
		Monsters: []engine.CardRef{"rag_mega"},
	})
	hurt(tb, 0, 1)
	tb.G.ForceRolls(6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "rag_mega", "bomb", "player 3 (samson)", "none")
	if handHas(tb.G, 1, "bomb") {
		t.Error("Cain discards when the attack is declared")
	}
	if !tb.G.Players[2].Dead || tb.G.Players[1].Dead {
		t.Error("Isaac kills Samson and stops")
	}
}
