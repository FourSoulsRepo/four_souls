package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Goat Head (Passive Treasure Card)
//
//	At the end of your turn, you may discard any number of loot cards, then loot equal to the number of cards discarded in this way.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var goatHead = engine.CardDef{
	Ref:    "goat_head",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
