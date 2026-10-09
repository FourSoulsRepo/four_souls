package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestGurdyJr(t *testing.T) {
	tb := slayTable(t, "gurdy_jr", engine.SituationPlayer{Character: "isaac", Items: items("mr_boom")}, seat("cain"))
	tb.G.ForceRolls(1)
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentAttack}); err != nil {
		t.Fatal(err)
	}
	tb.Pass(0)
	tb.Pass(1)
	tb.Choose("gurdy_jr")
	tb.Activate(0, "mr_boom", 0, "gurdy_jr", "mr_boom") // 2 damage kills Isaac: the penalty
	if d := tb.G.Players[0].Damage; d < 1 {
		t.Error("activating an item during the attack did not hurt")
	}
}
