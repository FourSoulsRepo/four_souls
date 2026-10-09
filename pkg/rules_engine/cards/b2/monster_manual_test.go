package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestMonsterManual(t *testing.T) {
	tb := itemTable(t, items("monster_manual"))
	tb.Activate(0, "monster_manual", 0, "leech")
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentEndTurn}); err == nil {
		t.Fatal("ended the turn without attacking the leech")
	}
	leech, _ := tb.G.Monsters[1].TopOf()
	tb.G.ForceRolls(6)
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentAttack}); err != nil {
		t.Fatal(err)
	}
	tb.Pass(0)
	tb.Pass(1)
	if opts := tb.Options(); len(opts) != 1 || opts[0] != "leech" {
		t.Fatalf("attack options %v, want only the leech", opts)
	}
	tb.Settle("leech")
	if tb.G.Object(leech).Zone.Kind == engine.ZoneInPlay {
		t.Error("the leech survived")
	}
}
