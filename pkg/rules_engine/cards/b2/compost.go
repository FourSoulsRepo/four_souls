package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Compost (Active Treasure Card)
//
//	{Tap Effect}The next time a player would loot, they loot from the top of the loot discard instead.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var compost = engine.CardDef{
	Ref:    "compost",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
