package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// XII. The Hanged Man (Wildcard Card)
//
//	Look at the top card of each deck. You may put any of those cards on the bottom of their deck, then loot 2.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var xiiTheHangedMan = engine.CardDef{
	Ref:    "xii_the_hanged_man",
	Kind:   engine.LootCard,
	Copies: 1,
}
