package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Boomerang (Active Treasure Card)
//
//	{Tap Effect}Choose another player. Steal a loot card from them at random.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var boomerang = engine.CardDef{
	Ref:    "boomerang",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
