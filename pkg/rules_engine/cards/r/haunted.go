package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Haunted! (Room Card)
//
//	Each time a player gains a soul, they skip their next turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var haunted = engine.CardDef{
	Ref:    "haunted",
	Kind:   engine.RoomCard,
	Copies: 1,
}
