package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestMegaSatanIndomitable(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), seat("cain"), seat("samson")},
		Monsters: []engine.CardRef{"fly", "leech"},
	})
	tb.G.PutOnDeck(engine.MonsterDeck, "mega_satan")
	tb.G.ForceRolls(6)
	tb.Attack("fly", 6) // the fly dies; Mega Satan would refill its slot
	if len(tb.G.Monsters) != 3 {
		t.Fatalf("monster slots = %d, want 3: Mega Satan takes a new one", len(tb.G.Monsters))
	}
	if top(tb, 2) != "mega_satan" || top(tb, 0) == "fly" {
		t.Error("Mega Satan in the new slot; slot 1 refilled with another card")
	}
}

func TestMegaSatanWins(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), {Character: "cain", Dead: true}, {Character: "samson", Dead: true}},
		Monsters: []engine.CardRef{"mega_satan", "fly"},
	})
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentEndTurn}); err != nil {
		t.Fatal(err)
	}
	for range 20 {
		if w := tb.G.Prompt(); w.Kind == engine.PromptPriority && !tb.G.Over {
			tb.Pass(int(w.Player))
		}
	}
	if !tb.G.Over || len(tb.G.Winners) != 1 || tb.G.Winners[0] != 0 {
		t.Error("the only living player wins at the end of the turn")
	}
}
