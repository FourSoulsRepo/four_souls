package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Carrion Queen (Boss Card)
//
//	This takes no combat damage on attack rolls of 4 or 5.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var carrionQueen = engine.CardDef{
	Ref:     "carrion_queen",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
}
