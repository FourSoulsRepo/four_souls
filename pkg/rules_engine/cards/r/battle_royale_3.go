package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Battle Royale!!!!! (Room Card)
//
//	The active player may attack other players. Attacked players have 5+{DC}.
//	When a player dies, if that player was attacked this turn, that player gives the active player the item they would destroy for the death penalty.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var battleRoyale3 = engine.CardDef{
	Ref:    "battle_royale_3",
	Kind:   engine.RoomCard,
	Copies: 1,
}
