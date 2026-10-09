package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Fat Bat (Basic Monster Card)
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var fatBat = engine.CardDef{
	Ref:     "fat_bat",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      5,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
}
