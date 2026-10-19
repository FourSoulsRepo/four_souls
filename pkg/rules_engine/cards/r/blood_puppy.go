package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Blood Puppy (Passive Treasure Card)
//
//	Each time a player dies, this levels up.
//	{LV4 Effect}You have +1{ATK}.
//	{LV8 Effect}At the start of your turn, deal 1 damage to a monster or player.
//	{LV10 Effect}At the start of your turn roll-
//	1-3: You die.
var bloodPuppy = engine.CardDef{
	Ref:    "blood_puppy",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player dies, this levels up.",
			Trigger: engine.WhenAPlayerDies(),
			Effects: []engine.Effect{levelUp},
		},
		{
			Kind:    engine.Triggered,
			Text:    "LV8: At the start of your turn, deal 1 damage to a monster or player.",
			Trigger: atStartOfYourTurnAtLevel(8),
			Effects: []engine.Effect{damageAnything(1)},
		},
		{
			Kind:    engine.Triggered,
			Text:    "LV10: At the start of your turn roll- 1-3: You die.",
			Trigger: atStartOfYourTurnAtLevel(10),
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.Results(1, 3, killYou))},
		},
	},
	Statics: []engine.Static{{Stat: engine.StatPlayerATK, Amount: 1, Applies: func(g *engine.Game, self engine.ObjectID, p engine.PlayerID, _ engine.ObjectID) bool {
		o := g.Object(self)
		return o.Controller == p && o.CountersOf(level) >= 4
	}}},
}
