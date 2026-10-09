package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Razor Blade (Active Treasure Card)
//
//	{Tap Effect}Deal 1 damage to a player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var razorBlade = engine.CardDef{
	Ref:    "razor_blade",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
