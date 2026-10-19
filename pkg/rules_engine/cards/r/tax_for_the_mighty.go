package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Tax For The Mighty (Room Card)
//
//	When any player controls a soul, players who control the most souls or tied for the most must pay each other player 2¢ to attack.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var taxForTheMighty = engine.CardDef{
	Ref:    "tax_for_the_mighty",
	Kind:   engine.RoomCard,
	Copies: 1,
}
