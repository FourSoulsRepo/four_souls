package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Holy Chest (Good Event Card)
//
//	Roll-
//	1-2: Prevent the next 2 damage you would take this turn. You may attack an additional time this turn.
//	3-5: Gain 7¢.
//	6: This becomes a soul. Gain it.
var holyChest = engine.CardDef{
	Ref:    "holy_chest",
	Kind:   engine.EventCard,
	Copies: 1,
	Soul:   1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Roll- 1-2: Prevent the next 2 damage you would take this turn. You may attack an additional time this turn. 3-5: Gain 7¢. 6: This becomes a soul. Gain it.",
		Trigger: engine.WhenThisEntersPlay(),
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
			Results(1, 2, engine.PreventDamage(2, engine.You), engine.ActivePlayerAttacksAgain()).
			Results(3, 5, engine.GainCents(7)).
			Results(6, 6, engine.BecomeSoul()))},
	}},
}
