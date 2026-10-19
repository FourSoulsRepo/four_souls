package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Eden’s Blessing (Room Card)
//
//	At the start of the turn, the active player may reroll an item they control.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var edensBlessingR = engine.CardDef{
	Ref:    "edens_blessing_r",
	Kind:   engine.RoomCard,
	Copies: 1,
}
