package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Chub (Boss Card)
//
//	Each time the attacking player rolls an attack roll of 1, this heals 2{HP}.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var chub = engine.CardDef{
	Ref:     "chub",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      3,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 7}},
}
