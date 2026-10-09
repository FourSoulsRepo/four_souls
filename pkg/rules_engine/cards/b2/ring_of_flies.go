package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Ring Of Flies (Basic Monster Card)
//
//	Each time the attacking player rolls an attack roll of 3, they must steal a loot card from from another player at random.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var ringOfFlies = engine.CardDef{
	Ref:     "ring_of_flies",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}, {Kind: engine.RewardCents, Amount: 3}},
}
