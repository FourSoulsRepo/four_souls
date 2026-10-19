package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Corny (Basic Monster Card)
//
//	When this dies, deal 2 damage to a non-active player chosen at random.
var corny = engine.CardDef{
	Ref:       "corny",
	Kind:      engine.MonsterCard,
	Copies:    1,
	HP:        1,
	DC:        4,
	ATK:       1,
	Rewards:   []engine.Reward{{Kind: engine.RewardCents, Amount: 4}},
	ThreePlus: true,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "When this dies, deal 2 damage to a non-active player chosen at random.",
		Trigger: engine.WhenThisDies(),
		Effects: []engine.Effect{randomPlayerTakes(2, true)},
	}},
}
