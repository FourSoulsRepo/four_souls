package rulesengine

// PurchaseState tracks a declared purchase (R-SHOP).
type PurchaseState struct {
	On bool `json:"on,omitempty"`
}

// basePrice is the default price of a shop item or the top of the
// treasure deck (R-SHOP-03).
const basePrice = 10

// Events about purchasing.
const (
	EvPurchaseDeclared EventKind = "purchase_declared"
	EvPurchased        EventKind = "purchased"
	EvPurchaseFailed   EventKind = "purchase_failed"
	EvPaid             EventKind = "paid"
)

// declarePurchase: priority passes before the item is chosen (R-SHOP-02).
func (g *Game) declarePurchase(p PlayerID) {
	if g.Turn.Purchases > 0 {
		g.Turn.Purchases--
	} else {
		g.Turn.BonusPurchasesUsed++
	}
	g.Purchase = PurchaseState{On: true}
	g.emit(Event{Kind: EvPurchaseDeclared, Player: p})
	g.openWindow(p)
}

// askPurchase offers the shop items and the top of the treasure deck.
func (g *Game) askPurchase() {
	var items []ObjectID
	for _, s := range g.Shop {
		if top, ok := s.TopOf(); ok {
			items = append(items, top)
		}
	}
	labels := append(g.labels(items), "treasure deck")
	g.ask(Choice{Purpose: ChoosePurchase, Player: g.Turn.Active, Rule: "R-SHOP-02", Objects: items, Deck: true}, labels)
}

// purchase pays the price and gains the choice, or fails (R-SHOP-03, R-SHOP-04).
// fromDeck is true for the top of the treasure deck.
func (g *Game) purchase(item ObjectID, fromDeck bool) {
	p := g.Turn.Active
	g.Purchase = PurchaseState{}
	price := basePrice
	if !fromDeck {
		price = max(price+g.bonus(StatShopPrice, p, item), 0) // e.g. Steamy Sale
	}
	if g.Players[p].Cents < price {
		g.emit(Event{Kind: EvPurchaseFailed, Player: p, Amount: price})
		return
	}
	g.Players[p].Cents -= price // pay into the game's pool (R-MECH-44)
	g.emit(Event{Kind: EvPaid, Player: p, Amount: price})
	if fromDeck {
		g.gainTreasure(p, 1)
		return
	}
	g.removeFromSlot(item)
	nid := g.move(item, Zone{Kind: ZoneInPlay}, p)
	o := g.Object(nid)
	g.enterAsItem(p, nid)
	g.emit(Event{Kind: EvPurchased, Player: p, Object: nid, Card: o.Card})
	g.enqueue(Action{Kind: ActRefillSlots, Player: NoPlayer}) // R-SHOP-06
}

// purchasesLeft is how many purchases the active player may still declare.
func (g *Game) purchasesLeft() int {
	return g.Turn.Purchases + max(g.bonus(StatPurchases, g.Turn.Active, 0)-g.Turn.BonusPurchasesUsed, 0)
}
