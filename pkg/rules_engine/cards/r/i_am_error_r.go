package r

import engine "github.com/FourSoulsRepo/rules_engine"

// I Am Error (Room Card)
//
//	At the start of your turn roll-
//	1-2: Put a monster into discard and replace it with the top card of the monster deck
//	3-4: Reroll an item you control.
//	5: Discard your hand and loot 3.
//	6: Do all of the above.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var iAmErrorR = engine.CardDef{
	Ref:    "i_am_error_r",
	Kind:   engine.RoomCard,
	Copies: 1,
}
