package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Diplopia (Active Treasure Card)
//
//	{Tap Effect}Choose a non-eternal passive item. This becomes a copy of that item till end of turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var diplopia = engine.CardDef{
	Ref:    "diplopia",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
