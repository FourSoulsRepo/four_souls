package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// XIII. Death (Wildcard Card)
//
//	Kill a player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var xiiiDeath = engine.CardDef{
	Ref:    "xiii_death",
	Kind:   engine.LootCard,
	Copies: 1,
}
