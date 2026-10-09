package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mom! (Epic Boss Card)
//
//	Combat damage this deals is doubled on attack rolls of 1.
//	When this dies, expand monsters slots by 1.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var mom = engine.CardDef{
	Ref:     "mom",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      5,
	DC:      4,
	ATK:     2,
	Soul:    2,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
}
