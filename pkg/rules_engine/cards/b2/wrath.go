package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Wrath (Boss Card)
//
//	When this dies, the active player rolls-
//	1-3: Each player takes 1 damage.
//	4-6: Each player takes 2 damage.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var wrath = engine.CardDef{
	Ref:     "wrath",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      3,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 6}},
}
