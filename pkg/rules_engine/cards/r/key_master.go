package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Key Master (Room Card)
//
//	{Paid Effect}Discard a loot card: Roll. If the discarded card had "key" in its name, roll 3x instead-
//	5-6: Gain +1 treasure.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var keyMaster = engine.CardDef{
	Ref:    "key_master",
	Kind:   engine.RoomCard,
	Copies: 1,
}
