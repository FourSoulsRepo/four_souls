package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Boil (Basic Monster Card)
//
//	Each time this deals combat damage, it heals 1{HP}.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var boil = engine.CardDef{
	Ref:     "boil",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 3}},
}
