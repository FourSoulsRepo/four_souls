package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Hornfel (Boss Card)
//
//	When this would deal combat damage to the active player, prevent it, then this deals 1 damage to a player chosen at random.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var hornfel = engine.CardDef{
	Ref:     "hornfel",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      5,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
}
