package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Slot Machine (Room Card)
//
//	{Paid Effect}Pay 3¢: roll-
//	2: Gain 1¢
//	3: Gain 2¢
//	4: Gain 5¢
//	5: Gain 7¢
//	6: Gain 15¢. Put this into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var slotMachine = engine.CardDef{
	Ref:    "slot_machine",
	Kind:   engine.RoomCard,
	Copies: 1,
}
