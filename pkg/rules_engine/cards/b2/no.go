package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// No! (Active Treasure Card)
//
//	{Tap Effect}Cancel the ↷ or $ ability of an item.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var no = engine.CardDef{
	Ref:    "no",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
