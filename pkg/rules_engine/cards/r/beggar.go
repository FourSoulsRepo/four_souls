package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Beggar (Room Card)
//
//	{Paid Effect}Pay 3¢: Roll-
//	3-4: Gain 2¢
//	5: Loot 2.
//	6: Gain +1 treasure. Put this into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var beggar = engine.CardDef{
	Ref:    "beggar",
	Kind:   engine.RoomCard,
	Copies: 1,
}
