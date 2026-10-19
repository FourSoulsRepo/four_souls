package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Dice Room!! (Room Card)
//
//	When this enters play, each player discards their hands and loots 3.
//	At the end of the turn, put this into discard.
var diceRoom2 = engine.CardDef{
	Ref:       "dice_room_2",
	Kind:      engine.RoomCard,
	Copies:    1,
	Abilities: diceRoom("When this enters play, each player discards their hands and loots 3.", engine.EachPlayer(discardHandLoot(3))),
}
