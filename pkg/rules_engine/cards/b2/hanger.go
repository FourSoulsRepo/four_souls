package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Hanger (Basic Monster Card)
//
//	When this dies, expand shop slots by 1.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var hanger = engine.CardDef{
	Ref:     "hanger",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      4,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 7}},
}
