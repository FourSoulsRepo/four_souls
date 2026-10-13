package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Hourglass (Passive Treasure Card)
//
//	Each time a player rolls a ❷, you may deactivate an item.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var hourglass = engine.CardDef{
	Ref:    "hourglass",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
