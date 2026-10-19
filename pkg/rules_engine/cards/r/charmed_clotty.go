package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Charmed Clotty (Holy/Charmed Monster Card)
//
//	When this dies, after gaining rewards, choose another player. They loot 1.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var charmedClotty = engine.CardDef{
	Ref:     "charmed_clotty",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
}
