package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Charged Penny (1 Cent Card)
//
//	Gain 1¢, then recharge an item.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var chargedPenny = engine.CardDef{
	Ref:    "charged_penny",
	Kind:   engine.LootCard,
	Copies: 1,
}
