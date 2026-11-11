package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Daddy Haunt (Passive Treasure Card)
//
//	{Curse Effect}If you would take any amount of damage, take that much damage +1 instead.
//	When you die, before paying penalties, give this to another player.
var daddyHaunt = engine.CardDef{
	Ref:    "daddy_haunt",
	Kind:   engine.TreasureCard,
	Copies: 1,
	// The engine's one rule for "damage you would take is increased"
	// (StatDamageTaken): it applies before cards that lower damage.
	Statics: []engine.Static{engine.YouHave(engine.StatDamageTaken, 1)},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When you die, before paying penalties, give this to another player.",
			Trigger: engine.WhenYouDie(),
			Effects: []engine.Effect{giveThisAway},
		},
	},
}
