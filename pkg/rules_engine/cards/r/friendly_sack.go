package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Friendly Sack (Active Treasure Card)
//
//	{Tap Effect}Look at the top 3 cards of the loot deck. Put one in your hand and put the rest in another player's hand.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var friendlySack = engine.CardDef{
	Ref:    "friendly_sack",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
