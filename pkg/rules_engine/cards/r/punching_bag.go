package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Punching Bag (Passive Treasure Card)
//
//	You have +2{HP}.
//	This item can be attacked.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var punchingBag = engine.CardDef{
	Ref:    "punching_bag",
	Kind:   engine.TreasureCard,
	Copies: 1,
	HP:     2,
	DC:     5,
	ATK:    1,
}
