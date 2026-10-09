package b2

import "testing"

func TestShopUpgrade(t *testing.T) {
	tb := monsterTable(t, items("fly"), seat("isaac"), seat("cain"))
	shop := len(tb.G.Shop)
	tb.RevealFromDeck("shop_upgrade")
	if n := len(tb.G.Shop); n != shop+2 {
		t.Errorf("%d shop slots, want %d", n, shop+2)
	}
	if n := tb.G.Turn.Attacks; n != 1 {
		t.Errorf("attacks left %d, want 1 (the additional one)", n)
	}
}
