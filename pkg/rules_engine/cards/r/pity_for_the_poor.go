package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Pity For The Poor (Room Card)
//
//	Players who control the fewest souls or tied for fewest may purchase a shop item for 0¢ on their turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var pityForThePoor = engine.CardDef{
	Ref:    "pity_for_the_poor",
	Kind:   engine.RoomCard,
	Copies: 1,
}
