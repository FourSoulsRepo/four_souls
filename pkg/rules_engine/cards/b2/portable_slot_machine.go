package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Portable Slot Machine (Paid Treasure Card)
//
//	{Paid Effect}Pay 3¢: Roll-
//	1-2: Loot 1.
//	3-4: Gain 4¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var portableSlotMachine = engine.CardDef{
	Ref:    "portable_slot_machine",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
