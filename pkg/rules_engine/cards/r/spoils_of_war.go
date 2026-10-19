package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Spoils Of War (Room Card)
//
//	Each time a player dies, each other player loots 1.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var spoilsOfWar = engine.CardDef{
	Ref:    "spoils_of_war",
	Kind:   engine.RoomCard,
	Copies: 1,
}
