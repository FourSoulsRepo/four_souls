package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Blessing Of Greed (Room Card)
//
//	At the start of each turn, the active player gains 3¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var blessingOfGreed = engine.CardDef{
	Ref:    "blessing_of_greed",
	Kind:   engine.RoomCard,
	Copies: 1,
}
