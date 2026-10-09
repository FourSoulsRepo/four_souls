package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Holy Squirt (Holy/Charmed Monster Card)
//
//	Each time a player rolls a ❺, they loot 1.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var holySquirt = engine.CardDef{
	Ref:     "holy_squirt",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
}
