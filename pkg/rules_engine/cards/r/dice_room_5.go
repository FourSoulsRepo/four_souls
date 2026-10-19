package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Dice Room!!!!! (Room Card)
//
//	When this enters play, put each monster into discard.
//	At the end of the turn, put this into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var diceRoom5 = engine.CardDef{
	Ref:    "dice_room_5",
	Kind:   engine.RoomCard,
	Copies: 1,
}
