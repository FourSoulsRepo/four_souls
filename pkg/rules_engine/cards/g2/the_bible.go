package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Bible (Active Treasure Card)
//
//	{Tap Effect}Choose one-
//	Destroy a curse.
//	Cancel your attack.
//	Kill a monster named Mom!, Mom's Heart!, or It Lives!
var theBible = engine.CardDef{
	Ref:    "the_bible",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{{
		Kind:  engine.Activated,
		Text:  "↷: Choose one- Destroy a curse. Cancel your attack. Kill a monster named Mom!, Mom's Heart!, or It Lives!",
		Costs: []engine.Cost{engine.Tap()},
		Modes: []engine.Mode{
			{Text: "Destroy a curse.", Targets: []engine.TargetSpec{engine.Choose(engine.TargetCurse)}, Effects: []engine.Effect{engine.DestroyTarget(0)}},
			{Text: "Cancel your attack.", Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				if c.Controller == c.G.Turn.Active {
					c.Do(engine.CancelAttack())
				}
			})}},
			{Text: "Kill a monster named Mom!, Mom's Heart!, or It Lives!", Targets: []engine.TargetSpec{engine.ChooseWhere(engine.TargetMonster, namedMom)}, Effects: []engine.Effect{engine.Kill(0)}},
		},
	}},
}
