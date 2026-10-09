package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// VII. The Chariot (Wildcard Card)
//
//	Choose a player.
//	They gain +1{ATK} and +1{HP} till end of turn.
var viiTheChariot = engine.CardDef{
	Ref:    "vii_the_chariot",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Choose a player. They gain +1 ATK and +1 HP till end of turn.",
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetPlayer)},
			Effects: []engine.Effect{engine.GainATKThisTurn(1, 0), engine.GainHPThisTurn(1, 0)},
		},
	},
}
