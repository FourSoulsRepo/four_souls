package r

import engine "github.com/FourSoulsRepo/rules_engine"

// MOTHER! (Epic Boss Card)
//
//	Each time this takes damage, it gains +1{DC} till end of turn.
//	When this dies, expand monster slots by 2.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var mother = engine.CardDef{
	Ref:     "mother",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      8,
	DC:      1,
	ATK:     8,
	Soul:    2,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
}
