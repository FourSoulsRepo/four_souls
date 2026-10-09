package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Placebo (Active Treasure Card)
//
//	{Tap Effect}This copies a ↷ ability of a non-eternal item.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var placebo = engine.CardDef{
	Ref:    "placebo",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
