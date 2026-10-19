package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Blood Money (Room Card)
//
//	Each time a player dies, each other player gains 4¢.
var bloodMoney = engine.CardDef{
	Ref:    "blood_money",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time a player dies, each other player gains 4¢.",
		Trigger: engine.WhenAPlayerDies(),
		Effects: []engine.Effect{eachOtherThanDead(engine.GainCents(4))},
	}},
}
