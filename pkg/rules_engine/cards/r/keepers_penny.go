package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Keeper’s Penny (Active Treasure Card)
//
//	When you gain or purchase this, gain 9¢.
//	{Tap Effect}Put a counter on this. Gain ¢ equal to the number of counters on this, then give this card to another player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var keepersPenny = engine.CardDef{
	Ref:    "keepers_penny",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
