package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Black Bony (Basic Monster Card)
//
//	When this dies, it deals 1 damage to the player who killed it.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var blackBony = engine.CardDef{
	Ref:     "black_bony",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Roll: true}},
}
