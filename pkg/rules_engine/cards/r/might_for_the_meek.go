package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Might For The Meek (Room Card)
//
//	Players who control the fewest souls or tied for fewest have +1{ATK} and may attack an additional time on their turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var mightForTheMeek = engine.CardDef{
	Ref:    "might_for_the_meek",
	Kind:   engine.RoomCard,
	Copies: 1,
}
