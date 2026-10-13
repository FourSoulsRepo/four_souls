package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Globin (Basic Monster Card)
//
//	When this reaches 1{HP}, the active player rolls-
//	1-4: This takes 1 damage.
//	5-6: This heals 2{HP}.
var globin = engine.CardDef{
	Ref:     "globin",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 3}},
	Abilities: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "When this reaches 1 HP, the active player rolls- 1-4: This takes 1 damage. 5-6: This heals 2 HP.",
		Trigger: engine.Trigger{On: engine.EvDamaged, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
			return e.Object == self && g.HP(self) == 1
		}},
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
			Results(1, 4, engine.EffectFunc(func(c *engine.Ctx) { c.G.DealDamageTo(engine.Target{Object: c.Source}, 1, engine.NoPlayer, c.Source) })).
			Results(5, 6, healThis(2)))},
	}},
}
