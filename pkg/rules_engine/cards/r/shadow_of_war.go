package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Shadow Of War (Room Card)
//
//	When a player dies, before paying penalties, they must destroy an item they control.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var shadowOfWar = engine.CardDef{
	Ref:    "shadow_of_war",
	Kind:   engine.RoomCard,
	Copies: 1,
}
