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

// TestDiceOptionKey: a dice value option has its own key, apart from an
// amount. The client draws it as a die, an amount as a button (18.27).
func TestDiceOptionKey(t *testing.T) {
	if got := DiceOption(6).Key; got != AskDice {
		t.Errorf("DiceOption(6) key = %q, want ask.dice", got)
	}
	if got := NumberOption(2).Key; got != AskNumber {
		t.Errorf("NumberOption(2) key = %q, want ask.number", got)
	}
	if got := (&Game{}).Plain(DiceOption(6)); got != "6" {
		t.Errorf("Plain(DiceOption(6)) = %q, want 6", got)
	}
}
