package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Peep (Boss Card)
//
//	When this dies, search the monster deck for a card named The Bloat and put it in a monster slot not being attacked, then shuffle the monster deck.
var peep = engine.CardDef{
	Ref:     "peep",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When this dies, search the monster deck for a card named The Bloat and put it in a monster slot not being attacked, then shuffle the monster deck.",
			Trigger: engine.WhenThisDies(),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				for _, id := range c.G.Decks[engine.MonsterDeck] {
					if c.G.Object(id).Card == "the_bloat" {
						if slots := freeSlots(c.G); len(slots) > 0 {
							c.G.PlaceFromDeck(id, slots[0])
						}
						break
					}
				}
				c.G.ShuffleDeck(engine.MonsterDeck)
			})},
		},
	},
}
