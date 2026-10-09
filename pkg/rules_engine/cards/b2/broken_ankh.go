package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Broken Ankh (Trinket Card)
//
//	When you would die, roll-
//	6: Prevent death. If it's your turn, cancel everything that hasn't resolved and end it.
//	-Trinket- This loot becomes an item under your control when it resolves.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var brokenAnkh = engine.CardDef{
	Ref:    "broken_ankh",
	Kind:   engine.LootCard,
	Copies: 1,
}
