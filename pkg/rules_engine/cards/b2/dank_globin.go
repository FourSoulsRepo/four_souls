package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dank Globin (Basic Monster Card)
//
//	When this dies, the active player forces a player to discard 2 loot cards.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var dankGlobin = engine.CardDef{
	Ref:     "dank_globin",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      4,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
}
