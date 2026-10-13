package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Cursed Tumor (Cursed Monster Card)
//
//	{Curse Effect}Each time a player rolls a ➃, each player takes 1 damage.
var cursedTumor = engine.CardDef{
	Ref:     "cursed_tumor",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Roll: true}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time a player rolls a 4, each player takes 1 damage.",
		Trigger: engine.OnRollOf(4),
		Effects: []engine.Effect{eachPlayerTakes(1)},
	}},
}
