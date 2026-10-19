package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Fortune Teller (Room Card)
//
//	{Paid Effect}Pay 3¢: roll-
//	3: Loot 1.
//	4: Loot 2.
//	5: Loot 3.
//	6: Loot 6. Put this into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var fortuneTeller = engine.CardDef{
	Ref:    "fortune_teller",
	Kind:   engine.RoomCard,
	Copies: 1,
}
