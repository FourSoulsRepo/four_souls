package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// III. The Empress (Wildcard Card)
//
//	Choose a player.
//	They gain +1{ATK} and +1 to dice rolls till end of turn.
var iiiTheEmpress = engine.CardDef{
	Ref:    "iii_the_empress",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Choose a player. They gain +1 ATK and +1 to dice rolls till end of turn.",
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetPlayer)},
			Effects: []engine.Effect{engine.GainATKThisTurn(1, 0), engine.GainRollBonusThisTurn(1, 0)},
		},
	},
}
