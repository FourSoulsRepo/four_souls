package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Blastocyst (Boss Card)
//
//	When this dies, expand monster slots by 2. The active player may attack the monster deck an additional time.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var blastocyst = engine.CardDef{
	Ref:     "blastocyst",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      5,
	DC:      4,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
}
