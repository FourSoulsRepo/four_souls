package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Blessing Of The Inner Eye (Room Card)
//
//	Each time a player declares an attack, before choosing what to attack, they may look at the top 3 cards of the monster deck and put them back in any order.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var blessingOfTheInnerEye = engine.CardDef{
	Ref:    "blessing_of_the_inner_eye",
	Kind:   engine.RoomCard,
	Copies: 1,
}
