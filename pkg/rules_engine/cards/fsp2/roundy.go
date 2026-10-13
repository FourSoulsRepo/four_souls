package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Roundy (Basic Monster Card)
//
//	Each time this deals combat damage, it deals 1 damage to each other monster.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var roundy = engine.CardDef{
	Ref:     "roundy",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      4,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 6}},
}
