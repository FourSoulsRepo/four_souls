package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Baby Haunt (Passive Treasure Card)
//
//	{Curse Effect}Monsters have +1{DC} on your turn.
//	When you die, before paying penalties, give this to another player.
var babyHaunt = engine.CardDef{
	Ref:    "baby_haunt",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Statics: []engine.Static{{Stat: engine.StatMonsterDC, Amount: 1, Applies: func(g *engine.Game, self engine.ObjectID, _ engine.PlayerID, _ engine.ObjectID) bool {
		return g.Turn.Active == g.Object(self).Controller
	}}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When you die, before paying penalties, give this to another player.",
			Trigger: engine.WhenYouDie(),
			Effects: []engine.Effect{giveThisAway},
		},
	},
}
