package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dark Arts (Eternal Treasure Card)
//
//	Each time a player rolls a ❻, gain 2¢.
//	Each time another player dies, loot 1.
//	-Eternal- This can't be destroyed or put into discard.
var darkArts = engine.CardDef{
	Ref:     "dark_arts",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Abilities: []engine.Ability{
		{Kind: engine.Triggered, Text: "Each time a player rolls a 6, gain 2¢.", Trigger: engine.OnRollOf(6), Effects: []engine.Effect{engine.GainCents(2)}},
		{Kind: engine.Triggered, Text: "Each time another player dies, loot 1.", Trigger: whenAnotherPlayerDies(), Effects: []engine.Effect{engine.Loot(1)}},
	},
}
