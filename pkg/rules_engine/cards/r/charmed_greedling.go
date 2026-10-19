package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Charmed Greedling (Holy/Charmed Monster Card)
//
//	When this dies, after gaining rewards, choose another player. They gain 10¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var charmedGreedling = engine.CardDef{
	Ref:     "charmed_greedling",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      6,
	ATK:     0,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 6}},
}
