package rulesengine

import "testing"

func declarePurchase(t *testing.T, g *Game) {
	t.Helper()
	if _, err := g.Apply(Intent{Player: g.Turn.Active, Kind: IntentPurchase}); err != nil {
		t.Fatal(err)
	}
	passWhile(t, g, func() bool { return true })
	if p := g.Prompt(); p.Kind != PromptChoose || p.Purpose != ChoosePurchase {
		t.Fatalf("prompt %+v; want the purchase choice (R-SHOP-02)", p)
	}
}

func TestBuyAShopItem(t *testing.T) {
	g := toAction(t, 2)
	a := g.Turn.Active
	g.Players[a].Cents = 12
	shopItem, _ := g.Shop[0].TopOf()
	card := g.Object(shopItem).Card
	declarePurchase(t, g)
	if _, err := g.Apply(Intent{Player: a, Kind: IntentChoose, Choice: 0}); err != nil {
		t.Fatal(err)
	}
	if g.Players[a].Cents != 2 {
		t.Errorf("cents %d, want 12-10 (R-SHOP-03)", g.Players[a].Cents)
	}
	found := false
	for _, id := range g.Players[a].InPlay {
		if g.Object(id).Card == card && g.Object(id).Role == RoleItem {
			found = true
		}
	}
	if !found {
		t.Error("the bought item is not under the buyer's control (R-SHOP-04)")
	}
	if top, ok := g.Shop[0].TopOf(); !ok || top == shopItem {
		t.Error("the shop slot was not refilled (R-SHOP-06)")
	}
	if g.Turn.Purchases != 0 {
		t.Error("one purchase per turn (R-SHOP-05)")
	}
	if _, err := g.Apply(Intent{Player: a, Kind: IntentPurchase}); err == nil {
		t.Error("purchased twice in one turn")
	}
}

func TestBuyFromTheDeck(t *testing.T) {
	g := toAction(t, 2)
	a := g.Turn.Active
	g.Players[a].Cents = 10
	items := len(g.Players[a].InPlay)
	deck := len(g.Decks[TreasureDeck])
	declarePurchase(t, g)
	choose(t, g, "treasure deck")
	if len(g.Players[a].InPlay) != items+1 || len(g.Decks[TreasureDeck]) != deck-1 || g.Players[a].Cents != 0 {
		t.Errorf("items %d->%d, deck %d->%d, cents %d", items, len(g.Players[a].InPlay), deck, len(g.Decks[TreasureDeck]), g.Players[a].Cents)
	}
}

func TestPurchaseFailsWithoutMoney(t *testing.T) {
	g := toAction(t, 2)
	a := g.Turn.Active
	g.Players[a].Cents = 9
	items := len(g.Players[a].InPlay)
	declarePurchase(t, g)
	if _, err := g.Apply(Intent{Player: a, Kind: IntentChoose, Choice: 0}); err != nil {
		t.Fatal(err)
	}
	if g.Players[a].Cents != 9 || len(g.Players[a].InPlay) != items {
		t.Errorf("a player who cannot pay gains nothing and pays nothing (R-SHOP-04)")
	}
	if !g.inOpenActionPhase() {
		t.Error("after a failed purchase the action phase goes on")
	}
}
