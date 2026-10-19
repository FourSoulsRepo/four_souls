package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Dice Room! (Room Card)
//
//	When this enters play, the active player rerolls each item they control.
//	At the end of the turn, put this into discard.
var diceRoom1 = engine.CardDef{
	Ref:    "dice_room_1",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: diceRoom("When this enters play, the active player rerolls each item they control.", engine.EffectFunc(func(c *engine.Ctx) {
		rerollAll(c.G, itemsOf(c.G, c.Controller))
	})),
}
