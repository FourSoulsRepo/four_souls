package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBlessingOfTheInnerEye(t *testing.T) {
	tb := roomTable(t, "blessing_of_the_inner_eye")
	tb.G.PutOnDeck(engine.MonsterDeck, "pooter")
	tb.G.PutOnDeck(engine.MonsterDeck, "leech")
	tb.G.ForceRolls(6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "pooter", "leech", "clotty", "ask.deck.monster", "ask.slot_card 1 fly")
	if !inDiscard(tb.G, engine.MonsterDeck, "pooter") {
		t.Error("Isaac put the pooter on top, revealed it and killed it")
	}
}
