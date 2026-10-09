package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Curse Of Amnesia (Curse Card)
//
//	{Curse Effect}At the end of your turn, discard 2 loot cards.
//	-Curse- When this enters play, give this to a player. When they die, they put this into discard.
var curseOfAmnesia = engine.CardDef{
	Ref:    "curse_of_amnesia",
	Kind:   engine.EventCard,
	Copies: 1,
	Curse:  true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "At the end of your turn, discard 2 loot cards.",
			Trigger: engine.AtEndOfYourTurn(),
			Effects: []engine.Effect{discardOne, discardOne},
		},
	},
}
