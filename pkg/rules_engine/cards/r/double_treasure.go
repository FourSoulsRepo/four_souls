package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Double Treasure! (Good Event Card)
//
//	Choose another player. You and that player each gain +1 treasure.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var doubleTreasure = engine.CardDef{
	Ref:    "double_treasure",
	Kind:   engine.EventCard,
	Copies: 1,
}
