package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Greedling (Basic Monster Card)
//
//	When this dies, the active player chooses a player. They lose 7¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var greedling = engine.CardDef{
	Ref:     "greedling",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      5,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 7}},
}
