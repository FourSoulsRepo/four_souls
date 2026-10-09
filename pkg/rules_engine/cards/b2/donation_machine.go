package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Donation Machine (Paid Treasure Card)
//
//	{Paid Effect}Give another non-eternal item you control to another player:
//	Gain 8¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var donationMachine = engine.CardDef{
	Ref:    "donation_machine",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
