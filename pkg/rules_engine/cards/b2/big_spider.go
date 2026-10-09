package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Big Spider (Basic Monster Card)
//
//	When this dies, the active player may attack the monster deck an additional time.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var bigSpider = engine.CardDef{
	Ref:     "big_spider",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
}
