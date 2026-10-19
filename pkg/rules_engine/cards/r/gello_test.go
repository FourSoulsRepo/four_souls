package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestGello(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), seat("the_harlot", "gello")},
		Monsters: []engine.CardRef{"pooter", "fly"},
	})
	pooter := tb.G.Monsters[0].Cards[0]
	tb.G.Object(pooter).Damage = 1
	tb.G.ForceRolls(1)
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentAttack}); err != nil {
		t.Fatal(err)
	}
	tb.Pass(0)
	tb.Pass(1)
	tb.Choose("pooter")
	tb.ThrowPending() // the attack roll waits on the stack
	tb.Pass(0)
	tb.Activate(1, "gello", 0, "pooter", "fly")
	if tb.G.Object(pooter).Damage != 0 || top(tb, 1) == "fly" {
		t.Error("the pooter heals 1; 1 damage kills the fly")
	}
	if tb.G.Attack.On || tb.G.Turn.Attacks != 1 {
		t.Error("not her turn: the attack is cancelled, Isaac may attack again")
	}
}
