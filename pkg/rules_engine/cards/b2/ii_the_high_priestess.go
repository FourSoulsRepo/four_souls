package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// II. The High Priestess (Wildcard Card)
//
//	Choose a player or monster, then roll-
//	Deal damage to them equal to the result.
var iiTheHighPriestess = engine.CardDef{
	Ref:    "ii_the_high_priestess",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Choose a player or monster, then roll- Deal damage to them equal to the result.",
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetMonsterOrPlayer)},
			Effects: []engine.Effect{engine.Roll(damageEqualToRoll())},
		},
	},
}
