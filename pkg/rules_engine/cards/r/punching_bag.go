package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Punching Bag (Passive Treasure Card)
//
//	You have +2{HP}.
//	This item can be attacked.
var punchingBag = engine.CardDef{
	Ref:        "punching_bag",
	Kind:       engine.TreasureCard,
	Copies:     1,
	HP:         2,
	DC:         5,
	ATK:        1,
	Attackable: true,
	Rewards:    []engine.Reward{{Kind: engine.RewardCents, Amount: 5}},
	Statics:    []engine.Static{engine.YouHave(engine.StatPlayerHP, 2)},
}
