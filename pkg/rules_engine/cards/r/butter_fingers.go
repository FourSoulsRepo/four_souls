package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Butter Fingers (Room Card)
//
//	At the end of each turn, the active player discards a loot card.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var butterFingers = engine.CardDef{
	Ref:    "butter_fingers",
	Kind:   engine.RoomCard,
	Copies: 1,
}
