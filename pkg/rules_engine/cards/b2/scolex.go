package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Scolex (Boss Card)
//
//	Each time this deals combat damage to a player, they discard a loot card.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var scolex = engine.CardDef{
	Ref:     "scolex",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      5,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
}
