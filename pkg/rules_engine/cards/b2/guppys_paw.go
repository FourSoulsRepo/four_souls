package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Guppy’s Paw (Active Treasure Card)
//
//	{Tap Effect}Pay 1{HP}. If you do, choose a player. Prevent the next instance of up to 2 damage they would take this turn.
//	-Guppy- The first player to control 2 or more Guppy items gains the Soul of Guppy.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var guppysPaw = engine.CardDef{
	Ref:    "guppys_paw",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
