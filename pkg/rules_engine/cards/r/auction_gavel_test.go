package r

import (
	"testing"

	"github.com/FourSoulsRepo/rules_engine/enginetest"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func gavelTable(t *testing.T) *enginetest.Table {
	t.Helper()
	return setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{
			seat("isaac", "auction_gavel"), {Character: "cain", Cents: 5}, {Character: "samson", Cents: 5},
		},
		Shop: []engine.CardRef{"spoon_bender"},
	})
}

func TestAuctionGavel(t *testing.T) {
	tb := gavelTable(t)
	// Cain bids 1¢, Samson tops it with 3¢, Cain passes: Samson's bid stands.
	tb.Activate(0, "auction_gavel", 0, "spoon_bender", "card.auction_gavel.bid 1", "card.auction_gavel.bid 3", "card.auction_gavel.pass")
	if !hasItem(tb.G, 2, "spoon_bender") {
		t.Fatal("Samson wins the item")
	}
	if tb.G.Players[0].Cents != 3 || tb.G.Players[2].Cents != 2 || tb.G.Players[1].Cents != 5 {
		t.Errorf("cents %d %d %d; Samson pays Isaac 3¢", tb.G.Players[0].Cents, tb.G.Players[1].Cents, tb.G.Players[2].Cents)
	}
}

func TestAuctionGavelNoBids(t *testing.T) {
	tb := gavelTable(t)
	tb.Activate(0, "auction_gavel", 0, "spoon_bender", "card.auction_gavel.pass", "card.auction_gavel.pass")
	if len(tb.G.ShopItems()) != 1 || tb.G.Object(tb.G.ShopItems()[0]).Card != "spoon_bender" {
		t.Error("nobody bid: the item stays in the shop")
	}
}
