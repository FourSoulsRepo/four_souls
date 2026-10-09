package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// VI. The Lovers (Wildcard Card)
//
//	Choose a player.
//	They gain +2{HP} till end of turn.
var viTheLovers = engine.CardDef{
	Ref:    "vi_the_lovers",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Choose a player. They gain +2 HP till end of turn.",
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetPlayer)},
			Effects: []engine.Effect{engine.GainHPThisTurn(2, 0)},
		},
	},
}
