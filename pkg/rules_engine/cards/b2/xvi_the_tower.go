package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// XVI. The Tower (Wildcard Card)
//
//	Roll-
//	1-2: Each player takes 1 damage.
//	3-4: Each monster takes 1 damage.
//	5-6: Each player takes 2 damage.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var xviTheTower = engine.CardDef{
	Ref:    "xvi_the_tower",
	Kind:   engine.LootCard,
	Copies: 1,
}
