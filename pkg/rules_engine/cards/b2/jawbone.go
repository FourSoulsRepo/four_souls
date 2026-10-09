package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Jawbone (Active Treasure Card)
//
//	{Tap Effect}Steal 3¢ from a player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var jawbone = engine.CardDef{
	Ref:    "jawbone",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
