package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Tumor (Basic Monster Card)
//
//	Each time this takes damage, it deals 1 damage to each non-active player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var tumor = engine.CardDef{
	Ref:     "tumor",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      5,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 6}},
}
