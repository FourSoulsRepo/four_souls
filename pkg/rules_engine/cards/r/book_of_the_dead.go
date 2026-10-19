package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Book Of The Dead (Active Treasure Card)
//
//	{Tap Effect}Put a non-event monster card in discard in a monster slot not being attacked.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var bookOfTheDead = engine.CardDef{
	Ref:    "book_of_the_dead",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
