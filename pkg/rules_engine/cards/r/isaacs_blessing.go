package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Isaac’s Blessing (Room Card)
//
//	Each time a player would roll a ❶ or ❻, they may reroll it.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var isaacsBlessing = engine.CardDef{
	Ref:    "isaacs_blessing",
	Kind:   engine.RoomCard,
	Copies: 1,
}
