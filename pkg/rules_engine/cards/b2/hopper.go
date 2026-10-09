package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Hopper (Basic Monster Card)
//
//	This takes no combat damage on attack rolls of 6.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var hopper = engine.CardDef{
	Ref:     "hopper",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 3}},
}
