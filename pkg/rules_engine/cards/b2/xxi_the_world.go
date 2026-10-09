package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// XXI. The World (Wildcard Card)
//
//	Look at each player's hand, then loot 2.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var xxiTheWorld = engine.CardDef{
	Ref:    "xxi_the_world",
	Kind:   engine.LootCard,
	Copies: 1,
}
