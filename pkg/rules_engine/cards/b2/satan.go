package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Satan! (Epic Boss Card)
//
//	Each time the attacking player rolls an attack roll of 6, they choose a living player. That player dies.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var satan = engine.CardDef{
	Ref:     "satan",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      6,
	DC:      4,
	ATK:     2,
	Soul:    2,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 2}},
}
