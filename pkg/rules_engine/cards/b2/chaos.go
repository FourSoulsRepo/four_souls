package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Chaos (Active Treasure Card)
//
//	{Tap Effect}Each player gives their hand to the player to their left.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var chaos = engine.CardDef{
	Ref:    "chaos",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
