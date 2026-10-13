package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dingle (Boss Card)
//
//	Each time the active player deals damage to this, they roll-
//	1-2: They take 1 damage.
//	3-4: Each player takes 1 damage.
//	5-6: This takes 1 damage.
var dingle = engine.CardDef{
	Ref:     "dingle",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      3,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 6}},
	Abilities: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "Each time the active player deals damage to this, they roll- 1-2: They take 1 damage. 3-4: Each player takes 1 damage. 5-6: This takes 1 damage.",
		Trigger: engine.Trigger{On: engine.EvDamaged, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
			return e.Object == self && e.By == int(g.Turn.Active)+1
		}},
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
			Results(1, 2, damageActive(1)).
			Results(3, 4, eachPlayerTakes(1)).
			Results(5, 6, engine.EffectFunc(func(c *engine.Ctx) { c.G.DealDamageTo(engine.Target{Object: c.Source}, 1, engine.NoPlayer, c.Source) })))},
	}},
}
