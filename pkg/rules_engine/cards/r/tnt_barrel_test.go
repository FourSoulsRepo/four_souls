package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTNTBarrel(t *testing.T) {
	tb := fight(t, "fly")
	tb.G.PutOnDeck(engine.MonsterDeck, "tnt_barrel")
	tb.G.ForceRolls(6, 5)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "fly", foe)
	if !tb.G.Players[1].Dead || tb.G.Players[0].Damage != 0 {
		t.Error("5-6: only Cain takes 4 damage")
	}
}
