package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Cage (Boss Card)
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theCage = engine.CardDef{
	Ref:     "the_cage",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      6,
	DC:      4,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 2}},
}
