package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dagaz (Pill/Rune Card)
//
//	Choose one- Destroy a curse. Choose a player. Prevent the next 1 damage they would take this turn.
var dagaz = engine.CardDef{
	Ref:    "dagaz",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind: engine.LootAbility,
			Text: "Choose one- Destroy a curse. Choose a player. Prevent the next 1 damage they would take this turn.",
			Modes: []engine.Mode{
				{Text: "Destroy a curse.", Targets: []engine.TargetSpec{engine.Choose(engine.TargetCurse)}, Effects: []engine.Effect{engine.DestroyTarget(0)}},
				{Text: "Choose a player. Prevent the next 1 damage they would take this turn.", Targets: []engine.TargetSpec{engine.Choose(engine.TargetPlayer)}, Effects: []engine.Effect{engine.PreventDamage(1, 0)}},
			},
		},
	},
}
