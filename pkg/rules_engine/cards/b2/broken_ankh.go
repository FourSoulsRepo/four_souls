package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Broken Ankh (Trinket Card)
//
//	When you would die, roll-
//	6: Prevent death. If it's your turn, cancel everything that hasn't resolved and end it.
//	-Trinket- This loot becomes an item under your control when it resolves.
var brokenAnkh = engine.CardDef{
	Ref:     "broken_ankh",
	Kind:    engine.LootCard,
	Copies:  1,
	Trinket: true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When you would die, roll- 6: Prevent death. If it's your turn, cancel everything that hasn't resolved and end it.",
			Trigger: engine.WhenYouWouldDie(),
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.Results(6, 6, engine.PreventYourDeath(), endYourTurn))},
		},
	},
}
