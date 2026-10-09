package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// VII. The Chariot (Wildcard Card)
//
//	Choose a player.
//	They gain +1{ATK} and +1{HP} till end of turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var viiTheChariot = engine.CardDef{
	Ref:    "vii_the_chariot",
	Kind:   engine.LootCard,
	Copies: 1,
}
