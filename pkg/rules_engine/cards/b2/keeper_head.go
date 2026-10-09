package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Keeper Head (Basic Monster Card)
//
//	Each time this deals combat damage to a player, they lose 2¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var keeperHead = engine.CardDef{
	Ref:     "keeper_head",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Roll: true}},
}
