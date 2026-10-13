package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Fallen (Boss Card)
//
//	When this dies, look at the top 5 cards of the monster deck and put them back in any order.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theFallen = engine.CardDef{
	Ref:     "the_fallen",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      5,
	ATK:     2,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
}
