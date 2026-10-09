package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mulligan (Basic Monster Card)
//
//	When this dies, expand monster slots by 1.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var mulligan = engine.CardDef{
	Ref:     "mulligan",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 3}},
}
