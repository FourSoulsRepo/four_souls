package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Charmed Fatty (Holy/Charmed Monster Card)
//
//	When this dies, after gaining rewards, choose another player. They loot 2.
//	Each time this deals damage to a player, they discard a loot card.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var charmedFatty = engine.CardDef{
	Ref:     "charmed_fatty",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      5,
	DC:      2,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
}
