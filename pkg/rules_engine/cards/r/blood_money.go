package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Blood Money (Room Card)
//
//	Each time a player dies, each other player gains 4¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var bloodMoney = engine.CardDef{
	Ref:    "blood_money",
	Kind:   engine.RoomCard,
	Copies: 1,
}
