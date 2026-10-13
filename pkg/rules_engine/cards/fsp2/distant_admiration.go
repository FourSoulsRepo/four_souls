package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Distant Admiration (Active Treasure Card)
//
//	{Tap Effect}Choose another player. They give you a loot card. Reveal it, then you must play that loot card if able. This doesn't use a loot play.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var distantAdmiration = engine.CardDef{
	Ref:    "distant_admiration",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
