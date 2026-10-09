package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Daddy Long Legs (Boss Card)
//
//	Each time the attacking player rolls an attack roll of 1, each monster gains +1{DC} till end of turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var daddyLongLegs = engine.CardDef{
	Ref:     "daddy_long_legs",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      4,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 7}},
}
