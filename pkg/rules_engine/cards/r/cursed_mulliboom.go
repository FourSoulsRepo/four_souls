package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Cursed Mulliboom (Cursed Monster Card)
//
//	Each time a monster dies, this deals 3 damage to a player chosen at random.
var cursedMulliboom = engine.CardDef{
	Ref:     "cursed_mulliboom",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      4,
	ATK:     3,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 2}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time a monster dies, this deals 3 damage to a player chosen at random.",
		Trigger: engine.WhenAMonsterDies(),
		Effects: []engine.Effect{randomPlayerTakes(3, false)},
	}},
}
