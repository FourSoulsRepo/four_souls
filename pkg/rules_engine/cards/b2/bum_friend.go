package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Bum Friend (Active Treasure Card)
//
//	{Tap Effect}Loot 1, then put a loot card from your hand on top of the loot deck.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var bumFriend = engine.CardDef{
	Ref:    "bum_friend",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
