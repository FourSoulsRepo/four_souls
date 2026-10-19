package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Equality (Room Card)
//
//	If the player would gain any amount of ¢, instead each player gains that much ¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var equality = engine.CardDef{
	Ref:    "equality",
	Kind:   engine.RoomCard,
	Copies: 1,
}
