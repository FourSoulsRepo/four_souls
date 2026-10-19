package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Birthright (Active Treasure Card)
//
//	{Tap Effect}This becomes a copy of an eternal item you control. This loses eternal.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var birthright = engine.CardDef{
	Ref:    "birthright",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
