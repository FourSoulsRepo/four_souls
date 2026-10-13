package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Krampus (Boss Card)
//
//	The active player must attack this once each turn if able.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var krampus = engine.CardDef{
	Ref:     "krampus",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      4,
	ATK:     2,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
}
