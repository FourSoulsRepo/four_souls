package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Boom Fly (Basic Monster Card)
//
//	When this dies, it deals 1 damage to each player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var boomFly = engine.CardDef{
	Ref:     "boom_fly",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 4}},
}
