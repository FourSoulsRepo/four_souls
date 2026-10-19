package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Charmed Bony (Holy/Charmed Monster Card)
//
//	When this dies, after gaining rewards, choose another player. They gain 5¢.
//	Each time the attacking player rolls an attack roll of 6, they deal 1 damage to a player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var charmedBony = engine.CardDef{
	Ref:     "charmed_bony",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 5}},
}
