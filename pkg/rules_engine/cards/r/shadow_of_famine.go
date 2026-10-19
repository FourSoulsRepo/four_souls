package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Shadow Of Famine (Room Card)
//
//	Players can't activate more than one ability each turn.
//	Players can't play more than one loot card each turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var shadowOfFamine = engine.CardDef{
	Ref:    "shadow_of_famine",
	Kind:   engine.RoomCard,
	Copies: 1,
}
