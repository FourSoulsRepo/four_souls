package rulesengine

import "testing"

// enterSet has an item that expands shop slots by 1 when it enters play.
var enterSet = CardSet{Name: "enter", Cards: []CardDef{{
	Ref: "welcome_mat", Kind: TreasureCard,
	Abilities: []Ability{{
		Kind: Triggered, Text: "When this enters play, expand shop slots by 1.", Trigger: WhenThisEntersPlay(),
		Effects: []Effect{EffectFunc(func(c *Ctx) { c.G.ExpandShop(1) })},
	}},
}}}

func enterTable(t *testing.T) *Game {
	t.Helper()
	g, err := BuildTable(SituationSetup{
		Players: []SituationPlayer{{Character: "hero_a", Cents: 10}, {Character: "hero_b"}},
		Shop:    []CardRef{"trinket"},
	}, testSet, enterSet)
	if err != nil {
		t.Fatal(err)
	}
	// Only the mat put on top: the deck's own copy would come up at random.
	var deck []ObjectID
	for _, id := range g.Decks[TreasureDeck] {
		if g.Object(id).Card != "welcome_mat" {
			deck = append(deck, id)
		}
	}
	g.Decks[TreasureDeck] = deck
	g.PutOnDeck(TreasureDeck, "welcome_mat")
	return g
}

// An item gained from the treasure deck enters play (19.25).
func TestGainedItemEntersPlay(t *testing.T) {
	g := enterTable(t)
	if _, err := g.Apply(Intent{Player: 0, Kind: IntentPurchase}); err != nil {
		t.Fatal(err)
	}
	if _, err := g.Settle("ask.deck.treasure"); err != nil {
		t.Fatal(err)
	}
	if len(g.Shop) != 2 {
		t.Errorf("shop slots = %d, want 2: the mat entered play under its buyer", len(g.Shop))
	}
}

// A shop item's "enters play" ability works in the shop (R-CARD-05).
func TestShopItemEntersPlay(t *testing.T) {
	g := enterTable(t)
	if _, err := g.Apply(Intent{Player: 0, Kind: IntentPurchase}); err != nil {
		t.Fatal(err)
	}
	if _, err := g.Settle("trinket"); err != nil { // the slot refills with the mat
		t.Fatal(err)
	}
	if len(g.Shop) != 2 {
		t.Errorf("shop slots = %d, want 2: the mat entered play in the shop", len(g.Shop))
	}
}
