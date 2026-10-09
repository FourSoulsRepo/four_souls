package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Pestilence (Boss Card)
//
//	When this dies, the active player deals 2 damage divided as they choose to any number of monsters or players.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var pestilence = engine.CardDef{
	Ref:     "pestilence",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      4,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
}
