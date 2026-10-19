package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Red Vise (Room Card)
//
//	At the end of the turn, if the active player has 1 or fewer loot cards in their hand, they take 1 damage.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var redVise = engine.CardDef{
	Ref:    "red_vise",
	Kind:   engine.RoomCard,
	Copies: 1,
}
