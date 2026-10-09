package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Cursed Gaper (Cursed Monster Card)
//
//	{Curse Effect}Each time a player rolls a ➃, each monster gains +1{ATK} till end of turn.
var cursedGaper = engine.CardDef{
	Ref:     "cursed_gaper",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 3}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player rolls a 4, each monster gains +1 ATK till end of turn.",
			Trigger: playerRollOf(4),
			Effects: []engine.Effect{eachMonsterGains(engine.StatMonsterATK)},
		},
	},
}
