package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Charged Baby (Passive Treasure Card)
//
//	Each time a player rolls a ❷, you may recharge an item.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var chargedBaby = engine.CardDef{
	Ref:    "charged_baby",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
