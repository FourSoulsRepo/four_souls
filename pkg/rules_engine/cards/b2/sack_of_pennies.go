package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Sack Of Pennies (Active Treasure Card)
//
//	{Tap Effect}Gain 1¢.
//	Each time a player rolls a ❶, you may recharge this.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var sackOfPennies = engine.CardDef{
	Ref:    "sack_of_pennies",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
