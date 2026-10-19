package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Blood Lust (Room Card)
//
//	The active player must attack each turn if able.
//	The first time the active player declares an attack each turn, they must attack an additional time this turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var bloodLustR = engine.CardDef{
	Ref:    "blood_lust_r",
	Kind:   engine.RoomCard,
	Copies: 1,
}
