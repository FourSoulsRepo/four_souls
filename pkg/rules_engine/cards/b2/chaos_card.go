package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Chaos Card (One-Use Treasure Card)
//
//	{Tap Effect}Destroy this. If you do, choose one- Kill a player or monster. Destroy an item or soul.
var chaosCard = engine.CardDef{
	Ref:    "chaos_card",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:  engine.Activated,
			Text:  "↷: Destroy this. If you do, choose one- Kill a player or monster. Destroy an item or soul.",
			Costs: []engine.Cost{engine.Tap()},
			Modes: []engine.Mode{
				{Text: "Kill a player or monster.", Targets: []engine.TargetSpec{engine.Choose(engine.TargetMonsterOrPlayer)}, Effects: []engine.Effect{engine.DestroyThis(engine.Kill(0))}},
				{Text: "Destroy an item or soul.", Targets: []engine.TargetSpec{engine.Choose(engine.TargetItemOrSoul)}, Effects: []engine.Effect{engine.DestroyThis(engine.DestroyTarget(0))}},
			},
		},
	},
}
