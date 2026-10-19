package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Blessing Of Gluttony (Room Card)
//
//	The active player loots +1 during their loot step.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var blessingOfGluttony = engine.CardDef{
	Ref:    "blessing_of_gluttony",
	Kind:   engine.RoomCard,
	Copies: 1,
}
