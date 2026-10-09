package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Relic (Passive Treasure Card)
//
//	Each time a player rolls a ❶, loot 1.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theRelic = engine.CardDef{
	Ref:    "the_relic",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
