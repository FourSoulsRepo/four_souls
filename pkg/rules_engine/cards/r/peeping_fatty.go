package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Peeping Fatty (Basic Monster Card)
//
//	When this dies, look at the top 5 cards of the loot deck and put them back in any order.
var peepingFatty = engine.CardDef{
	Ref:     "peeping_fatty",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "When this dies, look at the top 5 cards of the loot deck and put them back in any order.",
		Trigger: engine.WhenThisDies(),
		Effects: []engine.Effect{putBackInOrder(engine.LootDeck, 5)},
	}},
}
