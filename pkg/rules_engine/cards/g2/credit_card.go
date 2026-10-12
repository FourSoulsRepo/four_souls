package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Credit Card (Wildcard Card)
//
//	The first shop item a player purchases this turn costs 0¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var creditCard = engine.CardDef{
	Ref:    "credit_card",
	Kind:   engine.LootCard,
	Copies: 1,
}
