package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Curse Of Greed (Curse Card)
//
//	{Curse Effect}At the end of your turn, lose 4¢.
//	-Curse- When this enters play, give this to a player. When they die, they put this into discard.
var curseOfGreed = engine.CardDef{
	Ref:    "curse_of_greed",
	Kind:   engine.EventCard,
	Copies: 1,
	Curse:  true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "At the end of your turn, lose 4¢.",
			Trigger: engine.AtEndOfYourTurn(),
			Effects: []engine.Effect{engine.LoseCents(4)},
		},
	},
}
