package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Dogma (Epic Boss Card)
//
//	When an attack is declared on this, put each other monster into discard.
//	When this dies, put a card named The Harbingers from outside the game into a monster slot not being attacked.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var dogma = engine.CardDef{
	Ref:     "dogma",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      3,
	ATK:     3,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 2}},
}
