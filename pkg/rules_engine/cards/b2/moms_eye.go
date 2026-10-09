package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mom’s Eye (Basic Monster Card)
//
//	When this dies, the active player may look at a player's hand.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var momsEye = engine.CardDef{
	Ref:     "moms_eye",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      4,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
}
