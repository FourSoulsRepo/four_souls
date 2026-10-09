package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Greed (Boss Card)
//
//	Each time this deals damage, each player loses 4¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var greed = engine.CardDef{
	Ref:     "greed",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 9}},
}
