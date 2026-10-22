package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Charmed Fatty (Holy/Charmed Monster Card)
//
//	When this dies, after gaining rewards, choose another player. They loot 2.
//	Each time this deals damage to a player, they discard a loot card.
var charmedFatty = engine.CardDef{
	Ref:     "charmed_fatty",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      5,
	DC:      2,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
	Abilities: []engine.Ability{
		afterRewardsAnotherPlayer("They loot 2.", engine.Loot(2)),
		{
			Kind:    engine.Triggered,
			Text:    "Each time this deals damage to a player, they discard a loot card.",
			Takes:   engine.TakesEventPlayer,
			Trigger: cardkit.WhenThisDealsDamage(false),
			Effects: []engine.Effect{cardkit.RollerDiscards},
		},
	},
}
