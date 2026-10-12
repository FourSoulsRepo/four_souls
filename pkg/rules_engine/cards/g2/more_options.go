package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// More Options (Passive Treasure Card)
//
//	At the start of your turn, you may look at the top 2 cards of the treasure deck.
//	You may purchase an additional time on your turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var moreOptions = engine.CardDef{
	Ref:    "more_options",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
