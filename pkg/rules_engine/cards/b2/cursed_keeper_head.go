package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Cursed Keeper Head (Cursed Monster Card)
//
//	{Curse Effect}Each time a player rolls a ➀, they lose 2¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var cursedKeeperHead = engine.CardDef{
	Ref:     "cursed_keeper_head",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Roll: true}},
}
