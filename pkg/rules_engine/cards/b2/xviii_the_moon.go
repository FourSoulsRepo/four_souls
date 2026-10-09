package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// XVIII. The Moon (Wildcard Card)
//
//	Look at the top 5 cards of the loot deck. Put 1 on top and the rest on the bottom.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var xviiiTheMoon = engine.CardDef{
	Ref:    "xviii_the_moon",
	Kind:   engine.LootCard,
	Copies: 1,
}
