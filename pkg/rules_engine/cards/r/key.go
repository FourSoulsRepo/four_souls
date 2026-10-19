package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Key (Wildcard Card)
//
//	Choose one- The active player may attack an additional time this turn. Put the top card of the monster deck in a monster slot not being attacked.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var key = engine.CardDef{
	Ref:    "key",
	Kind:   engine.LootCard,
	Copies: 1,
}
