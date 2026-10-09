package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mom’s Box (Passive Treasure Card)
//
//	Each time a player rolls a ❹, you may loot 1, then discard a loot card.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var momsBox = engine.CardDef{
	Ref:    "moms_box",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
