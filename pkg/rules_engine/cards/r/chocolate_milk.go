package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Chocolate Milk (Passive Treasure Card)
//
//	At the start of your turn, put a counter on this. Then if this has 3+ counters, remove all counters from this and deal 1 damage to a player or monster.
//	While this has 2 counters on it, you have +1{ATK}.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var chocolateMilk = engine.CardDef{
	Ref:    "chocolate_milk",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
