package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dark One (Boss Card)
//
//	Each time this takes damage, it gains +1{ATK} till end of turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var darkOne = engine.CardDef{
	Ref:     "dark_one",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
}
