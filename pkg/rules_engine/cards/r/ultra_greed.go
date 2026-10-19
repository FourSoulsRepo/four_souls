package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Ultra Greed! (Epic Boss Card)
//
//	After each attack roll the active player makes against this, each other player in turn order makes an attack roll against this.
//	The player who kills this gains its rewards.
var ultraGreed = engine.CardDef{
	Ref:             "ultra_greed",
	Kind:            engine.MonsterCard,
	Copies:          1,
	HP:              8,
	DC:              5,
	ATK:             1,
	Soul:            1,
	RewardsToKiller: true,
	Rewards:         []engine.Reward{{Kind: engine.RewardCents, Amount: 40}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "After each attack roll the active player makes against this, each other player in turn order makes an attack roll against this.",
		Trigger: whenAttackedThis(),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			c.G.AddAttackHelpers(othersOfActive(c.G)...)
		})},
	}},
}
