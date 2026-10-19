package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Pageant Boy (One-Use Treasure Card)
//
//	{Paid Effect}Destroy this: Roll-
//	Gain x¢, where x is 4 times the result.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var pageantBoy = engine.CardDef{
	Ref:    "pageant_boy",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
