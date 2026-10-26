package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestAngelicIntervention(t *testing.T) {
	tb := roomTable(t, "angelic_intervention")
	tb.G.PutOnDeck(engine.MonsterDeck, "leech")
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "ask.deck.monster", "ask.slot_card 1 fly", "ask.yes")
	if tb.G.Attack.On || tb.G.Turn.AttackRolls != 0 {
		t.Error("Isaac may cancel the attack after the revealed monster takes a slot")
	}
}
