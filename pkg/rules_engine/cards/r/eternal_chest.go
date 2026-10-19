package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Eternal Chest (Room Card)
//
//	{Paid Effect}Discard a loot card: Roll-
//	3: Gain 4¢.
//	4-5: Gain 7¢.
//	6: Gain +1 treasure and loot 1.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var eternalChest = engine.CardDef{
	Ref:    "eternal_chest",
	Kind:   engine.RoomCard,
	Copies: 1,
}
