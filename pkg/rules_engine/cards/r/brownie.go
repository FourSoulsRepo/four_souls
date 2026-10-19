package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Brownie (Boss Card)
//
//	Each time this takes damage, it gains -1{DC} till end of turn.
var brownie = engine.CardDef{
	Ref:     "brownie",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      6,
	DC:      5,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time this takes damage, it gains -1 DC till end of turn.",
		Trigger: cardkit.WhenThisTakesDamage(false),
		Effects: []engine.Effect{thisGains(engine.StatMonsterDC, -1)},
	}},
}
