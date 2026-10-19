package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Mushroom (Basic Monster Card)
//
//	Each time this deals combat damage to the attacking player, it deals 1 damage to each other player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var mushroom = engine.CardDef{
	Ref:     "mushroom",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
}
