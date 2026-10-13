package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Gaper (Basic Monster Card)
//
//	When this dies on an attack roll of 6, the active player must make an additional attack on the monster deck.
var gaper = engine.CardDef{
	Ref:     "gaper",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 3}},
	Abilities: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "When this dies on an attack roll of 6, the active player must make an additional attack on the monster deck.",
		Trigger: engine.Trigger{On: engine.EvDied, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
			return e.Prev == self && g.Object(self).KilledOn == 6
		}},
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.ForceAttacks(1, true) })},
	}},
}
