package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Modeling Clay (Active Treasure Card)
//
//	{Tap Effect}Choose a non-eternal item. This becomes a copy of that item.
//	(This change is indefinite.)
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var modelingClay = engine.CardDef{
	Ref:    "modeling_clay",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
