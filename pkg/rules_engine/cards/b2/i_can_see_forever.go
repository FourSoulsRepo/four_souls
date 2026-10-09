package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// I Can See Forever! (Good Event Card)
//
//	Look at the top 6 cards of the loot deck. Put them back in any order, then loot 1.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var iCanSeeForever = engine.CardDef{
	Ref:    "i_can_see_forever",
	Kind:   engine.EventCard,
	Copies: 1,
}
