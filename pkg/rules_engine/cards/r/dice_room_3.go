package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Dice Room!!! (Room Card)
//
//	When this enters play, put each shop item into discard.
//	At the end of the turn, put this into discard.
var diceRoom3 = engine.CardDef{
	Ref:       "dice_room_3",
	Kind:      engine.RoomCard,
	Copies:    1,
	Abilities: diceRoom("When this enters play, put each shop item into discard.", discardShop),
}
