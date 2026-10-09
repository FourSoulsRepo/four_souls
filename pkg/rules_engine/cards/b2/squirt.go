package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Squirt (Basic Monster Card)
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var squirt = engine.CardDef{
	Ref:     "squirt",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
}
