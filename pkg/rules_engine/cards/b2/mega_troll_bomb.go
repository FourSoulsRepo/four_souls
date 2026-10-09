package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mega Troll Bomb! (Bad Event Card)
//
//	Each player takes 2 damage!
var megaTrollBomb = engine.CardDef{
	Ref:    "mega_troll_bomb",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each player takes 2 damage!",
			Trigger: engine.WhenThisEntersPlay(),
			Effects: []engine.Effect{engine.EachPlayer(engine.DealDamage(2, engine.You))},
		},
	},
}
