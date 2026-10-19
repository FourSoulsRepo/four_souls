package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Charmed Globin (Holy/Charmed Monster Card)
//
//	When this dies, after gaining rewards, choose another player. They loot 3.
//	Each time the attacking player misses an attack roll, this heals 1{HP}.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var charmedGlobin = engine.CardDef{
	Ref:     "charmed_globin",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 3}},
}
