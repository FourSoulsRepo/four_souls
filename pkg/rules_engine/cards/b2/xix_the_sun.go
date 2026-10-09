package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// XIX. The Sun (Wildcard Card)
//
//	Put this on the bottom of the loot deck. If you do, take an extra turn after this one if it's your turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var xixTheSun = engine.CardDef{
	Ref:    "xix_the_sun",
	Kind:   engine.LootCard,
	Copies: 1,
}
