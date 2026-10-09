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
	DamageMod: func(g *engine.Game, self engine.ObjectID, t engine.Target, n int) int {
		if t.IsPlayer && t.Player == g.Object(self).Controller {
			return n + 1
		}
		return n
	},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When you die, before paying penalties, give this to another player.",
			Trigger: engine.WhenYouDie(),
			Effects: []engine.Effect{giveThisAway},
		},
	},
}
