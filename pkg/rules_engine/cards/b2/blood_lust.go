package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Blood Lust (Eternal Treasure Card)
//
//	{Tap Effect}Choose a player or monster. They gain +1{ATK} till end of turn.
//	At the end of your turn, recharge this.
//	-Eternal- This can't be destroyed or put into discard.
var bloodLust = engine.CardDef{
	Ref:     "blood_lust",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Choose a player or monster. They gain +1 ATK till end of turn.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetMonsterOrPlayer)},
			Effects: []engine.Effect{engine.GainATKThisTurn(1, 0)},
		},
		rechargeAtEndOfTurn,
	},
}
