package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Gurdy (Boss Card)
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var gurdy = engine.CardDef{
	Ref:     "gurdy",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      5,
	DC:      4,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 7}},
}
