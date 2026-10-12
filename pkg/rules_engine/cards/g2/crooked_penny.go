package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Crooked Penny (Active Treasure Card)
//
//	{Tap Effect}Choose a player, then roll-
//	1-3: Double the amount of ¢ they have.
//	4-6: That player loses all ¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var crookedPenny = engine.CardDef{
	Ref:    "crooked_penny",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
