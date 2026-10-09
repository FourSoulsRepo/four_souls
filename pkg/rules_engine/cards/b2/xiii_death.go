package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// XIII. Death (Wildcard Card)
//
//	Kill a player.
var xiiiDeath = engine.CardDef{
	Ref:    "xiii_death",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Kill a player.",
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetPlayer)},
			Effects: []engine.Effect{engine.Kill(0)},
		},
	},
}
