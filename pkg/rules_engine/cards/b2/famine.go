package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Famine (Boss Card)
//
//	When this dies, the active player skips their next turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var famine = engine.CardDef{
	Ref:     "famine",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 3}},
}
