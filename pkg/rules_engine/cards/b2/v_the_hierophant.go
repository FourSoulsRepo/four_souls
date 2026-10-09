package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// V. The Hierophant (Wildcard Card)
//
//	Choose a player or monster. Prevent the next instance of up to 2 damage they would take this turn.
var vTheHierophant = engine.CardDef{
	Ref:    "v_the_hierophant",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Choose a player or monster. Prevent the next instance of up to 2 damage they would take this turn.",
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetMonsterOrPlayer)},
			Effects: []engine.Effect{engine.PreventDamage(2, 0)},
		},
	},
}
