package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Clog (Boss Card)
//
//	When this dies, each player gives their hand to the player to their left.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var clog = engine.CardDef{
	Ref:     "clog",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      5,
	DC:      4,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 3}},
}
