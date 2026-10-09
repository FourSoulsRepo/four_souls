package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Holy Mom’s Eye (Holy/Charmed Monster Card)
//
//	Each time a player rolls a ❷, they may recharge an item.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var holyMomsEye = engine.CardDef{
	Ref:     "holy_moms_eye",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      4,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
}
