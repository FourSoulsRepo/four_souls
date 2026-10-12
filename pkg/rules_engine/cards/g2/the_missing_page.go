package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Missing Page (Passive Treasure Card)
//
//	Each time a player would roll a ❺, you may force that player to reroll it.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theMissingPage = engine.CardDef{
	Ref:    "the_missing_page",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
