package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Devil Deal (Good Event Card)
//
//	Choose one- Put this into discard. Loot 2. Take 1 damage. Take 2 damage. Search the treasure deck for a guppy item, gain it, then shuffle the treasure deck.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var devilDeal = engine.CardDef{
	Ref:    "devil_deal",
	Kind:   engine.EventCard,
	Copies: 1,
}
