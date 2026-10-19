package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Dice Room!!!!!! (Room Card)
//
//	When this enters play, reroll each item in play, each player discards their hand and loots 3. Put each monster into discard.
//	At the end of the turn, put this into discard.
var diceRoom6 = engine.CardDef{
	Ref:    "dice_room_6",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: diceRoom("When this enters play, reroll each item in play, each player discards their hand and loots 3. Put each monster into discard.",
		engine.EffectFunc(func(c *engine.Ctx) {
			rerollAll(c.G, append(itemsOf(c.G, engine.NoPlayer), shopItems(c.G)...))
		}),
		engine.EachPlayer(discardHandLoot(3)),
		eachMonsterToDiscard(false)),
}
