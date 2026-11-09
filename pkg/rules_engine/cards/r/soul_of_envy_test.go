package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSoulOfEnvy(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:    []engine.SituationPlayer{seat("isaac"), seat("cain"), seat("samson")},
		Monsters:   []engine.CardRef{"chub", "leech"},
		BonusSouls: []engine.CardRef{"soul_of_envy"},
	})
	for range 2 {
		tb.G.Players[0].InPlay = append(tb.G.Players[0].InPlay, soulOf(tb, 0, "chub"))
	}
	tb.G.Players[2].InPlay = append(tb.G.Players[2].InPlay, soulOf(tb, 2, "chub"))
	hurt(tb, 0, 1)
	tb.G.ForceRolls(6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "chub", foe) // Isaac's 3rd soul; Cain has the fewest
	settleAny(tb)
	if s := tb.G.SoulValue(1); s != 1 {
		t.Errorf("Cain's soul value = %d, want 1", s)
	}
}

// The player who gains Soul of Envy reaches 3 souls with it; the gained
// soul's own ability has stopped, so it asks only once (R-CARD-23,
// 19.144).
func TestSoulOfEnvyOnce(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:    []engine.SituationPlayer{seat("isaac"), seat("cain")},
		Monsters:   []engine.CardRef{"chub", "leech"},
		BonusSouls: []engine.CardRef{"soul_of_envy"},
	})
	for p := range 2 {
		for range 2 {
			tb.G.Players[p].InPlay = append(tb.G.Players[p].InPlay, soulOf(tb, p, "chub"))
		}
	}
	hurt(tb, 0, 1)
	tb.G.ForceRolls(6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "chub", foe) // Isaac's 3rd soul; Cain has the fewest
	// Do answered the first question.
	asks := 1
	for range 200 {
		tb.ThrowPending()
		w := tb.G.Prompt()
		if w.Kind == engine.PromptChoose {
			if w.Text.Key == askWhoGainsThe {
				asks++
			}
			tb.Choose(tb.Options()[0])
			continue
		}
		if len(tb.G.Stack) == 0 && !tb.G.Attack.On {
			break
		}
		tb.Pass(int(w.Player))
	}
	if asks != 1 {
		t.Errorf("Soul of Envy asked %d times, want 1", asks)
	}
	if s := tb.G.SoulValue(1); s != 3 {
		t.Errorf("Cain's soul value = %d, want 3", s)
	}
}
