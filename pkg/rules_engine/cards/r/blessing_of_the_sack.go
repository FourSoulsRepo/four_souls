package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Blessing Of The Sack (Room Card)
//
//	Rewards are doubled.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var blessingOfTheSack = engine.CardDef{
	Ref:    "blessing_of_the_sack",
	Kind:   engine.RoomCard,
	Copies: 1,
}
