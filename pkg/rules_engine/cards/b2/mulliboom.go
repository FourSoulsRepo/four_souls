package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mulliboom (Basic Monster Card)
//
//	When this dies, the active player deals 3 damage to a player.
var mulliboom = engine.CardDef{
	Ref:     "mulliboom",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      2,
	ATK:     4,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 6}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When this dies, the active player deals 3 damage to a player.",
			Trigger: engine.WhenThisDies(),
			Effects: []engine.Effect{damageAPlayer(3, false)},
		},
	},
}
