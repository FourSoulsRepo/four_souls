package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Planetarium (Room Card)
//
//	At the start of the turn, the active player may gain +1 treasure.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var planetarium = engine.CardDef{
	Ref:    "planetarium",
	Kind:   engine.RoomCard,
	Copies: 1,
}
