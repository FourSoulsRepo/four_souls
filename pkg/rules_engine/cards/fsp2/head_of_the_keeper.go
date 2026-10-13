package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Head Of The Keeper (Passive Treasure Card)
//
//	Each time you deal damage, gain 1¢.
//	Each time you kill a monster or player, gain 2¢.
var headOfTheKeeper = engine.CardDef{
	Ref:    "head_of_the_keeper",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{Kind: engine.Triggered, Text: "Each time you deal damage, gain 1¢.", Trigger: whenYouDealDamage(), Effects: []engine.Effect{engine.GainCents(1)}},
		{Kind: engine.Triggered, Text: "Each time you kill a monster or player, gain 2¢.", Trigger: whenYouKill(), Effects: []engine.Effect{engine.GainCents(2)}},
	},
}
