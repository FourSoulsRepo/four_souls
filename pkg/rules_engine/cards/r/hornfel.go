package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Hornfel (Boss Card)
//
//	When this would deal combat damage to the active player, prevent it, then this deals 1 damage to a player chosen at random.
var hornfel = engine.CardDef{
	Ref:     "hornfel",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      5,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Abilities: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "When this would deal combat damage to the active player, prevent it, then this deals 1 damage to a player chosen at random.",
		Trigger: engine.Trigger{On: engine.EvDamagePending, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
			return e.Source == self && e.Text == "combat" && e.Player == g.Turn.Active
		}},
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			c.G.CancelStackItem(c.EventStack)
		}), randomPlayerTakes(1, false)},
	}},
}
