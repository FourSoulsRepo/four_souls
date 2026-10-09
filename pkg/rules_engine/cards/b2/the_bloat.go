package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Bloat (Boss Card)
//
//	Each time this deals combat damage, it deals 1 damage to each non-active player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theBloat = engine.CardDef{
	Ref:     "the_bloat",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      4,
	ATK:     2,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
}
