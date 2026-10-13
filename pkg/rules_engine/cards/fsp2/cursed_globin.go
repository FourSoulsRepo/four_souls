package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Cursed Globin (Cursed Monster Card)
//
//	{Curse Effect}Each time a player rolls a ➂, each monster heals 2{HP}.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var cursedGlobin = engine.CardDef{
	Ref:     "cursed_globin",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Roll: true}},
}
