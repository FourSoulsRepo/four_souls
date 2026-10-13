package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Lord Of The Pit (Eternal Treasure Card)
//
//	{Tap Effect}Cancel an attack. That player may attack an additional time this turn.
//	-Eternal- This can't be destroyed or put into discard.
var lordOfThePit = engine.CardDef{
	Ref:     "lord_of_the_pit",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Cancel an attack. That player may attack an additional time this turn.",
			Costs:   []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{engine.CancelAttack(engine.ActivePlayerAttacksAgain())},
		},
	},
}
