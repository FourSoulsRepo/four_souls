package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Gus (Room Card)
//
//	This room can be attacked.
//	When an attack is declared on this, each non-active player rolls:
//	4-6: They must make an attack roll against this after each attack roll the active player makes this attack.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var gus = engine.CardDef{
	Ref:    "gus",
	Kind:   engine.RoomCard,
	Copies: 1,
	HP:     9,
	DC:     3,
	ATK:    1,
	Soul:   1,
}
