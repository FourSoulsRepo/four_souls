package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestDeadBird(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		seat("isaac", "dead_bird"), {Character: "cain", Items: items("mystery_sack"), Hand: items("bomb")},
	}, Active: 1}, Set)
	tb.G.ForceRolls(3)
	tb.Activate(1, "mystery_sack", 0, "bomb")
	if h := handCards(tb.G, 0); len(h) != 1 || h[0] != "bomb" {
		t.Errorf("hand %v, want the stolen bomb", h)
	}
}

// The thief and the victim see the stolen card; the third player only
// sees that a card changed hands (V-14).
func TestDeadBirdShowsTheCardToTwo(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		seat("isaac", "dead_bird"), {Character: "cain", Items: items("mystery_sack"), Hand: items("bomb")}, seat("judas"),
	}, Active: 1}, Set)
	tb.G.ForceRolls(3)
	ev := tb.Activate(1, "mystery_sack", 0, "bomb")
	var saw []engine.PlayerID
	for p := range engine.PlayerID(3) {
		for _, e := range engine.FilterEvents(ev, engine.Viewer{Kind: engine.ViewPlayer, Player: p}) {
			if e.Kind == engine.EvGaveCard && e.Card == "bomb" {
				saw = append(saw, p)
			}
		}
	}
	if len(saw) != 2 || saw[0] != 0 || saw[1] != 1 {
		t.Errorf("players %v saw the stolen bomb, want [0 1]", saw)
	}
}
