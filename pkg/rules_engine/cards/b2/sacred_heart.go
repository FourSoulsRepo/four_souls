package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Sacred Heart (Passive Treasure Card)
//
//	When you would roll a 1, you may change the result to a 6.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var sacredHeart = engine.CardDef{
	Ref:    "sacred_heart",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
