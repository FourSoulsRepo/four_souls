package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Cursed Horf (Cursed Monster Card)
//
//	{Curse Effect}Each time a player rolls a ➁, they take 2 damage.
var cursedHorf = engine.CardDef{
	Ref:     "cursed_horf",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 3}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player rolls a 2, they take 2 damage.",
			Trigger: playerRollOf(2),
			Effects: []engine.Effect{forRoller(engine.DealDamage(2, engine.You))},
		},
	},
}
