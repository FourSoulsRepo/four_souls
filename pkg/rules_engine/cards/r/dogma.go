package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Dogma (Epic Boss Card)
//
//	When an attack is declared on this, put each other monster into discard.
//	When this dies, put a card named The Harbingers from outside the game into a monster slot not being attacked.
var dogma = engine.CardDef{
	Ref:     "dogma",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      3,
	ATK:     3,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 2}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When an attack is declared on this, put each other monster into discard.",
			Trigger: whenAttackedThis(),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				for _, id := range monstersInPlay(c.G) {
					if id != c.Source {
						c.G.DiscardMonster(id)
					}
				}
				c.G.RefillSlots()
			})},
		},
		{
			Kind:    engine.Triggered,
			Text:    "When this dies, put a card named The Harbingers from outside the game into a monster slot not being attacked.",
			Trigger: engine.WhenThisDies(),
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				if slots := freeSlots(c.G); a[0] >= 0 && a[0] < len(slots) {
					c.G.PlaceFromOutside("the_harbingers", slots[a[0]])
				}
			}, slotQuestion("Put The Harbingers in which monster slot?"))},
		},
	},
}
