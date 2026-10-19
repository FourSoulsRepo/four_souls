package r

import engine "github.com/FourSoulsRepo/rules_engine"

// All Hallows’ Eve (Room Card)
//
//	Players can't gain souls.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var allHallowsEve = engine.CardDef{
	Ref:    "all_hallows_eve",
	Kind:   engine.RoomCard,
	Copies: 1,
}
