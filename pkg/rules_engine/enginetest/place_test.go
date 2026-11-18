package enginetest

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/b2"
)

// TestPlaceFromDeck: a loot card goes from its deck into a hand, a
// treasure card into play as a charged item; the state stays sound; a
// card out of its deck, or of the wrong kind, is refused.
func TestPlaceFromDeck(t *testing.T) {
	g, _, err := engine.NewGame(engine.Setup{Seed: 3, Players: 2, Sets: []engine.CardSet{b2.Set}})
	if err != nil {
		t.Fatal(err)
	}
	loot := g.Decks[engine.LootDeck][0]
	card := g.Object(loot).Card
	hand := len(g.Players[1].Hand)
	id, err := PlaceFromDeck(g, loot, 1, false)
	if err != nil {
		t.Fatal(err)
	}
	if o := g.Object(id); o.Card != card || o.Zone.Kind != engine.ZoneHand || len(g.Players[1].Hand) != hand+1 || !slices.Contains(g.Players[1].Hand, id) {
		t.Errorf("loot: %+v, hand %v", o, g.Players[1].Hand)
	}
	if slices.Contains(g.Decks[engine.LootDeck], loot) || g.Object(loot).Zone.Kind != engine.ZoneGone {
		t.Error("the loot card is still in its deck")
	}
	treasure := g.Decks[engine.TreasureDeck][0]
	id, err = PlaceFromDeck(g, treasure, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if o := g.Object(id); o.Role != engine.RoleItem || o.Controller != 0 || !slices.Contains(g.Players[0].InPlay, id) {
		t.Errorf("item: %+v", o)
	}
	if _, err := PlaceFromDeck(g, treasure, 0, true); err == nil {
		t.Error("a card no longer in its deck was placed")
	}
	if _, err := PlaceFromDeck(g, g.Decks[engine.TreasureDeck][0], 0, false); err == nil {
		t.Error("a treasure card went into a hand")
	}
	if err := g.CheckInvariants(); err != nil {
		t.Fatal(err)
	}
}
