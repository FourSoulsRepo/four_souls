package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Baby Plum (Boss Card)
//
//	Each time this takes damage, the active player rolls-
//	1-2: You take 1 damage.
//	3-4: The player to your left takes 1 damage.
//	5-6: The player to your right takes 1 damage.
var babyPlum = engine.CardDef{
	Ref:       "baby_plum",
	Kind:      engine.MonsterCard,
	Copies:    1,
	HP:        4,
	DC:        2,
	ATK:       1,
	Soul:      1,
	Rewards:   []engine.Reward{{Kind: engine.RewardCents, Roll: true}},
	ThreePlus: true,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time this takes damage, the active player rolls- 1-2: You take 1 damage. 3-4: The player to your left takes 1 damage. 5-6: The player to your right takes 1 damage.",
		Trigger: cardkit.WhenThisTakesDamage(false),
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
			Results(1, 2, damageActive(1)).
			Results(3, 4, damageNeighbour(true)).
			Results(5, 6, damageNeighbour(false)))},
	}},
}
