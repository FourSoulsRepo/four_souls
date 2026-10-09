package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Cursed Fatty (Cursed Monster Card)
//
//	{Curse Effect}Each time a player rolls a ➄, they discard a loot card.
var cursedFatty = engine.CardDef{
	Ref:     "cursed_fatty",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      2,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player rolls a 5, they discard a loot card.",
			Trigger: playerRollOf(5),
			Effects: []engine.Effect{rollerDiscards},
		},
	},
}
