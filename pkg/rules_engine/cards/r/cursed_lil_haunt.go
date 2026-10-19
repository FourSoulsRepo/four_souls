package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Cursed Lil Haunt (Cursed Monster Card)
//
//	If a player would roll a dice, they instead roll 2 dice. The player to their left chooses one of the rolls as the result.
var cursedLilHaunt = engine.CardDef{
	Ref:       "cursed_lil_haunt",
	Kind:      engine.MonsterCard,
	Copies:    1,
	HP:        3,
	DC:        3,
	ATK:       1,
	Rewards:   []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
	ThreePlus: true,
	Statics:   []engine.Static{engine.MonstersHave(engine.StatTwoDiceLeftChooses, 1)},
}
