package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Delirium (Epic Boss Card)
//
//	Other monsters have +1{DC}.
//	When this dies, put it in the monster deck 6 cards from the top.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var delirium = engine.CardDef{
	Ref:     "delirium",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      5,
	DC:      4,
	ATK:     3,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 2}},
}
