package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Red Champions (Room Card)
//
//	Monsters have +1{HP}.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var redChampions = engine.CardDef{
	Ref:    "red_champions",
	Kind:   engine.RoomCard,
	Copies: 1,
}
