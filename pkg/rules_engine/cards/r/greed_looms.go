package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Greed Looms (Room Card)
//
//	At the end of the turn, if the active player didn't purchase a shop item, they discard their hand.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var greedLooms = engine.CardDef{
	Ref:    "greed_looms",
	Kind:   engine.RoomCard,
	Copies: 1,
}
