package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Joker (Wildcard Card)
//
//	Look at a player's hand. You may steal a loot card from them.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var joker = engine.CardDef{
	Ref:    "joker",
	Kind:   engine.LootCard,
	Copies: 1,
}
