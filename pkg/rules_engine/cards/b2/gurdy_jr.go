package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Gurdy Jr. (Boss Card)
//
//	Each time the attacking player activates an item, they take 1 damage.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var gurdyJr = engine.CardDef{
	Ref:     "gurdy_jr",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      5,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
}
