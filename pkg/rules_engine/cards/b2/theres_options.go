package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// There’s Options (Passive Treasure Card)
//
//	You may look at the top card of the treasure deck at any time on your turn.
//	You may purchase an additional time on your turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theresOptions = engine.CardDef{
	Ref:    "theres_options",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
