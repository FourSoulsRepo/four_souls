package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Cursed Chest (Bad Event Card)
//
//	Roll-
//	1-3: Take 1 Damage.
//	4-5: Take 2 Damage.
//	6: Search the treasure deck for a Guppy item, gain it, then shuffle the treasure deck.
var cursedChest = engine.CardDef{
	Ref:    "cursed_chest",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Roll- 1-3: Take 1 damage. 4-5: Take 2 damage. 6: Search the treasure deck for a Guppy item, gain it, then shuffle the treasure deck.",
			Trigger: engine.WhenThisEntersPlay(),
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
				Results(1, 3, engine.DealDamage(1, engine.You)).
				Results(4, 5, engine.DealDamage(2, engine.You)).
				Results(6, 6, searchGuppy))},
		},
	},
}
