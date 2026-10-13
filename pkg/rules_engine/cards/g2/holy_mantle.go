package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Holy Mantle (Eternal Treasure Card)
//
//	{Tap Effect}Choose a player. The next time that player would die this turn, prevent it. If it's their turn, cancel everything that hasn't resolved and end it.
//	At the end of your turn, recharge this.
//	-Eternal- This can't be destroyed or put into discard.
var holyMantle = engine.CardDef{
	Ref:     "holy_mantle",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Choose a player. The next time that player would die this turn, prevent it. If it's their turn, cancel everything that hasn't resolved and end it.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetPlayer)},
			Effects: []engine.Effect{engine.WardDeath(0)},
		},
		rechargeAtEndOfTurn,
	},
}
