package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dad’s Key (Paid Treasure Card)
//
//	{Paid Effect}Discard a loot card:
//	Put the top card of the monster deck in a monster slot not being attacked.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var dadsKey = engine.CardDef{
	Ref:    "dads_key",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
