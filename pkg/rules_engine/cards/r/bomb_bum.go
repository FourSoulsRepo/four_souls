package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Bomb Bum (Room Card)
//
//	{Paid Effect}Discard a loot card: Roll. if the discarded card had "bomb" in its name, roll 3x instead-
//	1-4: Deal 1 damage to a monster or player.
//	5-6: Take 1 damage.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var bombBum = engine.CardDef{
	Ref:    "bomb_bum",
	Kind:   engine.RoomCard,
	Copies: 1,
}
