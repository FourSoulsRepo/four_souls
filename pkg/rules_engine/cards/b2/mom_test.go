package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestMom(t *testing.T) {
	tb := slayTable(t, "mom", engine.SituationPlayer{Character: "isaac", Items: items("breakfast", "dinner", "meat")}, seat("cain"))
	m, _ := tb.G.Monsters[0].TopOf()
	tb.G.Object(m).Damage = 4
	tb.G.ForceRolls(6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "mom")
	if n := len(tb.G.Monsters); n != 2 {
		t.Errorf("%d monster slots, want 2", n)
	}
}
