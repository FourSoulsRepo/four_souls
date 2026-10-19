package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Charmed Keeper (Holy/Charmed Monster Card)
//
//	When this dies, after gaining rewards, choose another player. They roll and gain ¢ equal to the result.
//	Each time this deals combat damage to a player, they lose 2¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var charmedKeeper = engine.CardDef{
	Ref:     "charmed_keeper",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      3,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Roll: true}},
}
