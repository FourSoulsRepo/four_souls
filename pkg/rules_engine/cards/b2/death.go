package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Death (Boss Card)
//
//	When this dies, the active player kills a player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var death = engine.CardDef{
	Ref:     "death",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      4,
	ATK:     2,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
}
