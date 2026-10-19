package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Static Shock (Room Card)
//
//	Each time a player activates an item, they take 1 damage.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var staticShock = engine.CardDef{
	Ref:    "static_shock",
	Kind:   engine.RoomCard,
	Copies: 1,
}
