package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Monstro II (Boss Card)
//
//	Each time this deals combat damage, the active player rolls-
//	1-3: Deal 1 damage to the player to your right.
//	4-6: Deal 1 damage to the player to your left.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var monstroIi = engine.CardDef{
	Ref:       "monstro_ii",
	Kind:      engine.MonsterCard,
	Copies:    1,
	HP:        5,
	DC:        4,
	ATK:       1,
	Soul:      1,
	Rewards:   []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	ThreePlus: true,
}
