package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Black Champions (Room Card)
//
//	Monsters have +1{ATK}.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var blackChampions = engine.CardDef{
	Ref:    "black_champions",
	Kind:   engine.RoomCard,
	Copies: 1,
}
