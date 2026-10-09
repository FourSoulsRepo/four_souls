package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Monster Manual (Active Treasure Card)
//
//	{Tap Effect}Choose a monster. The active player must attack that monster this turn if able.
var monsterManual = engine.CardDef{
	Ref:    "monster_manual",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Choose a monster. The active player must attack that monster this turn if able.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetMonster)},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.Turn.MustAttack = c.Targets[0].Object })},
		},
	},
}
