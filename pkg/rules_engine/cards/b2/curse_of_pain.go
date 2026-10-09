package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Curse Of Pain (Curse Card)
//
//	{Curse Effect}At the start of your turn, take 1 damage.
//	-Curse- When this enters play, give this to a player. When they die, they put this into discard.
var curseOfPain = engine.CardDef{
	Ref:    "curse_of_pain",
	Kind:   engine.EventCard,
	Copies: 1,
	Curse:  true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "At the start of your turn, take 1 damage.",
			Trigger: engine.AtStartOfYourTurn(),
			Effects: []engine.Effect{engine.DealDamage(1, engine.You)},
		},
	},
}
