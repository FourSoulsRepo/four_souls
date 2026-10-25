package r

import engine "github.com/FourSoulsRepo/rules_engine"

// IBS (Eternal Treasure Card)
//
//	{Paid Effect}Discard a loot card:
//	Choose one- Cancel your attack on a monster. Gain 3¢.
//	{Paid Effect}Discard 2 loot cards:
//	Cancel the ↷ or $ ability of an item or loot being played.
//	-Eternal- This can't be destroyed or put into discard.
var ibs = engine.CardDef{
	Ref:     "ibs",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "$ Discard a loot card: Choose one- Cancel your attack on a monster. Gain 3¢.",
			Costs:   []engine.Cost{engine.DiscardChosen(0)},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetYourHandCard)},
			Modes: []engine.Mode{
				{Text: "Cancel your attack on a monster.", Targets: []engine.TargetSpec{engine.Choose(engine.TargetYourHandCard)}, Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
					g := c.G
					if g.Attack.Started && g.Turn.Active == c.Controller && g.Object(g.Attack.Target).Role == engine.RoleMonster {
						c.Do(engine.CancelAttack())
					}
				})}},
				{Text: "Gain 3¢.", Targets: []engine.TargetSpec{engine.Choose(engine.TargetYourHandCard)}, Effects: []engine.Effect{engine.GainCents(3)}},
			},
		},
		{
			Kind:  engine.Activated,
			Text:  "$ Discard 2 loot cards: Cancel the ↷ or $ ability of an item or loot being played.",
			Costs: []engine.Cost{engine.DiscardChosen(0), engine.DiscardChosen(1)},
			Targets: []engine.TargetSpec{
				engine.Choose(engine.TargetYourHandCard),
				engine.ChooseWhere(engine.TargetYourHandCard, engine.NotChosen),
				engine.Choose(engine.TargetStackAbility),
			},
			Effects: []engine.Effect{engine.CancelTarget(2)},
		},
	},
}
