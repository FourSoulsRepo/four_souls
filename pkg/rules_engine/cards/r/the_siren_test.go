package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTheSirenAttackAPlayer(t *testing.T) {
	tb := sirenTable(t, nil, nil)
	tb.G.ForceRolls(4, 4)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, foe)
	if !tb.G.Players[1].Dead {
		t.Fatal("two hits on DC 4 kill Cain (2 HP)")
	}
	if tb.G.Attack.On {
		t.Error("killing the attacked player ends the attack")
	}
}

func TestTheSirenMissHurtsTheAttacker(t *testing.T) {
	tb := sirenTable(t, nil, nil)
	tb.G.ForceRolls(3, 6, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, foe)
	if tb.G.Players[0].Damage != 1 {
		t.Errorf("a miss: Cain deals his 1 ATK to Isaac; damage %d", tb.G.Players[0].Damage)
	}
}
