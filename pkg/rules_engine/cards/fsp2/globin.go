package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Globin (Basic Monster Card)
//
//	When this reaches 1{HP}, the active player rolls-
//	1-4: This takes 1 damage.
//	5-6: This heals 2{HP}.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var globin = engine.CardDef{
	Ref:     "globin",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 3}},
}
