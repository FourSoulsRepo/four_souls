package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Potato Peeler (Active Treasure Card)
//
//	{Tap Effect}Put the top card of each deck into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var potatoPeeler = engine.CardDef{
	Ref:    "potato_peeler",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
