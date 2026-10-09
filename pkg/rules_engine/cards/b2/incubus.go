package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Incubus (Eternal Treasure Card)
//
//	{Tap Effect}Choose one- Look at a player's hand. You may swap a card from your hand with one of theirs. Loot 1, then put a card from your hand on top of the loot deck.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var incubus = engine.CardDef{
	Ref:     "incubus",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
}
