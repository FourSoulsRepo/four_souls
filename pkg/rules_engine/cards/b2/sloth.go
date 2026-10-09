package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Sloth (Boss Card)
//
//	When this dies, the player that killed it discards their hand.
var sloth = engine.CardDef{
	Ref:     "sloth",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 1}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When this dies, the player that killed it discards their hand.",
			Trigger: engine.WhenThisDies(),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				k := c.G.Object(c.Source).Killer()
				if k == engine.NoPlayer {
					return
				}
				for _, id := range append([]engine.ObjectID(nil), c.G.Players[k].Hand...) {
					c.G.DiscardFromHand(k, id)
				}
			})},
		},
	},
}
