package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestOuijaBoard(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"ouija_board"}}, seat("cain")},
		Monsters: []engine.CardRef{"monstro", "fatty"},
	})
	tb.G.PutOnDeck(engine.MonsterDeck, "fly", "shop_upgrade", "leech", "trite") // trite on top
	shop := len(tb.G.Shop)
	tb.Activate(0, "ouija_board", 0, "shop_upgrade", "slot 2")
	if n := len(tb.G.Shop); n != shop+2 {
		t.Errorf("shop slots = %d, want %d: Shop Upgrade should have entered play and resolved", n, shop+2)
	}
	top := tb.G.DeckTop(engine.MonsterDeck, 1)
	if len(top) == 1 && tb.G.Object(top[0]).Card == "trite" {
		t.Error("the revealed cards should go to the bottom")
	}
}
