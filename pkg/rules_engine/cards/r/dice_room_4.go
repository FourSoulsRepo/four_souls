package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Dice Room!!!! (Room Card)
//
//	When this enters play, each player rerolls each of their items.
//	At the end of the turn, put this into discard.
var diceRoom4 = engine.CardDef{
	Ref:    "dice_room_4",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: diceRoom("When this enters play, each player rerolls each of their items.", engine.EffectFunc(func(c *engine.Ctx) {
		rerollAll(c.G, itemsOf(c.G, engine.NoPlayer))
	})),
}
