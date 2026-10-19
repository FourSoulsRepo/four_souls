package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Spoils Of War (Room Card)
//
//	Each time a player dies, each other player loots 1.
var spoilsOfWar = engine.CardDef{
	Ref:    "spoils_of_war",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time a player dies, each other player loots 1.",
		Trigger: engine.WhenAPlayerDies(),
		Effects: []engine.Effect{eachOtherThanDead(engine.Loot(1))},
	}},
}
