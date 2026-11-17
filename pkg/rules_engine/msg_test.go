package rulesengine

import "testing"

// TestPileOptionKeys: a deck or discard pile option names its kind in
// its key, "ask.deck.<kind>" or "ask.discard.<kind>". The client reads
// the kind from it to show the pile, not a button (20.26).
func TestPileOptionKeys(t *testing.T) {
	kinds := map[DeckKind]string{TreasureDeck: "treasure", LootDeck: "loot", MonsterDeck: "monster", RoomDeck: "room"}
	for d, name := range kinds {
		if got := DeckOption(d).Key.String(); got != "ask.deck."+name {
			t.Errorf("DeckOption(%s) = %q", name, got)
		}
		if got := DiscardOption(d).Key.String(); got != "ask.discard."+name {
			t.Errorf("DiscardOption(%s) = %q", name, got)
		}
		if len(DeckOption(d).Args) != 0 || len(DiscardOption(d).Args) != 0 {
			t.Errorf("a %s pile option has values", name)
		}
	}
}
