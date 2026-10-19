package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Marbles (Active Treasure Card)
//
//	{Tap Effect}Choose a non-eternal card in play. Put a counter on it of a type already on it or remove a counter from it.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var marbles = engine.CardDef{
	Ref:    "marbles",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
