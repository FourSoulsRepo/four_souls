package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Swarm Of Flies (Basic Monster Card)
//
//	Each time the attacking player rolls an attack roll of 5, they take 1 damage.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var swarmOfFlies = engine.CardDef{
	Ref:     "swarm_of_flies",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      5,
	DC:      2,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 5}},
}
