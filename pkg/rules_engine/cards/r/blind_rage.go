package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Blind Rage (Room Card)
//
//	The active player must attack the monster deck once each turn if able.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var blindRage = engine.CardDef{
	Ref:    "blind_rage",
	Kind:   engine.RoomCard,
	Copies: 1,
}
