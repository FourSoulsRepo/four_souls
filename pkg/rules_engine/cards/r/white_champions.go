package r

import engine "github.com/FourSoulsRepo/rules_engine"

// White Champions (Room Card)
//
//	Monsters have +1{DC}.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var whiteChampions = engine.CardDef{
	Ref:    "white_champions",
	Kind:   engine.RoomCard,
	Copies: 1,
}
