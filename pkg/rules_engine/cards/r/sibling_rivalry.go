package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Sibling Rivalry (Eternal Treasure Card)
//
//	{Tap Effect}Choose a non-active player. The next time the active player declares an attack this turn, the chosen player must make an attack roll after each attack roll the active player makes for the attack. If that monster dies this attack, the chosen player also gains the rewards.
//	-Eternal- This can't be destroyed or put into discard.
var siblingRivalry = engine.CardDef{
	Ref:       "sibling_rivalry",
	Kind:      engine.TreasureCard,
	Copies:    1,
	Eternal:   true,
	Outside:   true,
	Tap:       true,
	ThreePlus: true,
	Abilities: []engine.Ability{{
		Kind:    engine.Activated,
		Text:    "↷: Choose a non-active player. The next time the active player declares an attack this turn, the chosen player must make an attack roll after each attack roll the active player makes for the attack. If that monster dies this attack, the chosen player also gains the rewards.",
		Costs:   []engine.Cost{engine.Tap()},
		Targets: []engine.TargetSpec{engine.ChooseWhere(engine.TargetPlayer, func(g *engine.Game, _ engine.ObjectID, c engine.Chosen) bool { return c.Player != g.Turn.Active })},
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			c.G.HelpNextAttack(c.Targets[0].Player, true)
		})},
	}},
}
