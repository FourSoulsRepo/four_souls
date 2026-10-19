package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Gutted Fatty (Basic Monster Card)
//
//	When this dies expand monster slots by 1. The active player must make an additional attack.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var guttedFatty = engine.CardDef{
	Ref:     "gutted_fatty",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      5,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
}
