package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Charmed Bony (Holy/Charmed Monster Card)
//
//	When this dies, after gaining rewards, choose another player. They gain 5¢.
//	Each time the attacking player rolls an attack roll of 6, they deal 1 damage to a player.
var charmedBony = engine.CardDef{
	Ref:     "charmed_bony",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 5}},
	Abilities: []engine.Ability{
		afterRewardsAnotherPlayer("They gain 5¢.", engine.GainCents(5)),
		{
			Kind:    engine.Triggered,
			Text:    "Each time the attacking player rolls an attack roll of 6, they deal 1 damage to a player.",
			Trigger: cardkit.AttackRollOnThis(6),
			Effects: []engine.Effect{cardkit.ForRoller(cardkit.DamageAPlayer(1, false))},
		},
	},
}
