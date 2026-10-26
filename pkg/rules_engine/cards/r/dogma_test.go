package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestDogma(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), seat("cain")},
		Monsters: []engine.CardRef{"dogma", "fly"},
	})
	hurt(tb, 0, 1)
	tb.G.ForceRolls(6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "dogma", "ask.slot 1")
	if !inDiscard(tb.G, engine.MonsterDeck, "fly") {
		t.Error("the other monsters go to discard when Dogma is attacked")
	}
	found := false
	for i := range tb.G.Monsters {
		if len(tb.G.Monsters[i].Cards) > 0 && top(tb, i) == "the_harbingers" {
			found = true
		}
	}
	if !found {
		t.Error("The Harbingers comes into play when Dogma dies")
	}
}
