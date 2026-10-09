package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Cursed Chest (Bad Event Card)
//
//	Roll-
//	1-3: Take 1 Damage.
//	4-5: Take 2 Damage.
//	6: Search the treasure deck for a Guppy item, gain it, then shuffle the treasure deck.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var cursedChest = engine.CardDef{
	Ref:    "cursed_chest",
	Kind:   engine.EventCard,
	Copies: 1,
}
