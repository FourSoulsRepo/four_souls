package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestMegaFatty(t *testing.T) {
	tb := slayTable(t, "mega_fatty", engine.SituationPlayer{Character: "isaac", Items: items("breakfast")}, seat("cain"))
	m, _ := tb.G.Monsters[0].TopOf()
	tb.G.Object(m).Damage = 1
	tb.G.ForceRolls(1, 6, 6, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "mega_fatty")
	if tb.G.Object(m).Zone.Kind == engine.ZoneInPlay {
		t.Fatal("it survived")
	}
	if tb.G.Turn.AttackRolls != 4 {
		t.Errorf("%d attack rolls, want 4: it healed 1", tb.G.Turn.AttackRolls)
	}
}
