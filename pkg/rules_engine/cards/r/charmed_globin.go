package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Charmed Globin (Holy/Charmed Monster Card)
//
//	When this dies, after gaining rewards, choose another player. They loot 3.
//	Each time the attacking player misses an attack roll, this heals 1{HP}.
var charmedGlobin = engine.CardDef{
	Ref:     "charmed_globin",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 3}},
	Abilities: []engine.Ability{
		afterRewardsAnotherPlayer("They loot 3.", engine.Loot(3)),
		{
			Kind: engine.Triggered,
			Text: "Each time the attacking player misses an attack roll, this heals 1 HP.",
			Trigger: engine.Trigger{On: engine.EvRollResolved, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
				return e.Text == "attack" && g.Attack.Started && g.Attack.Target == self && e.Amount < g.Evasion(self)
			}},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.HealObject(c.Source, 1) })},
		},
	},
}
