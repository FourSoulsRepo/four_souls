package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// XI. Strength (Wildcard Card)
//
//	Choose a player.
//	They gain +1{ATK} till end of turn and may attack an additional time this turn.
var xiStrength = engine.CardDef{
	Ref:    "xi_strength",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Choose a player. They gain +1 ATK till end of turn and may attack an additional time this turn.",
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetPlayer)},
			Effects: []engine.Effect{engine.GainATKThisTurn(1, 0), engine.AddAttacks(1, 0)},
		},
	},
}
