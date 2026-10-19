package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestMagicMarker(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"sack_of_pennies"}, Hand: []engine.CardRef{"magic_marker"}}, seat("cain")},
	})
	tb.Play(0, "magic_marker", "sack_of_pennies", "gain 1¢", "+1")
	tb.Activate(0, "sack_of_pennies", 0)
	if got := tb.G.Players[0].Cents; got != 2 {
		t.Errorf("cents = %d, want 2: the sack gains 2¢ this turn", got)
	}
}

// A roll result moves: Shell Game's "5: Gain 5¢" becomes "4: Gain 5¢".
func TestMagicMarkerRollResult(t *testing.T) {
	tb := lootRoomTable(t, "shell_game", "magic_marker")
	tb.Play(0, "magic_marker", "shell_game", "roll 5", "-1")
	tb.G.Players[0].Cents = 3
	tb.G.ForceRolls(4)
	useRoom(tb, 0, 0)
	if got := tb.G.Players[0].Cents; got != 5 {
		t.Errorf("cents = %d, want 5", got)
	}
}
