package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Lucky Foot (Active Treasure Card)
//
//	{Tap Effect}Add up to 2 to a non-attack roll.
var luckyFoot = engine.CardDef{
	Ref:    "lucky_foot",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Add up to 2 to a non-attack roll.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.ChooseWhere(engine.TargetDiceRoll, nonAttackRoll)},
			Effects: []engine.Effect{changeRollBy(1, 2)},
		},
	},
}

// nonAttackRoll keeps dice rolls that are not attack rolls.
func nonAttackRoll(g *engine.Game, _ engine.ObjectID, c engine.Chosen) bool {
	it, ok := g.StackItemByID(c.StackID)
	return ok && !it.Attack
}
