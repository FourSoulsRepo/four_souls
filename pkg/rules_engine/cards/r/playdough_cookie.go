package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Playdough Cookie (Passive Treasure Card)
//
//	At the start of your turn, you may choose a non-eternal item. Put a counter on that item or remove one from it. Otherwise, gain 1¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var playdoughCookie = engine.CardDef{
	Ref:    "playdough_cookie",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
