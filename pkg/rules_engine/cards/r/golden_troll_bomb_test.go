package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestGoldenTrollBomb(t *testing.T) {
	tb := fight(t, "fly")
	tb.G.PutOnDeck(engine.MonsterDeck, "golden_troll_bomb")
	tb.G.ForceRolls(6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "fly", foe)
	if !tb.G.Players[1].Dead {
		t.Error("2 damage kills Cain")
	}
	deck := tb.G.Decks[engine.MonsterDeck]
	// 4 cards from the top; the refill of its slot then drew one of them.
	if tb.G.Object(deck[len(deck)-4]).Card != "golden_troll_bomb" {
		t.Error("it goes 4 cards from the top")
	}
}
