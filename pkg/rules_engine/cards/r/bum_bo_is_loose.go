package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Bum-Bo Is Loose! (Room Card)
//
//	At the end of the turn, the active player loses 2¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var bumBoIsLoose = engine.CardDef{
	Ref:    "bum_bo_is_loose",
	Kind:   engine.RoomCard,
	Copies: 1,
}
