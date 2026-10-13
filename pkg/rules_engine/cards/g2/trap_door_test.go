package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTrapDoor(t *testing.T) {
	tb := fight(t, "fatty", "monstro")
	tb.G.ForceRolls(5)
	tb.RevealFromDeck("trap_door")
	if d := tb.G.Players[0].Damage; d != 1 {
		t.Errorf("4-6: damage = %d, want 1", d)
	}
	tb2 := fight(t, "fatty", "monstro")
	tb2.G.PutOnDeck(engine.MonsterDeck, "fly", "leech", "trite") // trite on top
	tb2.G.ForceRolls(2)
	tb2.RevealFromDeck("trap_door", "fly", "leech", "trite")
	if top := tb2.G.DeckTop(engine.MonsterDeck, 3); tb2.G.Object(top[0]).Card != "fly" || tb2.G.Object(top[2]).Card != "trite" {
		t.Error("the top 3 were not put back in the chosen order")
	}
	if tb2.G.Turn.Attacks != 1 {
		t.Errorf("1-3: attacks left = %d, want 1 more", tb2.G.Turn.Attacks)
	}
}
