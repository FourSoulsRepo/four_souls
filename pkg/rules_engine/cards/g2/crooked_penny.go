package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Crooked Penny (Active Treasure Card)
//
//	{Tap Effect}Choose a player, then roll-
//	1-3: Double the amount of ¢ they have.
//	4-6: That player loses all ¢.
var crookedPenny = engine.CardDef{
	Ref:    "crooked_penny",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{{
		Kind:    engine.Activated,
		Text:    "↷: Choose a player, then roll- 1-3: Double the amount of ¢ they have. 4-6: That player loses all ¢.",
		Costs:   []engine.Cost{engine.Tap()},
		Targets: []engine.TargetSpec{engine.Choose(engine.TargetPlayer)},
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
			Results(1, 3, forTarget(engine.EffectFunc(func(c *engine.Ctx) { c.Do(engine.GainCents(c.G.Players[c.Controller].Cents)) }))).
			Results(4, 6, forTarget(engine.EffectFunc(func(c *engine.Ctx) { c.G.LoseCentsNow(c.Controller, c.G.Players[c.Controller].Cents) }))))},
	}},
}
