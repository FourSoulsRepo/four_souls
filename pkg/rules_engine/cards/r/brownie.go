package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Brownie (Boss Card)
//
//	Each time this takes damage, it gains -1{DC} till end of turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var brownie = engine.CardDef{
	Ref:     "brownie",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      6,
	DC:      5,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
}
