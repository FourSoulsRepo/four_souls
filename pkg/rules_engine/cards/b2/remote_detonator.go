package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Remote Detonator (Active Treasure Card)
//
//	{Tap Effect}Each player votes on an item in play. Destroy the item with the most votes. If there is a tie, nothing happens.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var remoteDetonator = engine.CardDef{
	Ref:    "remote_detonator",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
