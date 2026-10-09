package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dead Bird (Passive Treasure Card)
//
//	Each time a player rolls a ❸, you may look at their hand and steal a loot card from them.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var deadBird = engine.CardDef{
	Ref:    "dead_bird",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
