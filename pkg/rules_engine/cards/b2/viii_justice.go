package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// VIII. Justice (Wildcard Card)
//
//	Choose a player. Loot and gain ¢ until you have the same number of each as they do.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var viiiJustice = engine.CardDef{
	Ref:    "viii_justice",
	Kind:   engine.LootCard,
	Copies: 1,
}
