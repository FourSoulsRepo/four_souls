package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Restock Machine (Room Card)
//
//	{Paid Effect}Pay 2¢:
//	Put each shop item or each monster not being attacked into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var restockMachine = engine.CardDef{
	Ref:    "restock_machine",
	Kind:   engine.RoomCard,
	Copies: 1,
}
