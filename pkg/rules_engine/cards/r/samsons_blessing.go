package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Samson’s Blessing (Room Card)
//
//	Players have +1{ATK}.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var samsonsBlessing = engine.CardDef{
	Ref:    "samsons_blessing",
	Kind:   engine.RoomCard,
	Copies: 1,
}
