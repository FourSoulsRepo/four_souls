package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Fire Mind (Passive Treasure Card)
//
//	Combat damage you deal is doubled on attack rolls of 6.
//	Combat damage you take is doubled on attack rolls of 1.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var fireMind = engine.CardDef{
	Ref:    "fire_mind",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
